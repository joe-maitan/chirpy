# Chirpy
Chirpy is a social network similar to Twitter. This project is a production style HTTP server in Go, without the use of a framework.

## Highlights
* Basic Go HTTP Server with social media API func.
* PostgreSQL database to manage users and chirps created.
* Use JWTs to handle authentication.
* Use Webhooks to allow premium users to edit posted chirps.

## Overview
This work was done as apart of a course on boot.dev. I took this course to expand my knowledge of Go development and learn how to make HTTP server in Go.

## Usage
Chirpy is technically a monolith project, but a lot of the development was focused on the API. In the future, I would like to make a seperate front end that communicates with this backend.

### Docker
1. Build and run locally.
    ```
    docker build . -t chirpy:dev
    docker run chirpy
    ```

### Locally
#### Install steps
##### PostgreSQL
PostgreSQL is a production-ready, open-source SQL database. You will need to look up the instructions for installing and setting up PostgreSQL on your machine before beginning.

##### Goose
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

##### SQLC
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

##### Running the program
1. Install all the `go.mod` dependencies.
    ```
    go mod tidy
    ```

2. Create your own `.env` file from the `.env.example` file:
    ```
    cp .env.example .env
    ```

    And then fill in the variables.

3. Give permissions to the `run.sh` script and run the program:
    ```
    chmod +x run.sh
    ./run.sh
    ```

## What I learned about
### Middleware
Middleware is a way to wrap a handler with additional functionality. It is a common pattern in web applications that allows us to write DRY code.

### Architecture
Choosing what architecture you want your project to be is an important step of the planning process. While planning this app it helped to understand what the structure of the app should be before getting to lost in the development.

Chirpy is a monolith, but to help with decoupling in the future I worked on developing mainly API functions.

### Hosting my own Database locally
Using PostgreSQL I hosted the database that stored user account information and the chirps they created. I also stored the refresh tokens of each user.

### Authentication
Used JWTs and refresh tokens for user authentciation. The JWTs were to authenticate the user making the requests to the API and the refresh tokens allowed the user to stay logged in for longer periods of time than a single session.

### Authorization
Through the use of the JWTs and refresh tokens we can know authorize a user to do certain things if they are trying to either update their email or password. Or if they would like to delete a chirpy they posted.

### Webhooks
A webhook is an event that is sent to your server by an external service when  something happens. 
