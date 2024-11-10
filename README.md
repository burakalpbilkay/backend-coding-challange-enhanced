
# Back-end Coding Challenge Enhanced - Go API

This project is a back-end API service built with **Go**, **PostgreSQL**, and **Redis**. It provides user and action data handling, including caching and rate limiting. The API reads data from a PostgreSQL database and exposes endpoints to retrieve user information, action counts, action probabilities, and referral indexes.

## Project Structure

- **/cmd/app**: Contains the main application entry point.
  - **main.go**: The main file to start the application.
  
- **/internal**: Contains the core application logic.
  - **/constants**: Defines constants.
  - **/handlers**: Handles HTTP route logic.
  - **/services**: Business logic for managing user and action data.
  - **/helpers**: Utility functions, including `JSONError` for sending JSON error responses.
  - **/middleware**: Middleware components (rate limiting).
  - **/models**: Defines the data models used within the application.
  - **/repositories**: Manages data access for specific entities.

- **/migrations**: Contains SQL migration files for initializing the PostgreSQL database and data files used for initial database population.
- **/misc**: Contains `users.csv` and `actions.csv` files to populate the database tables.

- **/unit_test**: Test files for unit testing.

- **Dockerfile**: Dockerfile for building the application image.
- **docker-compose.yml**: Docker Compose file for setting up multi-container services.
- **README.md**: Project documentation, including setup instructions.
- **go.mod**: Go module file for dependency management.


## API Endpoints

### 1. Fetch a User by ID
- **Endpoint**: `/user/{id}`
- **Method**: GET
- **Description**: Retrieves a user based on the user ID.
- **Example**: `curl http://localhost:8080/user/1`

### 2. Get Total Number of Actions for a User
- **Endpoint**: `/user/{id}/actions/count`
- **Method**: GET
- **Description**: Returns the total number of actions a user has done.
- **Example**: `curl http://localhost:8080/user/1/actions/count`

### 3. Get Next Action Probabilities
- **Endpoint**: `/action/{type}/next`
- **Method**: GET
- **Description**: Retrieves the probability of possible next actions based on an action type.
- **Example**: `curl http://localhost:8080/action/REFER_USER/next`

### 4. Get Referral Index for All Users
- **Endpoint**: `/users/referral-index`
- **Method**: GET
- **Description**: Returns the referral index of all users. The referral index is the total number of unique users invited directly and indirectly by each user.
- **Example**: `curl http://localhost:8080/users/referral-index`

## Installation and Setup

### Prerequisites
- Docker and Docker Compose
- Go 1.19 or higher

### Steps to Run

1. Clone the repository:
   `git clone https://github.com/burakalpbilkay/backend-coding-challenge-enhanced.git`
   `cd backend-coding-challenge-enhanced`

2. Start services with Docker Compose:
`docker compose up --build`

3. PostgreSQL and Redis will be set up automatically, and the API will be available at `http://localhost:8080`.

## Valid Action Types and Error Handling

The project defines valid action types as constants in `internal/constants/constants.go`. When making a request to `/action/{type}/next`, the `type` parameter is validated against these constants. If an invalid action type is provided, the API returns a `400 Bad Request` error with a message indicating an invalid action type.

For requests to `/user/{id}/actions/count` and `/user/{id}`, the API performs the following validations on `user ID`:
- If the `user ID` is non-numerical, the API returns a `400 Bad Request` error with a message indicating an invalid user ID.
- If the `user ID` is valid but does not match with any user in the database, the API returns a `404 Not Found` error indicating that the user was not found. 

This validation is designed to improve error handling and ensure meaningful responses for invalid requests.


## Testing the API

You can test the API using **curl** or **Postman**.

Example commands:
```bash
curl http://localhost:8080/user/1
curl http://localhost:8080/user/1/actions/count
curl http://localhost:8080/action/REFER_USER/next
curl http://localhost:8080/users/referral-index

```
## Running Unit Tests

To run the unit tests for this project, use the following command:

```bash
go test ./...

```
## Notes

- **Rate Limiting**: Redis is used for rate limiting to restrict request frequency per user.
- **Caching**: Redis caches frequently accessed data to reduce database load.
- **Database Initialization**: SQL migration files in the `/migrations` and  `/misc` directories are run at startup to set up tables and import data.