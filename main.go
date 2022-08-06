package main

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jeri06/majoo/config"
	"github.com/jeri06/majoo/response"
	"github.com/jeri06/majoo/server"
	_ "github.com/joho/godotenv/autoload" //
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

var (
	cfg          *config.Config
	utcLocation  *time.Location
	indexMessage string = "Application is running properly"
)

func init() {
	utcLocation, _ = time.LoadLocation("UTC")
	cfg = config.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	// vld := validator.New()

	db, err := sql.Open(cfg.Mysql.Driver, cfg.Mysql.DSN)
	if err != nil {
		logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logger.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(cfg.Mysql.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.Mysql.MaxIdleConnections)

	router := mux.NewRouter()
	router.HandleFunc("/order", index)

	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	srv := server.NewServer(logger, handler, cfg.Application.Port)
	srv.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	<-sigterm

	srv.Close()

}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
