package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcuswu/app-reviews/config"
	"github.com/marcuswu/app-reviews/models"
)

func TestReviewIntegration(t *testing.T) {
	req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:%d/articles", config.SERVER_PORT), nil)

	w := httptest.NewRecorder()
	reviewRequestHandler(w, req)

	response := w.Result()
	if http.StatusOK != response.StatusCode {
		t.Errorf("expected OK status code (200), got %d", response.StatusCode)
	}
	defer response.Body.Close()

	reviews, err := models.LoadReviews(response.Body)
	if err != nil {
		t.Errorf("expected no error reading from body, got %s", err)
	}

	for _, review := range reviews {
		if time.Since(review.Updated).Hours() > 48 {
			t.Errorf("expected no review older than 48 hours. found one %f hours old", time.Since(review.Updated).Hours())
		}
	}
}
