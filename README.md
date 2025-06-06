# RSS Blog Aggregator

A simple RSS blog aggregator cli written in Go using Postgres SQL as the database.

## Requirements
1. [Go 1.24 or later](https://go.dev/doc/install)
2. [Postgres SQL 15 or later](https://www.postgresql.org/download/)
3. [Goose](https://github.com/pressly/goose)
4. [Git](https://git-scm.com/downloads)

## Installation
You can install gator in two ways: (1) from source or (2) using Go directly. Afterwards, you need to set up the database and run migrations using goose. 

### (1) Install gator from source
1. Clone the repository:
   ```bash
   git clone https://github.com/donaldnguyen99/gator.git
   ```
2. Change to the project directory:
   ```bash
   cd gator
   ```
3. Install the binary (defaults to `$GOPATH/bin or $HOME/go/bin if $GOPATH is not set`):
   ```bash
   go install
   ```

### (2) Or install gator using Go
1. Install the binary directly:
   ```bash
   go install github.com/donaldnguyen99/gator@latest
   ```

### Migrate the Postgres database
1. Create a connection string like `postgres://username:password@host:port/database`. This connection string will be used to connect to your Postgres database. You can set it as an environment variable, which will also be needed to run gator for the first time:
   ```bash
   export GATOR_POSTGRES_URL="postgres://username:password@host:port/database"
   ```
   - Replace `username`, `password`, `host`, `port`, and `database` with your actual system credentials and Postgres host, port, and database name.
   - Example: `postgres://postgres:postgres@localhost:5432/gator`
   - To create the database, you can use run the `psql` command to log into Postgres and create the database named `gator`. With the example connection string above without a database name, you can run:
     ```bash
        psql postgres://postgres:postgres@localhost:5432
        CREATE DATABASE gator;
        exit
     ```
   
2. Run goose
   ```bash
   goose postgres <connection_string> up
   ```
   This will run the initial migrations to create the necessary tables.

## Usage
- For the first time running gator, make sure the `GATOR_POSTGRES_URL` environment variable is set to your Postgres connection string.
  ```bash
  echo $GATOR_POSTGRES_URL
  ```
  If it is not set, you can set it using `export` from the migration step above.

- You will need to register a user to run gator:
  ```bash
  gator register <your_name>
  ```
  After registering, you will be logged in automatically.

- To login with an existing user:
  ```bash
  gator login <your_name>
  ```

- To list users:
  ```bash
  gator users
  ```

- To add a blog feed, you need to be logged in. Run:
  ```bash
  gator addfeed "Blog Name" <blog_url>
  ```

- To list all blog feeds across all users:
  ```bash
  gator feeds
  ```

- To follow a blog feed from another user, you need to be logged in. Run:
  ```bash
  gator follow <blog_url>
  ```

- To unfollow a blog feed, you need to be logged in. Run:
  ```bash
  gator unfollow <blog_url>
  ```

- To list all blog feeds for the current user:
  ```bash
  gator following
  ```

- To aggregate posts from all followed blogs, run the following as a background process or in a separate terminal. The agg command will require a duration argument to specify how long to aggregate posts for. The duration can be in the format of `1h20m30s` (1 hour, 20 minutes, and 30 seconds) or `1h` (1 hour). Minimum duration is 1 second. The maximum duration is 24 hours.
  ```bash
  gator agg 1h20m30s
  ```

- To list all posts of followed feeds, you need to be logged in. Default number of posts displayed is 2. Run:
  ```bash
  gator browse [num_posts]
  ```