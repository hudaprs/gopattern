# GoPattern

Go pattern with Golang Native using PostgreSQL, JWT & GORM.

## Usage

Run the server in project root

```
    go run main.go
```

Change _.env.example_ to _.env_

| KEY         | Value          |
| ----------- | -------------- |
| DB_HOST     | 127.0.0.1      |
| DB_PORT     | 5432           |
| DB_USER     | postgres       |
| DB_NAME     | gopattern      |
| DB_PASSWORD | yourdbpassword |
| SECRET      | secretJWT      |

## List of users

| Email                 | Password | Role         |
| --------------------- | -------- | ------------ |
| highadmin@gmail.com   | password | High Admin   |
| normaladmin@gmail.com | password | Normal Admin |

## List of Endpoints

List of endpoints for this starter

### Public Routes

| URL                          | Method | Description                  |
| ---------------------------- | ------ | ---------------------------- |
| /api/register                | POST   | Register a new user          |
| /api/login                   | POST   | Logging a user               |
| /api/forgot-password         | POST   | Forgot password user         |
| /api/change-password/{token} | PATCH  | Change / Reset password user |

### High Admin Routes

Only High admin can access this & need token to access this

| URL                | Method | Description         |
| ------------------ | ------ | ------------------- |
| /api/v1/roles      | GET    | Get all roles       |
| /api/v1/roles      | POST   | Creating a new role |
| /api/v1/roles/{id} | GET    | Get one role        |
| /api/v1/roles/{id} | PATCH  | Update role         |
| /api/v1/roles/{id} | DELETE | Delete role         |
| /api/v1/users      | GET    | Get All Users       |

### Protected Routes

Protected routes & need token to access this

| URL                           | Method | Description                          |
| ----------------------------- | ------ | ------------------------------------ |
| /api/v1/users/me              | GET    | Get profile / get authenticated user |
| /api/v1/users/me/upload-image | POST   | Upload image of authenticated user   |
| /api//users/me/delete-image   | GET    | Delete image of authenticated user   |
