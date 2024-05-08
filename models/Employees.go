package models

import "errors"

type Employee struct {
	ID       int     `json:"ID"`
	Name     string  `json:"Name"`
	Position string  `json:"position"`
	Salary   float64 `json:"Salary"`
}

func (e *Employee) CheckFeilds() error {
	// Check if all fields are present
	if e.Name == "" || e.Position == "" || e.Salary < 1 {
		return errors.New("all fields are mandatory")
	}

	return nil
}

func (e *Employee) CheckId() error {
	// Check if all fields are present
	if e.ID < 1 {
		return errors.New("all fields are mandatory")
	}

	return nil
}
