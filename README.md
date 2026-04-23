# aggregatr
## Overview


Aggregatr is an RSS Feed Aggregator that was a guided project with [Boot.dev](https://www.boot.dev/tracks/backend) to create the "gator" tool. The goal of the project was to learn how to use [Postgres](https://www.postgresql.org/) Databases and Queries within [Golang](https://go.dev/) by building a CLI tool.


## Installation Instructions


(I will only give Linux commands because we are all adults here)
Before the installaton of Aggregatr you must first download Go:
`curl -sS https://webi.sh/golang | sh`


After downloading go ensure that Postgres is downloaded:
`sudo apt install postgresql postgresql-contrib`


Once all those boxes are checked now it is time to install the tool:
`go install github.com/PlinyTheYounger0/aggregatr`


## Configuration

Aggregatr assumes there is a .gatorconfig.json file present in your home (~) directory. However you can configure the filename in internal/config/config.go. 


The .gatorconfig.json holds the db_url and the current_user_name. The db_url must be set by the user but the current_user_name is set by the program itself so it doesn't require any tinkering.

## Commands

The Bread and Butter Baby

- Register: `aggregatr register <user-name>`
Register takes the username and creates an account. It then logs in the username as the current user.


- Login: `aggregatr login <user-name>`
Logs in the desired user or will error if the account hasn't been registered yet.


- Reset: `aggregatr reset`
Resets the users registered.


- Users: `aggregatr users`
Lists users registered with the current user highlighted.


- Add Feed: `aggregatr addfeed <feed-name> <feed-url>`
Adds an RSS feed to be aggregated with the feed-name and feed-url.


- Feeds: `aggregatr feeds`
Lists the feeds being aggregated by the program.


- Follow: `aggregatr follow <feed-url>`
Adds the desired RSS feed to the users followed list.


- Unfollow `aggregatr unfollow <feed-url>`
Removes the desired RSS feed from the users followed list.


- Following `aggregatr following`
Lists the RSS feeds being followed by the current user.


- Aggregate: `aggregatr agg <time-duration>`
Fetches the least recent RSS Feed and updates it while storing the posts to be browsed by the user. Time between fetches is dictated by time-duration and only fetches one RSS Feed at a time. Time-duration formatting can be found at [Go Time](https://pkg.go.dev/time#ParseDuration).


- Browse: `aggregatr browse <posts>`
Prints out the users most recent posts. Will display the number of posts specified. Posts variable must be a positive integer. 


## Next Steps
TBD