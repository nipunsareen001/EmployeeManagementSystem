package dbHelperProvider

import (
	"Techiebulter/interview/backend/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// CreateEmployee creates a new employee record in the database.
func (dh *DBHelper) CreateEmployee(employee models.Employee) error {
	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the SQL query for creating the employees table if it doesn't exist
	createTableQuery := `
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

	// Define the SQL query for inserting values into the employees table
	insertQuery := `
        INSERT INTO employees (name, position, salary)
        VALUES ($1, $2, $3)
    `

	// Execute the create table query to ensure the table exists
	_, err := dh.pgClient.ExecContext(ctx, createTableQuery)
	if err != nil {
		log.Print("CreateEmployee: unable to create employees table:", err)
		return err
	}

	// Execute the insert query to add the new employee
	_, err = dh.pgClient.ExecContext(ctx, insertQuery, employee.Name, employee.Position, employee.Salary)
	if err != nil {
		log.Print("CreateEmployee: unable to insert employee into database:", err)
		return err
	}

	// Employee successfully created
	return nil
}

// GetEmployeeById retrieves an employee from the database by their ID.
func (dh *DBHelper) GetEmployeeById(id int) (models.Employee, error) {
	// Initialize an empty Employee struct to store the result
	var emp models.Employee

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the SQL query to select an employee by ID
	query := `
        SELECT id, name, position, salary
        FROM employees
        WHERE id = $1
    `

	// Execute the SQL query to retrieve the employee by ID
	err := dh.pgClient.QueryRowContext(ctx, query, id).Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no employee with the given ID is found, return a specific error
			return emp, fmt.Errorf("employee with ID %d not found", id)
		}
		// If there's an error other than "no rows", return it
		log.Println("GetEmployeeById: error retrieving employee from database:", err)
		return emp, err
	}

	// Return the retrieved employee and nil error
	return emp, nil
}

// UpdateEmployee selectively updates an employee's details in the database based on non-zero and non-empty fields.
func (dh *DBHelper) UpdateEmployee(employee models.Employee) (models.Employee, error) {
	// Initialize an empty Employee struct to store the updated details
	var updatedEmployee models.Employee

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construct the base SQL query for updating an employee's details
	query := "UPDATE employees SET"
	var args []interface{}

	// Check each field of the employee struct and add non-zero and non-empty fields to the query
	if employee.Name != "" {
		query += " name = $1,"
		args = append(args, employee.Name)
	}
	if employee.Position != "" {
		query += " position = $2,"
		args = append(args, employee.Position)
	}
	if employee.Salary != 0 {
		query += " salary = $3,"
		args = append(args, employee.Salary)
	}

	// Remove the trailing comma from the query
	query = strings.TrimSuffix(query, ",")

	// Add the WHERE clause to update the employee with the given ID
	query += " WHERE id = $4 RETURNING id, name, position, salary"
	args = append(args, employee.ID)

	// Execute the SQL query to update the employee's details and retrieve the updated record
	err := dh.pgClient.QueryRowContext(ctx, query, args...).
		Scan(&updatedEmployee.ID, &updatedEmployee.Name, &updatedEmployee.Position, &updatedEmployee.Salary)
	if err != nil {
		log.Println("UpdateEmployee: error updating employee details in database:", err)
		return updatedEmployee, err
	}

	// Return the updated employee details and nil error
	return updatedEmployee, nil
}

// DeleteEmployeeById deletes an employee from the database by their ID.
func (dh *DBHelper) DeleteEmployeeById(id int) error {
	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the SQL query to delete an employee by ID
	query := `
        DELETE FROM employees
        WHERE id = $1
    `

	// Execute the SQL query to delete the employee by ID
	_, err := dh.pgClient.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("DeleteEmployeeById: error deleting employee from database:", err)
		return err
	}

	// Employee successfully deleted
	return nil
}

// GetAllEmployees retrieves all employees from the database with pagination.
func (dh *DBHelper) GetAllEmployees(page string, limit string) ([]models.Employee, error) {
	// Initialize a slice of Employee structs to store the results
	var employees []models.Employee

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse page and limit parameters to integers
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return nil, fmt.Errorf("invalid page number: %s", page)
	}
	limitNumber, err := strconv.Atoi(limit)
	if err != nil {
		return nil, fmt.Errorf("invalid limit: %s", limit)
	}

	// Calculate the offset based on the page number and limit
	offset := (pageNumber - 1) * limitNumber

	// Define the SQL query to select employees with pagination
	query := `
        SELECT id, name, position, salary
        FROM employees
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	// Execute the SQL query to retrieve employees with pagination
	rows, err := dh.pgClient.QueryContext(ctx, query, limitNumber, offset)
	if err != nil {
		log.Println("GetAllEmployees: error getting results from database:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result rows and scan each employee into the slice
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary)
		if err != nil {
			log.Println("GetAllEmployees: error scanning row:", err)
			return nil, err
		}
		employees = append(employees, emp)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Println("GetAllEmployees: error iterating over rows:", err)
		return nil, err
	}

	// Return the slice of employees and nil error
	return employees, nil
}
