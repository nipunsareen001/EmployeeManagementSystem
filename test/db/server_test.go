package server_test

import (
	"Techiebulter/interview/backend/models"
	"Techiebulter/interview/backend/providers/dbHelperProvider"
	"Techiebulter/interview/backend/providers/dbProvider"
	"Techiebulter/interview/backend/server"
	"Techiebulter/interview/backend/utils"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
)

// MockDBClient is a mock implementation of the database client.
type MockDBClient struct {
	ExecContextFn        func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetEmployeeByIdFn    func(id int) (interface{}, error)
	UpdateEmployeeFn     func(employee models.Employee) (models.Employee, error)
	DeleteEmployeeByIdFn func(id int) error
}

// ExecContext executes the given SQL query with optional arguments.
func (m *MockDBClient) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.ExecContextFn != nil {
		return m.ExecContextFn(ctx, query, args...)
	}
	return nil, errors.New("not implemented")
}

// GetEmployeeById simulates the database operation to retrieve employee details by ID.
func (m *MockDBClient) GetEmployeeById(id int) (interface{}, error) {
	if m.GetEmployeeByIdFn != nil {
		return m.GetEmployeeByIdFn(id)
	}
	return nil, errors.New("not implemented")
}

// UpdateEmployee simulates the database operation to update employee details.
func (m *MockDBClient) UpdateEmployee(employee models.Employee) (models.Employee, error) {
	if m.UpdateEmployeeFn != nil {
		return m.UpdateEmployeeFn(employee)
	}
	return models.Employee{}, errors.New("not implemented")
}

// DeleteEmployeeById simulates the database operation to delete an employee by ID.
func (m *MockDBClient) DeleteEmployeeById(id int) error {
	if m.DeleteEmployeeByIdFn != nil {
		return m.DeleteEmployeeByIdFn(id)
	}
	return errors.New("not implemented")
}

// psql database connection
var pgClient = dbProvider.ConnectDB(utils.GetPGSQLConnectionString())

// dbHelpProvider contains all db related helper functions aka repository layer
var dbHelper = dbHelperProvider.NewDBHelper(pgClient.Client())

// Create an instance of DBHelper with the mock database client
var dh = &server.Server{
	PGClient: pgClient,
	DBHelper: dbHelper,
}

func TestCreateEmployee(t *testing.T) {
	// Mock database client
	mockClient := &MockDBClient{}

	// Test case 1: Successful employee creation
	t.Run("CreateEmployee_Success", func(t *testing.T) {
		// Define the expected SQL queries
		expectedCreateTableQuery := `
            DO $$ BEGIN
                CREATE TABLE IF NOT EXISTS employees (
                    id SERIAL PRIMARY KEY,
                    name VARCHAR(255) NOT NULL,
                    position VARCHAR(255) NOT NULL,
                    salary NUMERIC(10, 2) NOT NULL
                );
            EXCEPTION
                WHEN duplicate_table THEN
                    -- Table already exists, do nothing
                    NULL;
            END $$;
        `
		expectedInsertQuery := `
            INSERT INTO employees (name, position, salary)
            VALUES ($1, $2, $3)
        `

		// Mock ExecContext function to assert the expected SQL queries
		mockClient.ExecContextFn = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			switch query {
			case expectedCreateTableQuery:
				// Assert create table query
				// No need to return a result for DDL queries
			case expectedInsertQuery:
				// Assert insert query
				if len(args) != 3 || args[0] != "Trehan" || args[1] != "Software Engineer" || args[2] != 5000000 {
					t.Errorf("Unexpected arguments for insert query. Got: %v", args)
				}
			default:
				t.Errorf("Unexpected SQL query: %s", query)
			}
			return nil, nil
		}

		// Create an employee
		employee := models.Employee{
			Name:     "Trehan",
			Position: "Software Engineer",
			Salary:   5000000,
		}
		err := dh.DBHelper.CreateEmployee(employee)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	// Test case 2: Error handling
	t.Run("CreateEmployee_Error", func(t *testing.T) {
		// Mock ExecContext function to return an error
		mockClient.ExecContextFn = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return nil, errors.New("mock error")
		}

		// Create an employee
		employee := models.Employee{
			Name:     "Trehan",
			Position: "Software Engineer",
			Salary:   5000000,
		}
		err := dh.DBHelper.CreateEmployee(employee)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})
}

