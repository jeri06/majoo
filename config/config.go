package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Application struct {
		Port           string
		Name           string
		AllowedOrigins []string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Mysql struct {
		Driver             string
		Host               string
		Port               string
		Username           string
		Password           string
		Database           string
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
}

func Load() *Config {
	cfg := new(Config)
	cfg.logFormatter()
	cfg.app()
	cfg.mysql()
	return cfg
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}
func (cfg *Config) app() {
	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	rawAllowedOrigins := strings.Trim(os.Getenv("ALLOWED_ORIGINS"), " ")

	allowedOrigins := make([]string, 0)
	if rawAllowedOrigins == "" {
		allowedOrigins = append(allowedOrigins, "*")
	} else {
		allowedOrigins = strings.Split(rawAllowedOrigins, ",")
	}

	cfg.Application.Port = port
	cfg.Application.Name = appName
	cfg.Application.AllowedOrigins = allowedOrigins
}

func (cfg *Config) mysql() {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MYSQL_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MYSQL_MAX_IDLE_CONNECTIONS"), 10, 64)

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	connVal := url.Values{}
	connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	cfg.Mysql.Driver = "mysql"
	cfg.Mysql.Host = host
	cfg.Mysql.Port = port
	cfg.Mysql.Username = username
	cfg.Mysql.Password = password
	cfg.Mysql.Database = database
	cfg.Mysql.DSN = dsn
	cfg.Mysql.MaxOpenConnections = int(maxOpenConnections)
	cfg.Mysql.MaxIdleConnections = int(maxIdleConnections)
}
