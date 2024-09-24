Task Manager REST API

This is a Go-built REST API for task managers. The API has concurrency and authentication capabilities and lets users control tasks. After logging in, users can add, read, edit, and delete tasks. They can also mark tasks as completed in parallel.

The coding pattern followed in developing these REST apis are Service Layer Pattern and Middleware Pattern, both commonly used in building RESTful applications.

Contents:
1.Features
2.Prerequisites
3.Installation
4.Database Setup
5.Running the Application
6.API Endpoints
7.Bonus Features
8.Unit Tests

1.Features:

User Authentication: Register and login with JWT-based authentication.
Task Management: Create, read, update, and delete tasks.
Concurrency: Mark multiple tasks as "done" concurrently using Goroutines.
Middleware: Logging of incoming requests and authentication for task routes.
Pagination, Sorting, and Filtering: List tasks with sorting and filtering options (Bonus).

2.Prerequisites:

Go 1.18 or higher
PostgreSQL 12 or higher
Postman or cURL (for testing API endpoints)

3.Installation:
Step 1: Clone the Git repository

Step 2: Initailize the go module and Install Dependencies :
Use command : go mod tidy
The dependencies include:
gorm for ORM
gorilla/mux for routing
golang-jwt for JWT-based authentication

Step 3: Database Setup
Install PostgreSQL.
Create a new PostgreSQL database
Update the database configuration in database/database.go:
Replace User and password with PostgreSql Credentials.

Step 4: Run Database Migrations
Make sure the database schema is created by running GORM auto-migrations.

Step 5: Run the Application
Start the Go server by running:
go run main.go
The server will start on http://localhost:3000.

For running the application start the server with command go run main.go


API Documentation :
All endpoints are prefixed with /api only for task-related routes, where as  authentication routes do not have this prefix.

API Endpoints :

1. Authentication :

1.Register a new user :
POST /register
Request body:
{
  "username": "abhishek",
  "password": "abh262000"
}
Response:"User registered successfully"

2.Login a user :
Login
POST /login
Request body :
{
  "username": "abhishek",
  "password": "abh262000"
}
Response:
{
  "token": "jwt_token"
}

2. Task Management : 
All task-related endpoints require a JWT token to be passed in header in a key named Authorization.

1.Create a Task: 

POST /api/tasks
Request body:
{
  "title": "Task Title",
  "description": "Task Description",
  "status": "todo" 
}
Response:
{
  "id": 1,
  "title": "Task Title",
  "description": "Task Description",
  "status": "todo",
  "user_id": 1,
  "created_at": "2023-09-21T15:14:15Z",
  "updated_at": "2023-09-21T15:14:15Z"
}

2.Get All Tasks for a particular User

GET /api/gettasks
Request body: {
    "page": 1,
    "pageSize": 5,
    "sortBy": "status",
    "sortOrder": "asc"
}

Response:
[
  {
    "id": 1,
    "title": "Task 1",
    "description": "description",
    "status": "todo",
    "user_id": 1
  },
  {
    "id": 2,
    "title": "Task 2",
    "description": "description",
    "status": "in_progress",
    "user_id": 1
  }
]

3.Update a Task

PUT /api/updatetasks/{id}
Request body:
{
  "title": "Updated Title",
  "description": "Updated Description",
  "status": "done"
}
Response:
{
  "id": 1,
  "title": "Updated Title",
  "description": "Updated Description",
  "status": "done"
}

4.Delete a Task:

DELETE /api/tasks/{id}
Response: task deleted 

5.Concurrency:

Mark Multiple Tasks as Done Concurrently.
POST /api/concurrent
Request body:
{
  "task_ids": [1, 2, 3]
}
Response:"All tasks marked as done!"

6.Bonus Features:

Pagination for Task Listing implemented pass the page number and page size(meaning how much amount of data to be shown on the page).

Sorting and Filtering.
You can sort tasks by status or created_at.

7.Test:

Unit test cases are written for the implemented RESTful apis in tests folder categorised as controllers_test.go,middlewares_test.go and services_test.go files.This categorization makes it easier to understand for what functions test cases would have been implemented in a file.

3. Deployment:
 
1.Create a Dockerfile:

dockerfile
FROM golang:1.18-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o task-manager .
CMD ["./task-manager"]

2.Build and run the Dockerfile:

docker build -t task-manager-api .
docker run -p 8080:8080 task-manager-api






# task-manager
