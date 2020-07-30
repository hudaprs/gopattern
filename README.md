# GoPattern

Go pattern with Golang Native using PostgreSQL, JWT & GORM.

## What this starter include?

1. MVC Pattern (View not included)
2. Pagination
3. Role Management
4. User Management
5. Easy to manage your code

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

## Controller
I'm gonna recommend you,to make the controller inside *app/controllers*, it will make your controller more easy to import to use it in routing.

## Model
Everything that gonna interact with the database, you will place your model inside *app/models*.

## Middlewares
For the middleware, it's located in *app/middlewares*. 
The default is contains for setting the type of response header and verifying the JWT.

## Database
You can configure your database setting in *config/database.go*.

## Routing
The routing's located in *routes* folder.

## Pagination

All the core pagination functions located in _app/helpers/pagination.go_

Here's the main code:

```go
package helpers

import (
	"net/http"
	"strconv"
)

/**
@desc PaginationResponse make pagination response way better

@param r *http.Request
@param limit int query string from param for limiting data
*/
func Pagination(r *http.Request, limit int) (int, int) {
	keys := r.URL.Query()

	if keys.Get("page") == "" {
		return 1, 0
	}

	page, _ := strconv.Atoi(keys.Get("page"))
	if page < 1 {
		return 1, 0
	}
	begin := (limit * page) - limit
	return page, begin
}

/**
@desc PaginationResponse make pagination response way better

@param *http.Request
@param page int query string from param for getting data from different page
@param pages int list of total pages (total / limit)
@param limit int query string from param for limiting data
@param total int total data in database
@param data interface{} presenting the data there's gonna be output
*/
func PaginationResponse(r *http.Request, page, pages, limit, total int, data interface{}) interface{} {
	return map[string]interface{}{
		"Links": map[string]interface{}{
			"First": r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(page),
			"Last":  r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(pages),
			"Prev":  r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(page-1),
			"Next":  r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(page+1),
		},
		"Meta": map[string]interface{}{
			"Limit":        limit,
			"Total":        total,
			"TotalPage":    pages,
			"CurrentPage":  page,
			"NextPage":     page + 1,
			"PreviousPage": page - 1,
			"LastPage":     pages,
		},
		"Results": data,
	}
}
```

The code is self-explanatory in there

## Pagination Usage

Here's the example for implementing pagination in the controller.
You just insert the code down here.
The URL only accept **_limit_** and **_name_** param

```go
package controllers

// Paginate the roles
queryParams := r.URL.Query()
limit, _ := strconv.Atoi(queryParams.Get("limit"))
nameParam := queryParams.Get("name")
if limit < 1 {
	limit = 10
}
page, begin := helpers.Pagination(r, limit)
// @info total variable's from counting the roles in the model
pages := total / limit
if (total % limit) != 0 {
	pages++
}

// Return the paginate
roles, err := role.GetRoles(begin, limit, nameParam, config.DB)
if err != nil {
	helpers.Error(w, http.StatusBadRequest, err.Error())
	return
}

// Map the data for better response
mapRoles := helpers.PaginationResponse(r, page, pages, limit, total, roles)
helpers.Success(w, http.StatusOK, "Roles list", mapRoles)
```

## Formatted HTTP Response

I'm making a helper for generating better and consistent http response with JSON. The file's located in _app/helpers/json.go_

```go
package helpers
/**
@desc Success response JSON

@param w http.ResponseWriter
@param statusCode int
@param message string
@param data interface{}
 */
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"Status":  "Success",
		"Message": message,
		"Data": data,
	}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

/**
@desc Error response JSON

@param w http.ResponseWriter
@param statusCode int
@param message string
*/
func Error(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if message == "" {
		message = "Something went wrong"
	}
	response := map[string]interface{}{
		"Status": "Error",
		"Message": message,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
```

## Formatted HTTP Response Usage

You can simply use that http helpers JSON by write the example code below

```go
package somewhere

import (

"gopattern/app/helpers"
"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
    // This is for returning a success response
    someData := map[string]interface{}{"Name": "John Doe"}
    helpers.Success(w, http.StatusOK, "Success!", someData)
    return

    // This is for returning a error response
    explicitError := true
    if explicitError == true {
        helpers.Error(w, http.StatusInternalServerError, "Something went wrong!")
    }
}
```

Author Huda Prasetyo 2020, All Right Reserved.
