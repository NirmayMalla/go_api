# User API

A REST API built with Go and Fiber for managing users.

## Features

- Create a user
- Get all users
- Get a user by ID
- Update a user
- Delete a user
- Store date of birth (`dob`)
- Return age calculated dynamically
- Input validation using go-playground/validator
- Structured logging using Uber Zap
- Database access generated using SQLC

## Tech Stack

- Go
- Fiber
- MySQL
- SQLC
- go-playground/validator
- Uber Zap

## Project Structure

```
cmd/server/
config/
db/
├── migrations/
└── sqlc/

internal/
├── handler/
├── logger/
├── middleware/
├── models/
├── repository/
├── routes/
└── service/
```

## API Endpoints

### Create User

**POST** `/users`

Request:

```json
{
  "name": "Nirmay",
  "dob": "2002-12-14"
}
```

Response:

```json
{
  "id": 1,
  "name": "Nirmay",
  "dob": "2002-12-14"
}
```

---

### Get All Users

**GET** `/users`

Response:

```json
[
  {
    "id": 1,
    "name": "Nirmay",
    "dob": "2002-12-14",
    "age": 23
  }
]
```

---

### Get User By ID

**GET** `/users/:id`

---

### Update User

**PUT** `/users/:id`

Request:

```json
{
  "name": "Updated Name",
  "dob": "1999-04-03"
}
```

---

### Delete User

**DELETE** `/users/:id`

Response:

```
204 No Content
```

## Running the Project

### Clone the repository

```bash
git clone <repository-url>
cd go-api-project
```

### Install dependencies

```bash
go mod tidy
```

### Create the database

```sql
CREATE DATABASE user_api;
```

Create the `users` table:

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    dob DATE NOT NULL
);
```

### Generate SQLC code

```bash
sqlc generate
```

### Start the server

```bash
go run ./cmd/server
```

The API runs on:

```
http://localhost:3000
```

## Author

Nirmay Malla
