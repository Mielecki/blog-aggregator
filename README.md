# Blog aggregator

Blog aggregator is a command-line tool built with Go, designed to fetch and menage data from RSS feeds, storing it in a PostgreSQL database.

This project uses [sqlc](https://docs.sqlc.dev/en/latest/index.html) to generate Go code based on the provided SQL queries. It also uses [goose](https://github.com/pressly/goose/) for easy up/down database migrations.

## Commands

* `register [username]`: Registers new user in the database.
* `login [username]`: Sets the user at the current user.
* `users`: Lists all registered users.
* `agg [time]`: Fetches the oldest feed.
* `addfeed [name] [url]`: Adds a new feed to users following list.
* `feeds`: Lists all existing feeds.
* `follow [url]`: Add the feed to users following list.
* `following`: Lists all feeds followed by user.
* `unfollow [url]`: Removes the feed from users following list.
* `browse [limit]`: Browses feeds from the following list.

## Learing goals

* Learning Go,
* Working with APIs,
* Web Scraping,
* SQL & PostgreSQL

This project was developed as part of the Boot.dev online course.