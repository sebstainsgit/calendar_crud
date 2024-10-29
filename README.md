# Calendar API

This is a **Calendar API** built using **Go**, **PostgreSQL**, and the **Chi router**. The API supports various user operations, including managing events, authentication, and administrative tasks.

## Prerequisites

- **Go 1.17+** - [Download Go](https://golang.org/doc/install)
- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/)
- **Go-Chi Router** - [Chi GitHub](https://github.com/go-chi/chi)
- Environment configuration using **dotenv** - [dotenv GitHub](https://github.com/joho/godotenv)

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/sebstainsgit/calendar_crud.git

cd calendar
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Set Up Environment Variables
Create a .env file in the root directory and add the following values:
Copy code:
```bash
PORT=8080                   # The port on which the server will run
DB_URL="postgres://user:pass@localhost/dbname?sslmode=disable"  # PostgreSQL connection string
JWT_SECRET="your_jwt_secret" # JWT secret for user authentication
```

### 4. Run Database Migrations
Make sure the database is properly set up with the required tables. You can use a migration tool like goose to set up the schema if required. (I used goose)

### 5. Run The Project
```bash
go run main.go
```

## API Endpoints

### General Endpoints

- **GET `/api/ready`**  
  Health check endpoint. Returns an HTTP 200 status code if the service is running properly.

- **GET `/api/error`**  
  Sample error endpoint for testing error handling.

### Authentication Endpoint

- **POST `/api/login`**  
  Allows users to log in. Requires valid user credentials in the request body.

### User Endpoints

- **GET `/user/events`**  
  Retrieves all events associated with the logged-in user.

- **POST `/user/events`**  
  Creates a new event for the logged-in user. Requires event details in the request body.

- **POST `/user/group_event`**  
  Creates a group event. Requires event details in the request body.

- **GET `/user/refresh`**  
  Generates a new JWT using a valid refresh token.

- **POST `/user/update_event`**  
  Updates an existing event for the logged-in user. Requires event details in the request body.

- **DELETE `/user/delete_event`**  
  Deletes a specified event associated with the logged-in user.

- **DELETE `/user/remove_from_event`**  
  Removes the logged-in user from a specified event.

- **POST `/user/users`**  
  Registers a new user. Requires user information in the request body.

- **POST `/user/update_self`**  
  Updates the logged-in user’s information. Requires updated user details in the request body.

- **DELETE `/user/delete_self`**  
  Deletes the logged-in user's account.

### Admin Endpoints

- **GET `/admin/remove_expired_tokens`**  
  Removes expired refresh tokens from the system.

- **POST `/admin/admins`**  
  Creates a new admin user. Requires admin details in the request body.

- **DELETE `/admin/delete_user`**  
  Deletes a specified user account. Requires the user ID of the target user.

- **GET `/admin/users`**  
  Retrieves a list of all users.

- **POST `/admin/promote`**  
  Promotes a specified user to admin status. Requires the user ID of the target user.
