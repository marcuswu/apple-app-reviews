# Recent App Review Backend #

This service keeps a local cache of reviews from the Apple App Review RSS feed. It refreshes its cache periodically so that frequently requested app reviews do not cause undue traffic. The backend can be stopped
and restarted without losing its cache, but it will request new data if reviews are requested for an app whose
cache is stale.

Only standard libraries are used, but I would have added in zerolog, dotenv, and mux or gin to avoid
re-inventing the wheel if it was not requested that I avoid it. I would have also checked to see if there was
anything that already handles local caching, but if this were a real backend, something like redis would likely
be a better choice.

Though the instructions asked for the project to work on a specific app, the frontend and backend does support
entering and requesting 48 hours worth of reviews for an arbitrary app.

The majority of the code is procedural because there didn't seem to be a need to instantiate an object to handle
most of the review updating logic. The models are more object oriented in nature.

I really wanted to get to some automated testing, but I skipped it initiall to avoid taking too much time on
this. I have added some tests after submitting for review.

From an architecture standpoint, this is pretty simple. For a project of this size, this is probably acceptable, but if this grew much larger some additional separation of concerns could improve life.

File / cache management is also extremely simple. Additional work for removing older cache to ensure file system
health would be good to add. This could be as simple as removing any App-[0-9]*.json older than some predefined
duration.

Finally, the service just looks for `GET /{appId}`. In an ideal restful setup, it would be more like `GET /reviews/{appId}`. This could be accomplished by setting a proxy in front of it like such as AWS API Gateway or
by rewriting this project a little bit.

## Running it ##
Just run `go run main.go`

By default, it is configured to run on port `8000`

## Design review exercise ##
### Reflective Thoughts ###
Giving myself a short timeframe, I expected to have some flaws to the approach and implementation. I wrote this with that spirit in mind. I took an agile, incremental approach to writing something quickly with iteration for improvement in mind.
* ### Testing ###
    The testing and defect resolution are the two ways I allowed myself to make significant changes after completion of my deadline.
	* The `models` package is ok. Load and Save take streams making testing easy.
    * `updater` package:
        Much of this package is integration work, but some of its functionality could be separated to make testing of those integrations easier
        * `nextApp` was harder and requires the test to set up filesystem prerequisites
            * Improvement: create an interface for fs interoperability. Pass a test implementation to avoid real fs access during tests
        * `FetchAppReviews` suffers a similar problem to nextApp. It creates an http request.
            * Could pass req in, but we may need multiples pages limiting the usefulness of that approach
            * Extending that idea, we can pass a req. factory function or factory interface.
        * `fileForAppId` was so simple, tests seemed less useful, but could easily be written if an edge case was found
            * Technically an empty appId functions, but we may want to filter those and return an error which would make testing more useful
        * `SaveReviews` is really just a wrapper around `models.SaveReviews` to send in the file stream. There isn't much logic to test. Integration testing covers this.
        * `LoadReviews` has more logic. The logic of this method has more than one concern. Separating them would allow for improved testing. The two concerns:
            * loading the reviews from file
            * deciding when to refresh cache
        * `UpdateNext` by nature is an integration of previous functionality, but with the earlier stated improvement to FetchAppReviews, it could more easily be tested.
    * `main` package:
        * Use a named function for the go routine to be able to integration test it

* ### Higher level architecture / design considerations: ###
	* Went for as minimal initially and worked my way towards complexity
		* Some of it was done in my head rather than with code and commits to show
			* Would have been better to show that evolution
			* Start with a
				* constant app id
				* static cache
			* Use commits to show more thought process and evolution of the project
	* I started with the models, knowing that I would need to read and return reviews. Kind of bottom up approach here.
		* Used a different storage model to make it easier to deal with cache
		* Processing / filtering up front makes the optimal path (having recent cached reviews) faster and adds hopefully negligible delay to longer path
	* My approach might have been more complex than required
		* The minimal pproach probably was to pre-fetch a static app id with no cache invalidation
		* That didn't feel good to me, so I added:
			* minimalistic cache handling
				* simple file modification date check
				* go routine cache refresher for kicks
    				* probably needs to ensure requests don't get partial data
        				* locks for reading / writing
        				* write to separate file, move it in when done (use the fs / os to handle atomicity)
    				* in a more complex project, a document database or redis store might be used
      				* a more complex cache refresher could be its own process
			* app id selection
				* just pass it in; fill url template
				* url really should be in config
					* it felt arbitrary not to but I had to draw the line or forever improve it
						* in normal dev, there would be collab with other engs & stakeholders to determine where to draw these lines
* ### After completion, I found: ###
    After each issue, I list how long after introduction they were found and fixed. Each issue also describes the impact of the defect.
    * Off by one error when filtering app reviews by age (2 days after introduction)
        * Before this fix, one review within 48 hrs would be left out
        * After this fix, backend crashes when no reviews are left after filter
    * Index out of bounds (-1) due to previous fix (1 day after introduction)
        * After first update and before fix requests for reviews where there were none would crash the server
    * Finding files by age used the wrong date extreme (3 days after introduction)
        * Prior to this, the app cache checks would never invalidate files
        * Cache would still be updated by app review requests
    * Used wrong age constant when looking for next app to update (4 days after introduction)
        * The code functioned as intended, but with the wrong value
        * The cache took longer to invalidate so some stale data could be reported
        * The issue was masked by previous defects
    * Called UpdateNext without App Id list (2 days after introduction)
        * File changes were missed getting added to the previous commit
        * Should have been caught by:
            * Just running it :rage:
            * More testing
            * Ensuring no outstanding repo changes after commit
            * On a larger, continuous project by CICD
