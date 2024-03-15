package models

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

type LabeledField string

func (lf *LabeledField) UnmarshalJSON(data []byte) error {
	var labeledField struct {
		Label string `json:"label"`
	}

	if err := json.Unmarshal(data, &labeledField); err != nil {
		return err
	}

	*lf = LabeledField(labeledField.Label)
	return nil
}

type ReviewTime time.Time

func (rt *ReviewTime) UnmarshalJSON(data []byte) error {
	var reviewTime struct {
		Label time.Time `json:"label"`
	}

	if err := json.Unmarshal(data, &reviewTime); err != nil {
		return err
	}

	*rt = ReviewTime(reviewTime.Label)
	return nil

}

type ReviewRating int

func (rr *ReviewRating) UnmarshalJSON(data []byte) error {
	var err error
	var rating int
	var reviewRating struct {
		Label string `json:"label"`
	}

	if err = json.Unmarshal(data, &reviewRating); err != nil {
		return err
	}

	if rating, err = strconv.Atoi(reviewRating.Label); err != nil {
		return err
	}

	*rr = ReviewRating(rating)

	return nil
}

type ReviewLink string

func (rl *ReviewLink) UnmarshalJSON(data []byte) error {
	var appLink struct {
		Attributes struct {
			HRef string `json:"href"`
		} `json:"attributes"`
	}

	if err := json.Unmarshal(data, &appLink); err != nil {
		return err
	}

	*rl = ReviewLink(appLink.Attributes.HRef)
	return nil
}

type AppleAuthor struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

func (a *AppleAuthor) UnmarshalJSON(data []byte) error {
	var appleAuthor struct {
		Name LabeledField `json:"name"`
		Uri  LabeledField `json:"uri"`
	}

	if err := json.Unmarshal(data, &appleAuthor); err != nil {
		return err
	}

	a.Name = string(appleAuthor.Name)
	a.Uri = string(appleAuthor.Uri)

	return nil
}

type Author struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type AppleAppReview struct {
	Author  Author    `json:"author"`
	Updated time.Time `json:"updated"`    // updated/label
	Rating  int       `json:"im:rating"`  // im:rating/label
	Version string    `json:"im:version"` // im:version/label
	Id      string    `json:"id"`         // id/label
	Title   string    `json:"title"`      // title/label
	Content string    `json:"content"`    // content/label
	Link    string    `json:"link"`       // link/attributes/href
}

func (ar *AppleAppReview) String() string {
	return fmt.Sprintf("{\n\tAuthor Name: \"%s\",\n\tAuthor URL: \"%s\",\n\tLast Update: \"%s\",\n\tRating: %d,\n\tVersion: \"%s\",\n\tId: \"%s\",\n\tTitle: \"%s\",\n\tContent: \"%s\",\n\tLink: \"%s\"\n}", ar.Author.Name, ar.Author.Uri, ar.Updated, ar.Rating, ar.Version, ar.Id, ar.Title, ar.Content, ar.Link)
}

func (ar *AppleAppReview) UnmarshalJSON(data []byte) error {
	var appleAppReview struct {
		Author  AppleAuthor  `json:"author"`
		Updated ReviewTime   `json:"updated"`    // updated/label
		Rating  ReviewRating `json:"im:rating"`  // im:rating/label
		Version LabeledField `json:"im:version"` // im:version/label
		Id      LabeledField `json:"id"`         // id/label
		Title   LabeledField `json:"title"`      // title/label
		Content LabeledField `json:"content"`    // content/label
		Link    ReviewLink   `json:"link"`       // link/attributes/href
	}

	if err := json.Unmarshal(data, &appleAppReview); err != nil {
		return err
	}

	ar.Author = Author(appleAppReview.Author)
	ar.Updated = time.Time(appleAppReview.Updated)
	ar.Rating = int(appleAppReview.Rating)
	ar.Version = string(appleAppReview.Version)
	ar.Id = string(appleAppReview.Id)
	ar.Title = string(appleAppReview.Title)
	ar.Content = string(appleAppReview.Content)
	ar.Link = string(appleAppReview.Link)

	return nil
}

type AppReview struct {
	Author  Author    `json:"author"`
	Updated time.Time `json:"updated"`
	Rating  int       `json:"rating"`
	Version string    `json:"version"`
	Id      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Link    string    `json:"link"`
}

type AppReviews []AppReview
type AppleAppReviews []AppleAppReview

func (r AppReviews) Uniq() AppReviews {
	uniqueReviewIds := make(map[string]struct{})
	uniqueReviews := make([]AppReview, 0, len(r))
	for _, review := range r {
		if _, exists := uniqueReviewIds[string(review.Id)]; exists {
			continue
		}
		uniqueReviewIds[string(review.Id)] = struct{}{}
		uniqueReviews = append(uniqueReviews, review)
	}

	return uniqueReviews
}

func (r AppReviews) After(minTime time.Time) AppReviews {
	sort.Slice(r, func(i, j int) bool { return time.Time(r[i].Updated).After(time.Time(r[j].Updated)) })
	fmt.Printf("First updated date is %f hours ago\n", time.Since(r[0].Updated).Hours())
	fmt.Printf("Last updated date is %f hours ago\n", time.Since(r[len(r)-1].Updated).Hours())
	end := 0
	for idx, review := range r {
		end = idx
		if time.Time(review.Updated).Before(minTime) {
			break
		}
	}

	return r[:end]
}

type AppReviewFeed struct {
	Reviews []AppleAppReview
}

func (arf *AppReviewFeed) UnmarshalJSON(data []byte) error {
	var appReviewFeed struct {
		Feed struct {
			Entry AppleAppReviews `json:"entry"`
		} `json:"feed"`
	}

	if err := json.Unmarshal(data, &appReviewFeed); err != nil {
		fmt.Printf("Error unmarshalling app review feed: %s\n", err)
		return err
	}

	arf.Reviews = appReviewFeed.Feed.Entry

	return nil
}

func LoadReviews(stream io.ReadCloser) (AppReviews, error) {
	data, err := io.ReadAll(stream)

	if err != nil {
		fmt.Printf("LoadReviews could not read stream: %s\n", err)
		return nil, err
	}

	reviews := make(AppReviews, 0, 10)
	if err = json.Unmarshal(data, &reviews); err != nil {
		fmt.Printf("LoadReviews could not unmarshal json: %s\n", err)
		return nil, err
	}

	return reviews, nil
}

func SaveReviews(stream io.WriteCloser, reviews AppReviews) error {
	data, err := json.Marshal(reviews)

	if err != nil {
		return err
	}

	if _, err = stream.Write(data); err != nil {
		return err
	}

	return nil
}
