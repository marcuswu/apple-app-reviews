// Package updater provides most of the logic for fetching reviews and managing app review cache
package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/marcuswu/app-reviews/config"
	"github.com/marcuswu/app-reviews/models"
)

// appFiles is a small utility function for returning a list of cache files
func appFiles() ([]string, error) {
	return filepath.Glob("./App-[0-9]*.json")
}

// nextApp returns the next app cache to refresh or an error if there is nothing to update
func nextApp(files []string) (string, error) {
	oldest := time.Unix(0, 0)
	oldestId := ""
	for _, filename := range files {
		fi, err := os.Stat(filename)
		if err != nil {
			continue
		}
		if !fi.Mode().IsRegular() {
			continue
		}
		if fi.ModTime().Before(oldest) {
			oldest = fi.ModTime()
			oldestId = strings.TrimSuffix(strings.TrimPrefix(fi.Name(), "App-"), ".json")
		}
	}

	if len(oldestId) < 1 {
		return oldestId, errors.New("could not find an app to refresh")
	}

	if oldest.After(time.Now().Add(time.Duration(-config.OLDEST_REVIEW_HOURS) * time.Minute)) {
		// The oldest file has been refreshed too recently to refresh again
		return oldestId, errors.New("could not find an app to refresh")
	}

	return oldestId, nil
}

// FetchAppReviews retrieves reviews within config.OLDEST_REVIEW_HOURS age for the provided app id
func FetchAppReviews(appId string) (models.AppReviews, error) {
	page := 1
	reviews := make(models.AppReviews, 0, config.OLDEST_REVIEW_HOURS)
	for needMore := true; needMore; page++ {
		url := fmt.Sprintf("https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=%d/json",
			appId, page)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return []models.AppReview{}, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return []models.AppReview{}, err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return []models.AppReview{}, err
		}

		feed := models.AppReviewFeed{}
		if err := json.Unmarshal(resBody, &feed); err != nil {
			return []models.AppReview{}, err
		}

		for _, review := range feed.Reviews {
			reviews = append(reviews, models.AppReview(review))
		}
		// Keep requesting more reviews until we find a page with a review older than we need
		if len(reviews) < 1 {
			needMore = false
			continue
		}
		needMore = time.Since(reviews[len(reviews)-1].Updated).Hours() < config.OLDEST_REVIEW_HOURS

		fmt.Printf("Have %d reviews after page %d\n", len(reviews), page)
		reviews = reviews.After(time.Now().Add(time.Duration(-config.OLDEST_REVIEW_HOURS) * time.Hour))
		fmt.Printf("Have %d reviews after filtering\n", len(reviews))
	}
	fmt.Printf("Returning %d reviews", len(reviews))

	return reviews, nil
}

// fileForAppId returns the filename to store or retrieve app reviews to for a given app id
func fileForAppId(appId string) string {
	return fmt.Sprintf("App-%s.json", appId)
}

// SaveReviews saves a list of app reviews to cache
func SaveReviews(appId string, reviews models.AppReviews) error {
	filename := fileForAppId(appId)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	return models.SaveReviews(file, reviews)
}

// LoadReviews loads an app's cached app reviews.
// Returns an error if unable to read the reviews or if the cache is too stale to use
func LoadReviews(appId string) (models.AppReviews, error) {
	filename := fileForAppId(appId)

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Unable to find file %s\n", filename)
		return nil, err
	}

	modifiedtime := fileInfo.ModTime()
	if time.Since(modifiedtime).Minutes() > config.MAX_REVIEW_FILE_AGE_MINUTES {
		fmt.Printf("Refresh stale file\n")
		return nil, errors.New("stale file -- refresh it")
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0000)
	if err != nil {
		fmt.Printf("Unable to find file %s\n", filename)
		return nil, err
	}

	return models.LoadReviews(file)
}

// Look at cached app reviews and refresh the oldest one that is expired (if any)
func UpdateNext() error {
	apps, err := appFiles()
	if err != nil {
		return err
	}

	app, err := nextApp(apps)
	if err != nil {
		return err
	}

	reviews, err := FetchAppReviews(app)
	if err != nil {
		return err
	}

	reviews = reviews.After(time.Now().Add(time.Duration(-config.OLDEST_REVIEW_HOURS) * time.Hour))
	return SaveReviews(app, reviews)
}