func TestGetEmployeeById(t *testing.T) {
	// Mock fiber.Ctx object
	ctx := &fiber.Ctx{}

	// Mock DBHelper
	mockDBHelper := &MockDBClient{}

	// Test case 1: Successful retrieval of employee details
	t.Run("GetEmployeeById_Success", func(t *testing.T) {
		// Define the expected employee details
		expectedEmployeeDetails := map[string]interface{}{
			"id":       1,
			"name":     "Trehan",
			"position": "Software Engineer",
			"salary":   5000000,
		}

		// Mock GetEmployeeByIdFn to return the expected employee details
		mockDBHelper.GetEmployeeByIdFn = func(id int) (interface{}, error) {
			return expectedEmployeeDetails, nil
		}

		// Call the GetEmployeeById method
		_, err := dh.DBHelper.GetEmployeeById(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Assert the response status code and body
		assert.JSONEq(t, `{"status":"success","employeeDetails":{"id":1,"name":"Trehan","position":"Software Engineer","salary":5000000}}`, ctx.Body())
	})

	// Test case 2: Error handling
	t.Run("GetEmployeeById_Error", func(t *testing.T) {
		// Mock GetEmployeeByIdFn to return an error
		mockDBHelper.GetEmployeeByIdFn = func(id int) (interface{}, error) {
			return nil, errors.New("mock error")
		}

		// Call the GetEmployeeById method
		_, err := dh.DBHelper.GetEmployeeById(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Assert the response status code and body

		assert.JSONEq(t, `{"status":"success","data":{"error":"mock error"}}`, ctx.Body())
	})
}

func TestUpdateEmployee(t *testing.T) {
	// Mock fiber.Ctx object
	ctx := &fiber.Ctx{}

	// Mock DBHelper
	mockDBHelper := &MockDBClient{}

	// Test case 1: Successful update of employee details
	t.Run("UpdateEmployee_Success", func(t *testing.T) {
		// Define the expected updated employee details
		expectedUpdatedEmployeeDetails := models.Employee{
			ID:       1,
			Name:     "Trehan",
			Position: "Senior Software Engineer",
			Salary:   60000,
		}

		// Mock UpdateEmployee function to return the expected updated employee details
		mockDBHelper.UpdateEmployeeFn = func(employee models.Employee) (models.Employee, error) {
			return expectedUpdatedEmployeeDetails, nil
		}

		// Set the request body
		// ctx.Add(`{"id":1,"name":"Trehan","position":"Senior Software Engineer","salary":60000}`)

		// Call the UpdateEmployee method
		_, err := dh.DBHelper.UpdateEmployee(expectedUpdatedEmployeeDetails)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Assert the response body
		assert.JSONEq(t, `{"status":"success","updatedEmployeeDetails":{"id":1,"name":"Trehan","position":"Senior Software Engineer","salary":60000}}`, ctx.Body())
	})

	// Test case 2: Error handling
	t.Run("UpdateEmployee_Error", func(t *testing.T) {
		// Mock UpdateEmployee function to return an error
		mockDBHelper.UpdateEmployeeFn = func(employee models.Employee) (models.Employee, error) {
			return models.Employee{}, errors.New("mock error")
		}

		// Set the request body
		// ctx.Add(`{"id":1,"name":"Trehan","position":"Senior Software Engineer","salary":60000}`)

		// Call the UpdateEmployee method
		_, err := dh.DBHelper.UpdateEmployee(models.Employee{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Assert the response body
		assert.JSONEq(t, `{"status":"fail","data":{"error":"mock error"}}`, ctx.Body())
	})
}

func TestDeleteEmployee(t *testing.T) {
	// Mock fiber.Ctx object
	ctx := &fiber.Ctx{}

	// Mock DBHelper
	mockDBHelper := &MockDBClient{}

	// Test case 1: Successful deletion of employee
	t.Run("DeleteEmployee_Success", func(t *testing.T) {
		// Mock DeleteEmployeeById function to return nil (no error)
		mockDBHelper.DeleteEmployeeByIdFn = func(id int) error {
			return nil
		}

		// // Set the request parameter
		// ctx.Add("id", "1")

		// Call the DeleteEmployee method
		err := dh.DBHelper.DeleteEmployeeById(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		assert.JSONEq(t, `{"status":"success"}`, ctx.Body())
	})

	// Test case 2: Error handling
	t.Run("DeleteEmployee_Error", func(t *testing.T) {
		// Mock DeleteEmployeeById function to return an error
		mockDBHelper.DeleteEmployeeByIdFn = func(id int) error {
			return errors.New("mock error")
		}

		// Set the request parameter
		ctx.Params("id", "1")

		// Call the DeleteEmployee method
		err := dh.DBHelper.DeleteEmployeeById(1)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Assert the response body
		assert.JSONEq(t, `{"status":"fail","data":{"error":"mock error"}}`, ctx.Body())
	})
}
