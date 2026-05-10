# ptsd-crud-lab

Simple CRUD lab API for guitars.

Current routes:
- POST /guitar
- GET /guitar/{id}

## What happens on startup

When the app starts, it does this automatically:
1. Reads config from environment variables or .env
2. Runs SQL migrations from migrations/
3. Connects to PostgreSQL
4. Starts HTTP server on PORT (default 8080)

You do not need to create the guitar table manually.

## Prerequisites

- Go installed
- Docker Desktop running
- Bruno installed

## Run the app

From project root, run in PowerShell:

1. Start PostgreSQL in Docker

	docker compose up -d db

2. Create local env file

	Copy-Item .env.example .env -ErrorAction SilentlyContinue

3. Start API

	go run .

Default API base URL is:
http://localhost:8080

## Test with Bruno

Collection is already prepared in:
- bruno/CRUDlab

### Open collection and environment

1. Open Bruno
2. Choose Open Collection
3. Select folder bruno/CRUDlab
4. In the top right, select environment local

Important: if no environment is selected, Bruno will show Invalid URL because {{baseUrl}} cannot be resolved.

### Request flow and expected responses

Run requests in this order.

1. Create Guitar

	File: bruno/CRUDlab/guitar/01-create-guitar.bru

	Expected:
	- Status 201 Created
	- JSON body with generated id

	This request has a post-response script that saves response id into environment variable guitarId.

2. Get Guitar By ID

	File: bruno/CRUDlab/guitar/02-get-guitar-by-id.bru

	Expected:
	- Status 200 OK
	- JSON body for the created guitar

3. Create Guitar Invalid Body

	File: bruno/CRUDlab/guitar/03-create-guitar-invalid-body.bru

	Expected:
	- Status 400 Bad Request
	- Validation message, for example Manufacturer is required

4. Get Guitar Invalid UUID

	File: bruno/CRUDlab/guitar/04-get-guitar-invalid-uuid.bru

	Expected:
	- Status 400 Bad Request
	- Message Invalid UUID

## Quick troubleshooting

1. Invalid URL in Bruno
	- Select local environment in top right

2. connect: connection refused
	- Check that app is running with go run .
	- Check baseUrl is http://localhost:8080

3. Guitar not found on Get Guitar By ID
	- Run Create Guitar first
	- Confirm guitarId was updated in active environment

4. Database connection errors on app start
	- Ensure docker compose up -d db is running
	- Ensure .env exists and DATABASE_URL is valid

## Stop services

1. Stop app with Ctrl+C in terminal where go run . is running
2. Stop database container:

	docker compose down