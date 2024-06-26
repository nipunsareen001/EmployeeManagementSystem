# Employee Management System 

This application provides RESTful APIs for managing employees. It allows you to perform CRUD operations (Create, Read, Update, Delete) on employee records stored in a database.

## Endpoints

### 1. Create Employee

- **URL:** `/api/CreateEmpolyee`
- **Method:** `POST`
- **Description:** Create a new employee record.
- **Request Body:** JSON object containing employee details (name, position, salary).
- **Response:** JSON object with status and message indicating success or failure.

### 2. Get Employee by ID

- **URL:** `/api/GetEmployeeById/:id`
- **Method:** `GET`
- **Description:** Retrieve employee details by employee ID.
- **Query Parameters:** `id` (employee ID)
- **Response:** JSON object with status, employee details, or error message.

### 3. Update Employee

- **URL:** `/api/UpdateEmployee`
- **Method:** `PUT`
- **Description:** Update an existing employee record.
- **Request Body:** JSON object containing updated employee details (Id, name, position, salary).
- **Response:** JSON object with status, updated employee details, or error message.

### 4. Delete Employee

- **URL:** `/api/DeleteEmployee/:id`
- **Method:** `DELETE`
- **Description:** Delete an existing employee record.
- **Query Parameters:** `id` (employee ID)
- **Response:** JSON object with status indicating success or failure.

### 5. List Employees with Pagination

- **URL:** `/api/GetAllEmployees/:page/:limit`
- **Method:** `GET`
- **Description:** Retrieve a list of employees with support for pagination.
- **Query Parameters:** `page` (page number), `limit` (number of records per page)
- **Response:** JSON object with status, list of employees for the requested page, or error message.



## Dependencies

- **Fiber:** Fast and Expressive Go web framework
- **PostgreSQL:** Database for storing employee records
- **GORM:** ORM library for database operations

## Setup

1. Clone the repository: `git clone <repository-url>`
2. Install dependencies: `go mod tidy`
3. Set up PostgreSQL database and update the connection details in `config.go`
4. Build and run the application: `go run main.go`

## Testing

- Unit tests are provided for each CURD operation.
 
 
