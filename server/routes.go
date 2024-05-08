package server

import (
	"github.com/gofiber/fiber/v2"
	lg "github.com/gofiber/fiber/v2/middleware/logger"
)

// InjectRoutes function keeps all the fiber router end point for the server
func (srv *Server) InjectRoutes() *fiber.App {
	app := fiber.New()
	app.Use(lg.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:3000",
	// 	AllowHeaders:     "Origin, Content-Type, Accept",
	// 	AllowMethods:     "GET, POST, PATCH, DELETE",
	// 	AllowCredentials: true,
	// }))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("you are on /")
	})

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("you are on /api")
	})

	api.Get("/healthchecker", srv.HealthCheck)

	api.Post("/CreateEmpolyee", srv.CreateEmployee)
	api.Get("/GetEmployeeById/:id", srv.GetEmployeeById)
	api.Put("/UpdateEmployee", srv.UpdateEmployee)
	api.Delete("/DeleteEmployee/:id", srv.DeleteEmployee)

	api.Get("/GetAllEmployees/:page/:limit", srv.GetAllEmployees)

	return app
}
