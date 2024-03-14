Your service should store data about the reviews it fetches for an app (something as simple as writing to an external file is perfectly fine). The app should be able to be stopped/restarted without losing its progress and state.
heReviews fetched and displayed should be ordered by newest first, and for each review the output should include the review content, author, score, and time the review was submitted.

You should be able to do all of this using the standard libraries of the language you’re building in. If you do use 3rd party libraries, just be able to justify their usage.

Think about how to support any number of apps, and how this would affect your design.

For this assignment, you’ll be creating:
- A backend service/app that polls an iOS app’s App Store Connect RSS feed to fetch and store App Store reviews for a specific iOS app
- A React app that calls an endpoint on the backend app to fetch and display new reviews from the last 48 hours

We'll need two main pieces:
* A go routine to periodically fetch reviews
  * Load stored reviews
  * Purge entries older than 48 hours
  * Fetch new reviews
  * Merge the two lists
  * Save the results
  * A cancel channel to listen to for exiting
* An http service to serve reviews
  * Receives a request for reviews
  * Reads the review list
  * Sends the review list
* Listen for an interrupt to close the go routine
  * On receipt of interrupt, send cancel signal to channel
  * exit