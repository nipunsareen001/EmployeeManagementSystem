package providers

import "Techiebulter/interview/backend/models"

type DbHelperProvider interface {
	CreateEmployee(employee models.Employee) error
	GetEmployeeById(id int) (models.Employee, error)
	UpdateEmployee(empolyee models.Employee) (models.Employee, error)
	DeleteEmployeeById(id int) error
	GetAllEmployees(page string, limit string) ([]models.Employee, error)
}
