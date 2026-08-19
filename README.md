# Food App

This is something I started to learn more about creating a backend Go API.

Much of the initial focus was setting up user login/logout with sessions.

Things I have considered:
- I would clearly not allow the willy nilly creation of admin accounts by anyone in a production environment, that would be dumb.
- Better data validation is certainly needed, among other things.

## Features

- User authentication via sessions
- Access levels for endpoints
- Ingredient creation

## How To Run Locally

Requirements:
- Docker
- Docker Compose

Copy `.env.example` and make any desired changes.
*EMAIL_ENCRYPTION_KEY must be exactly 16, 24, or 32 characters long

To start the application and its database:
`docker compose up --build`

The API will be available at `http://localhost:8080`

To stop the application:
`docker compose down`

## How To Use

### Health

>**GET** `/health`

Gets the health of the application.

**Access Level**: none

### Users

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

### Login

>**POST** `/login`

Logs in the specified user.

**Access Level**: user

| Parameter Name | Type | Required | Description |
| --- | --- | --- | --- |
| email | string | yes | Email of the user to log in |
| password | string | yes | Password of the user to log in |

>**POST** `/logout`

Logs out the specified user.

**Access Level**: user

### Ingredients

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