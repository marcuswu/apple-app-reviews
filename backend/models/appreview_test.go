package models

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestModelMarshalling(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get current directory: %s", err)
	}
	feedData, err := os.OpenFile(filepath.Join(dir, "app-feed.json"), os.O_RDONLY, 0000)
	if err != nil {
		t.Errorf("failed to read test data: %s", err)
	}
	defer feedData.Close()

	time1, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-13T04:25:02-07:00")
	time2, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-12T10:10:58-07:00")
	time3, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-10T15:27:53-07:00")
	time4, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-10T13:42:31-07:00")
	time5, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-09T08:50:25-07:00")
	var expectedReviews = AppReviews{
		AppReview{
			Author:  Author{Name: "Test Author Foo", Uri: "unused"},
			Updated: time1,
			Rating:  1,
			Id:      "11039586140",
			Title:   "Test Review Title 1",
			Content: "Test Review One",
		},
		AppReview{
			Author:  Author{Name: "Test Author Blah", Uri: "unused"},
			Updated: time2,
			Rating:  5,
			Id:      "11037094603",
			Title:   "Test Review Title 2",
			Content: "Test Review Two",
		},
		AppReview{
			Author:  Author{Name: "Test Author Bar", Uri: "unused"},
			Updated: time3,
			Rating:  4,
			Id:      "11030840001",
			Title:   "Test Review Title 3",
			Content: "Test Review Three",
		},
		AppReview{
			Author:  Author{Name: "Test Author Baz", Uri: "unused"},
			Updated: time4,
			Rating:  1,
			Id:      "11030579850",
			Title:   "Test Review Title 4",
			Content: "Test Review Four",
		},
		AppReview{
			Author:  Author{Name: "Review Author Blah", Uri: "unused"},
			Updated: time5,
			Rating:  3,
			Id:      "11026038445",
			Title:   "Test Review Title 5",
			Content: "Test\nReview\nFive",
		},
	}

	data, err := io.ReadAll(feedData)

	if err != nil {
		t.Errorf("expected no error reading feed data: %s", err)
	}

	reviews := AppReviewFeed{}
	if err = json.Unmarshal(data, &reviews); err != nil {
		t.Errorf("expected no error parsing feed data: %s", err)
	}
	feed := reviews.Reviews

	if len(feed) != 5 {
		t.Errorf("expected just one review in the feed, got %d", len(feed))
	}

	for index, review := range feed {
		if expectedReviews[index].Author.Name != review.Author.Name {
			t.Errorf("unexpected author name. expected %s and found %s", expectedReviews[index].Author.Name, review.Author.Name)
		}

		if !expectedReviews[index].Updated.Equal(review.Updated) {
			t.Errorf("unexpected review update date. expected %s and found %s", expectedReviews[index].Updated, review.Updated)
		}

		if expectedReviews[index].Rating != review.Rating {
			t.Errorf("unexpected review rating. expected %d and found %d", expectedReviews[index].Rating, review.Rating)
		}

		if expectedReviews[index].Title != review.Title {
			t.Errorf("unexpected review title. expected %s and found %s", expectedReviews[index].Title, review.Title)
		}

		if expectedReviews[index].Content != review.Content {
			t.Errorf("unexpected review content. expected %s and found %s", expectedReviews[index].Content, review.Content)
		}
	}
}

