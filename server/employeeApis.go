package server

import (
	"Techiebulter/interview/backend/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) CreateEmployee(c *fiber.Ctx) error {
	var Employee models.Employee

	if err := c.BodyParser(&Employee); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Check if all fields of Employee are present
	if err := Employee.CheckFeilds(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	// Use a channel to communicate errors back from goroutines
	errChan := make(chan error, 1)

	// Start a goroutine to execute the database operation
	go func() {
		errChan <- s.DBHelper.CreateEmployee(Employee)
	}()

	// Wait for the database operation to complete
	err := <-errChan

	if err != nil {
		log.Println("CreateEmployee: error inserting data in the database", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "data": fiber.Map{"error": err}})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (s *Server) GetEmployeeById(c *fiber.Ctx) error {
	id_string := c.Params("id")
	id, err := strconv.Atoi(id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Use a channel to communicate results and errors back from goroutines
	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	// Start a goroutine to execute the database operation
	go func() {
		employeeDetails, err := s.DBHelper.GetEmployeeById(id)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- employeeDetails
	}()

	// Wait for the database operation to complete
	select {
	case employeeDetails := <-resultChan:
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "employeeDetails": employeeDetails})
	case err := <-errChan:
		log.Println("GetEmployeeById: error getting results from DB", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "success", "data": fiber.Map{"error": err}})
	}
}

func (s *Server) UpdateEmployee(c *fiber.Ctx) error {
	var Employee models.Employee

	if err := c.BodyParser(&Employee); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Check if Id of Employee is present
	if err := Employee.CheckId(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": err.Error()})
	}

	// Use a channel to communicate errors and results back from the goroutine
	resultChan := make(chan models.Employee, 1)
	errChan := make(chan error, 1)

	// Start a goroutine to execute the database operation
	go func() {
		updatedEmployeeDetails, err := s.DBHelper.UpdateEmployee(Employee)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- updatedEmployeeDetails
	}()

	// Wait for the database operation to complete
	select {
	case updatedEmployeeDetails := <-resultChan:
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "updatedEmployeeDetails": updatedEmployeeDetails})
	case err := <-errChan:
		log.Println("UpdateEmployee: error updating data in the database", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "data": fiber.Map{"error": err}})
	}
}

func (s *Server) DeleteEmployee(c *fiber.Ctx) error {
	id_string := c.Params("id")
	id, err := strconv.Atoi(id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Use a channel to communicate errors back from the goroutine
	errChan := make(chan error, 1)

	// Start a goroutine to execute the database operation
	go func() {
		errChan <- s.DBHelper.DeleteEmployeeById(id)
	}()

	// Wait for the database operation to complete
	err = <-errChan
	if err != nil {
		log.Println("DeleteEmployeeById: error deleting employee from DB", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "data": fiber.Map{"error": err}})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (s *Server) GetAllEmployees(c *fiber.Ctx) error {
	page_string := c.Params("page")
	limit_string := c.Params("limit")

	// Use a channel to communicate errors and results back from the goroutine
	resultChan := make(chan []models.Employee, 1)
	errChan := make(chan error, 1)

	// Start a goroutine to execute the database operation
	go func() {
		allEmployees, err := s.DBHelper.GetAllEmployees(page_string, limit_string)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- allEmployees
	}()

	// Wait for the database operation to complete
	select {
	case allEmployees := <-resultChan:
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "employees": allEmployees})
	case err := <-errChan:
		log.Println("GetAllEmployees: error getting results from DB", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "data": fiber.Map{"error": err}})
	}
}
