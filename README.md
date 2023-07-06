# API Service to perform the basic operations of a Bank

## Technology Stack

- Golang
- Postgres

## Setup Instructions

### Clone App from Git

> Open terminal and execute below commands:

```bash
# Clone the repository on desired path.
$ git clone git@github.com:jass-walia/bank_ops.git

# Move to cloned directory.
$ cd bank_ops
```

## Run via Docker

> *Note*:
You may want to look into `.env` file to review the environment setting.

```bash
# Run via docker-compose.
$ docker-compose up -d
```

### Logging

To view app logs:

```bash
# Optional use of flag `--follow` to tail the logs.
$ docker logs bank --follow
```

## Run Manually

### Install the Pre-requisites on your machine

- [Golang v1.14 or higher](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)

### Setup database

Note: You need to edit the value of `DB_HOST` to *localhost* and `DB_PORT` to the local postgres instance port, in `.env` file.

```bash
# This will create the required database and user for the app.
#
# Give executable permission to the bash script and run.
$ chmod +x setup-db.sh
$ ./setup-db.sh
```

### Build, Test and Run

```bash
# Run all tests.
$ go test ./... -v

# Build the binary.
$ go build -o main .

# Finally, run the binary with the flag to record logs to standard error.
$ ./main -logtostderr
```

## API Documentation

> Note: Params with asterisk sign (*) are mandatory.

### Create Account

`[POST] /api/v1/accounts`

Body Params:

- `*account_holder_name` *string*
- `*account_type` *string*: Allowed values: `saving` or `current`

Request:

```json
{
    "account_holder_name": "David",
    "account_type": "saving"
}
```

Success Response:

```json
{
    "created_at": "2022-02-15T14:03:42.7151735+05:30",
    "uid": "a6d08354-d327-4266-9e94-8f32d69719ab",
    "account_holder_name": "David",
    "account_number": 1,
    "account_type": "saving"
}
```

Error Response:

```json
{"errors": []}
```

Status Codes:

- `201`: Account Created
- `400`: Bad Request
- `422`: Client Validation Error
- `500`: Internal Server Error

### Make Transaction

`[POST] /api/v1/accounts/:uid/transaction`

Path Params:
> `*uid` *string*: Account UID returned on creating an account.

Body Params:

- `*amount` *float*: Upto 2 decimal places.
- `*type` *string*: Type of transaction, allowed values: `credit` or `debit`
- `narration` *string*: Optional comment about transaction.

Request:

```json
{
    "amount": 100,
    "type": "credit",
    "narration": "My first saving"
}
```

Success Response:

```json
{
    "created_at": "2022-02-15T14:08:19.8218828+05:30",
    "uid": "a08efdc4-c304-4706-bb2b-d243c8c61dd7",
    "narration": "My first saving",
    "credit": 100,
    "debit": 0,
    "account": {
        "created_at": "2022-02-15T14:03:42.715173+05:30",
        "uid": "a6d08354-d327-4266-9e94-8f32d69719ab",
        "account_holder_name": "David",
        "account_number": 1,
        "account_type": "saving"
    }
}
```

Error Response:

```json
{"errors": []}
```

Status Codes:

- `201`: Transaction Completed
- `400`: Bad Request
- `404`: Account Not Found
- `422`: Client Validation Error
- `500`: Internal Server Error

### Get Balance

`[GET] /api/v1/accounts/:uid/balance`

Params:
> `*uid` *string*: Account UID returned on creating an account.

Success Response:

```json
{
    "balance": 100
}
```

Error Response:

```json
{"errors": []}
```

Status Codes:

- `200`: Success
- `404`: Account Not Found
- `500`: Internal Server Error

## Trade-offs

### Unit Testing

#### *Problem*: Unit tests are using test database requires database setup

I'm using some test database to execute the unit tests, which violates the unit test definition and gives it a flavour of integration test.

#### *Solution*

One should think of mocking database to eliminate the need of database setup for running tests.

### Integration Testing

I believe integrations tests should be performed in a clean way without leaving any test data behind after tests are finished.
I personally like `dockertest`, a golang package which spins up the database container on-the-fly for testing and removes it after completion. And, I think it is even a good alternative of mocking database.

## Improvements

- Provide http handlers for liveness and readiness signals, will be useful for distributed systems like kubernetes who generally demands these endpoints to understand the pods' health more efficiently.

- Automate test encouraging Continous Integration using GitHib Actions feature or even Jenkins.

## Who do I talk to?

- jaspreet.surmount@gmail.com
