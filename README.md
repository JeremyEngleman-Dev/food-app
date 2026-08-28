# Food App

This is something I started to learn more about creating a backend Go API.

Much of the initial focus was setting up user authentication. I just switched from session based auth to JWT.

---

- [Features](#features)
- [How To Run](#how-to-run)
    - [Docker](#run-docker)
    - [Locally](#run-locally)
- [How To Use](#how-to-use)
    - [Health](#use-health)
    - [Users](#use-users)
    - [Authentication](#use-authentication)
    - [Ingredients](#use-ingredients)

## Features <a id="features"></a>

- User authentication via sessions
- Access levels for endpoints
- Ingredient creation

## How To Run <a id="how-to-run"></a>

Options:
- Docker
- Locally

Copy `.env.example` to `.env` and make any desired changes.

NOTE: EMAIL_ENCRYPTION_KEY must be exactly 16, 24, or 32 characters long

### Docker <a id="run-docker"></a>

Ensure that in `.env`, `DB_HOST` is set to "db"

Start the application and its database:

    `docker compose up --build`

The API will be available at `http://localhost:8080`

Stop the application:

    `docker compose down`

### Locally <a id="run-locally"></a>

Ensure that in `.env`, `DB_HOST` is set to "localhost"

Start the app by running:

    `go run main.go`

The API will be available at `http://localhost:8080`

Stop the app using `ctrl` + `c`

## How To Use <a id="how-to-use"></a>

### Health <a id="use-health"></a>

>**GET** `/health`

Gets the health of the application.

**Access Level**: none

---

### Users <a id="use-users"></a>

>**POST** `/users` 

Creates and returns a new user account.

**Access Level**: user

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| displayName | string | yes | Display name of the user |
| username | string | yes | Account identifier |
| password | string | yes | Password to access the account |
| role | string | yes | Access level, 'user' or 'admin' |


>**GET** `/users/{id}`

Retrieves the specified user by id.

**Access Level**: admin

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| id | int | yes | ID of the user to retrieve |

>**GET** `/users`

Retrieves a list of all users.

**Access Level**: admin

>**DELETE** `/users/{id}`

Deletes the specified user.

**Access Level**: admin

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| id | int | yes | ID of the user to delete |

---

### Authentication <a id="use-authentication"></a>

>**POST** `/auth/login`

Logs in the specified user.

**Access Level**: user

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| email | string | yes | Email of the user to log in |
| password | string | yes | Password of the user to log in |

>**POST** `/auth/logout`

Logs out the specified user.

**Access Level**: user

>**POST** `/auth/refresh`

Refreshes the JWT token and refresh token.

**Access Level**: user

---

### Ingredients <a id="use-ingredients"></a>

>**POST** `/ingredients`

Creates an ingredient for the user.

**Access Level**: user

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| name | string | yes | Name of the ingredient |
| category | int | no | Category ID of the ingredient |
| defaultMeasurementType | int | no | Default measurement type ID of the ingredient | 
| description | string | yes | Description of the ingredient |
| createdBy | int | no | ID of the user who created the ingredient |

>**GET** `/ingredients/{id}`

Retrieves an ingredient by ID.

**Access Level**: user

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| id | int | yes | ID of the ingredient to retrieve |

>**GET** `/ingredients`

Retrieves a list of ingredients.

**Access Level**: user

>**DELETE** `/ingredients/{id}`

Deletes the specified ingredient.

**Access Level**: user

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| id | int | yes | ID of the ingredient to delete |