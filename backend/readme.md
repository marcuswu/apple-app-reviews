# Recent App Review Backend #

This service keeps a local cache of reviews from the Apple App Review RSS feed. It refreshes its cache periodically so that frequently requested app reviews do not cause undue traffic. The backend can be stopped
and restarted without losing its cache, but it will request new data if reviews are requested for an app whose
cache is stale.

Only standard libraries are used, but I would have added in zerolog, dotenv, and mux to avoid re-inventing the
wheel if it was not requested that I avoid it. I would have also checked to see if there was anything that
already handles local caching, but if this were a real backend, something like redis would likely be a better
choice.

Though the instructions asked for the project to work on a specific app, the frontend and backend does support
entering and requesting 48 hours worth of reviews for an arbitrary app.

The majority of the code is procedural because there didn't seem to be a need to instantiate an object to handle
most of the review updating logic. The models are more object oriented in nature.

I really wanted to get to some automated testing, but I am skipping it to avoid taking too much time on this.

From an architecture standpoint, this is pretty simple. For a project of this size, this is probably acceptable, but if this grew much larger some additional separation of concerns could improve life.

File / cache management is also extremely simple. Additional work for removing older cache to ensure file system
health would be good to add. This could be as simple as removing any App-[0-9]*.json older than some predefined
duration.

Finally, the service just looks for `GET /{appId}`. In an ideal restful setup, it would be more like `GET /reviews/{appId}`. This could be accomplished by setting a proxy in front of it like such as AWS API Gateway or
by rewriting this project a little bit.

## Running it ##
Just run `go run main.go`