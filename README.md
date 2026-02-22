# Chirpy
Chirpy is a social network similar to twitter. This project is a production style HTTP server in Go, without the use of a framework.

## Highlights
* Basic Go HTTP Server with social media API func.
* PostgreSQL database to manage users and chirps created.
* Use JWTs to handle authentication.
* Use Webhooks to allow premium users to edit posted chirps.

## Overview
This work was done as apart of a course on boot.dev.

## Usage
Chirpy is technical a monolith project, but a lot of the development was focused on the API. In the future, I would like to make a seperate front end that communicates with this backend.

### Docker

### Locally

## Installing
### Goose
Goose is a database migration tool written in Go. It runs migrations from a set of SQL files, making it a perfect fit for this project (we wanna stay close to the raw SQL).

A migration is just a set of changes to your database table. You can have as many migrations as needed as your requirements change over time. For example, one migration might create a new table, one might delete a column, and one might add 2 new columns.

An "up" migration moves the state of the database from its current schema to the schema that you want. So, to get a "blank" database to the state it needs to be ready to run your application, you run all the "up" migrations.

If something breaks, you can run one of the "down" migrations to revert the database to a previous state. "Down" migrations are also used if you need to reset a local testing database to a known state.

1. Install Goose.
    ```
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

    Verify installation:
    ```
    goose -version
    ```

2. Create a database migration in a new `sql/schema` directory. A "migration" in Goose is just a .sql file with some SQL queries and some special comments. Our first migration should just create a users table. The simplest format for these files is:
    ```
    number_name.sql
    ```

    For example, I created a file in sql/schema called 001_users.sql with the following contents:
    ```
    -- +goose Up
    CREATE TABLE ...

    -- +goose Down
    DROP TABLE users;
    ```

3. Get your connection string to your database.
    ```
    protocol://username:password@host:port/database
    ```

4. Run the up migration
    ```
    goose postgres <connection_string> up
    ```

### SQLC
1. Install. Requires Go 1.21+.
    ```
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```

    Verify installation:
    ```
    sqlc version
    ```

2. Configure SQLC. You'll always run the `sqlc` command from the root of the project. Create a `sqlc.yml` or `.yaml` in the root of the project.

3. Generate queries
    ```
    sqlc generate
    ```