func TestLoadReviews(t *testing.T) {
	file := bytes.NewBuffer([]byte{})

	time1, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-13T04:25:02-07:00")
	time2, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-12T10:10:58-07:00")
	time3, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-10T15:27:53-07:00")
	time4, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-10T13:42:31-07:00")
	time5, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-09T08:50:25-07:00")
	var expectedReviews = AppReviews{
		AppReview{
			Author:  Author{Name: "Test Author Foo", Uri: "unused"},
			Updated: time1,
			Rating:  1,
			Id:      "11039586140",
			Title:   "Test Review Title 1",
			Content: "Test Review One",
		},
		AppReview{
			Author:  Author{Name: "Test Author Blah", Uri: "unused"},
			Updated: time2,
			Rating:  5,
			Id:      "11037094603",
			Title:   "Test Review Title 2",
			Content: "Test Review Two",
		},
		AppReview{
			Author:  Author{Name: "Test Author Bar", Uri: "unused"},
			Updated: time3,
			Rating:  4,
			Id:      "11030840001",
			Title:   "Test Review Title 3",
			Content: "Test Review Three",
		},
		AppReview{
			Author:  Author{Name: "Test Author Baz", Uri: "unused"},
			Updated: time4,
			Rating:  1,
			Id:      "11030579850",
			Title:   "Test Review Title 4",
			Content: "Test Review Four",
		},
		AppReview{
			Author:  Author{Name: "Review Author Blah", Uri: "unused"},
			Updated: time5,
			Rating:  3,
			Id:      "11026038445",
			Title:   "Test Review Title 5",
			Content: "Test\nReview\nFive",
		},
	}

	err := SaveReviews(file, expectedReviews)
	if err != nil {
		t.Errorf("expected no error saving feed data: %s", err)
	}

	feed, err := LoadReviews(file)
	if err != nil {
		t.Errorf("expected no error reading feed data: %s", err)
	}

	if len(feed) != 5 {
		t.Errorf("expected just one review in the feed, got %d", len(feed))
	}

	for index, review := range feed {
		if expectedReviews[index].Author.Name != review.Author.Name {
			t.Errorf("unexpected author name. expected %s and found %s", expectedReviews[index].Author.Name, review.Author.Name)
		}

		if !expectedReviews[index].Updated.Equal(review.Updated) {
			t.Errorf("unexpected review update date. expected %s and found %s", expectedReviews[index].Updated, review.Updated)
		}

		if expectedReviews[index].Rating != review.Rating {
			t.Errorf("unexpected review rating. expected %d and found %d", expectedReviews[index].Rating, review.Rating)
		}

		if expectedReviews[index].Title != review.Title {
			t.Errorf("unexpected review title. expected %s and found %s", expectedReviews[index].Title, review.Title)
		}

		if expectedReviews[index].Content != review.Content {
			t.Errorf("unexpected review content. expected %s and found %s", expectedReviews[index].Content, review.Content)
		}
	}
}

func TestReviewFiltering(t *testing.T) {
	time1, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-13T04:25:02-07:00")
	time2, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-12T10:10:58-07:00")
	time3, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-10T15:27:53-07:00")
	time4, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-10T13:42:31-07:00")
	time5, _ := time.Parse("2006-01-02T15:04:05-07:00", "2024-03-09T08:50:25-07:00")
	var reviews = AppReviews{
		AppReview{
			Author:  Author{Name: "Test Author Foo", Uri: "unused"},
			Updated: time1,
			Rating:  1,
			Id:      "11039586140",
			Title:   "Test Review Title 1",
			Content: "Test Review One",
		},
		AppReview{
			Author:  Author{Name: "Test Author Blah", Uri: "unused"},
			Updated: time2,
			Rating:  5,
			Id:      "11037094603",
			Title:   "Test Review Title 2",
			Content: "Test Review Two",
		},
		AppReview{
			Author:  Author{Name: "Test Author Bar", Uri: "unused"},
			Updated: time3,
			Rating:  4,
			Id:      "11030840001",
			Title:   "Test Review Title 3",
			Content: "Test Review Three",
		},
		AppReview{
			Author:  Author{Name: "Test Author Baz", Uri: "unused"},
			Updated: time4,
			Rating:  1,
			Id:      "11030579850",
			Title:   "Test Review Title 4",
			Content: "Test Review Four",
		},
		AppReview{
			Author:  Author{Name: "Review Author Blah", Uri: "unused"},
			Updated: time5,
			Rating:  3,
			Id:      "11026038445",
			Title:   "Test Review Title 5",
			Content: "Test\nReview\nFive",
		},
	}

	filterTime, _ := time.Parse(time.RFC3339, "2024-03-10T15:27:00-07:00")
	filtered := reviews.After(filterTime)

	if len(filtered) > 3 {
		for _, r := range filtered {
			t.Logf("date %s", r.Updated)
		}
		t.Errorf("expected to find 3 reviews and found %d", len(filtered))
	}

	filtered = AppReviews{}.After(filterTime)
	if len(filtered) > 0 {
		t.Errorf("expected to find no reviews and found %d", len(filtered))
	}
}
