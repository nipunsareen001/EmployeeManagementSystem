package server

import (
	"Techiebulter/interview/backend/providers"
	"Techiebulter/interview/backend/providers/dbHelperProvider"
	"Techiebulter/interview/backend/providers/dbProvider"
	"Techiebulter/interview/backend/utils"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Server struct {
	PGClient providers.PgClientProvider
	DBHelper providers.DbHelperProvider
	Handler  *fiber.App
}

func SrvInit() *Server {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// psql database connection
	pgClient := dbProvider.ConnectDB(utils.GetPGSQLConnectionString())

	// dbHelpProvider contains all db related helper functions aka repository layer
	dbHelper := dbHelperProvider.NewDBHelper(pgClient.Client())

	return &Server{
		PGClient: pgClient,
		DBHelper: dbHelper,
	}
}

func (srv *Server) Start() {
	addr := ":" + utils.GetFIBERPORTString()
	Handler := srv.InjectRoutes()
	Handler.Handler()

	srv.Handler = Handler

	_ = srv.PGClient.Ping()

	logrus.Info("Server running at PORT ", addr)
	if err := Handler.Listen(addr); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Start %v", err)
		return
	}

}

func (srv *Server) Stop() {
	logrus.Info("closing postgresql...")
	_ = srv.PGClient.Close()

	logrus.Info("closing server...")
	_ = srv.Handler.Shutdown()
}
