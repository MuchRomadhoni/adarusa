package app

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"fp_pinjaman_online/config"
	"fp_pinjaman_online/config/cloudinary"
	"fp_pinjaman_online/model/dto"
	"fp_pinjaman_online/pkg/validation"
	"fp_pinjaman_online/router"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initEnv() (dto.ConfigData, error) {
	var configData dto.ConfigData
	if err := godotenv.Load(".env"); err != nil {
		return configData, err
	}

	if port := os.Getenv("PORT"); port != "" {
		configData.AppConfig.Port = port
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbMaxIdle := os.Getenv("MAX_IDLE")
	dbMaxConn := os.Getenv("MAX_CONN")
	dbMaxLifeTme := os.Getenv("MAX_LIFE_TIME")
	logMode := os.Getenv("LOG_MODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" || dbMaxIdle == "" || dbMaxConn == "" || dbMaxLifeTme == "" || logMode == "" {
		return configData, errors.New("DB config is not set")
	}

	var err error
	configData.DbConfig.MaxConn, err = strconv.Atoi(dbMaxConn)
	if err != nil {
		return configData, err
	}

	configData.DbConfig.MaxIdle, err = strconv.Atoi(dbMaxIdle)
	if err != nil {
		return configData, err
	}

	configData.DbConfig.Host = dbHost
	configData.DbConfig.DbPort = dbPort
	configData.DbConfig.User = dbUser
	configData.DbConfig.Pass = dbPass
	configData.DbConfig.Database = dbName
	configData.DbConfig.MaxLifeTime = dbMaxLifeTme
	configData.DbConfig.LogMode, err = strconv.Atoi(logMode)
	if err != nil {
		return configData, err
	}

	return configData, nil
}

func initializeDomainModule(r *gin.Engine, db *sql.DB) {
	apiGroup := r.Group("/api")
	v1Group := apiGroup.Group("/v1")
	// init route
	router.InitRoute(v1Group, db)
}

func RunService() {
	configData, err := initEnv()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msg(fmt.Sprintf("config data %v", configData))

	conn, err := config.ConnectToDB(configData, log.Logger)
	if err != nil {
		log.Error().Msg("RunService.ConnectToDB.err : " + err.Error())
		return
	}

	duration, err := time.ParseDuration(configData.DbConfig.MaxLifeTime)
	if err != nil {
		log.Error().Msg("RunService.duration.err : " + err.Error())
		return
	}

	conn.SetConnMaxLifetime(duration)
	conn.SetMaxIdleConns(configData.DbConfig.MaxIdle)
	conn.SetMaxOpenConns(configData.DbConfig.MaxConn)

	defer func() {
		errClose := conn.Close()
		if errClose != nil {
			log.Error().Msg(errClose.Error())
		}
	}()

	time.Local = time.FixedZone("Asia/Jakarta", 7*60*60)
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"POST", "DELETE", "GET", "OPTIONS", "PUT"},
		AllowHeaders: []string{
			"Origin", "Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           120 * time.Second,
	}))

	log.Logger = log.With().Caller().Logger()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validation.ValidationPassword)
	}

	r.Use(logger.SetLogger(
		logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(os.Stdout).With().Caller().Logger()
		}),
	))

	r.Use(gin.Recovery())

	initializeDomainModule(r, conn)
	cloudinary.InitCloudinary()

	version := "0.0.1"
	log.Info().Msg(fmt.Sprintf("Service Running version %s", version))
	addr := flag.String("port: ", ":"+configData.AppConfig.Port, "Address to listen and serve")
	if err := r.Run(*addr); err != nil {
		log.Error().Msg(err.Error())
		return
	}
}
