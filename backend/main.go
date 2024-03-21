package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/marcuswu/app-reviews/config"
	"github.com/marcuswu/app-reviews/updater"
)

// Request handler for looking up app reviews for an app.
// Prefers local cache if within config.MAX_REVIEW_FILE_AGE_MINUTES
// If local cache doesn't exist or is stale, fetch reviews from Apple and cache them
func reviewRequestHandler(res http.ResponseWriter, req *http.Request) {
	appId := req.PathValue("appId")
	fmt.Printf("handling request for app id %s\n", appId)

	reviews, err := updater.LoadReviews(appId)
	if err != nil {
		reviews, err = updater.FetchAppReviews(appId)
		if err != nil {
			http.Error(res, fmt.Sprintf("Failed to fetch app reviews: %s", err), http.StatusFailedDependency)
			return
		}
		updater.SaveReviews(appId, reviews)
	}
	json.NewEncoder(res).Encode(reviews)
}

func main() {
	var wg sync.WaitGroup
	exitchan := make(chan bool, 1)

	// *** Start up review fetching ***
	wg.Add(1)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			updater.UpdateNext(updater.ListAppCache())

			// If there is anything on exitchan, we should stop
			select {
			case <-exitchan:
				wg.Done()
				return
			default:
				continue
			}
		}
	}()

	// *** Start up interrupt handler so we can shut down gracefully ***
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan // Wait for a SIGINT

		// Close out our go routine gracefully
		exitchan <- true
		wg.Wait()

		os.Exit(0)
	}()

	// *** Start up request handler ***
	http.HandleFunc("/{appId}", reviewRequestHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", config.SERVER_PORT), nil)
}
