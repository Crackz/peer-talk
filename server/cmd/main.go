package main

import (
	"fmt"
	"peer-talk/config"
	db "peer-talk/db/sqlc"
	"peer-talk/internal/auth"
	"peer-talk/internal/common"
	"peer-talk/internal/hub"
	"peer-talk/internal/middlewares"
	"peer-talk/internal/socket_handler"
	"peer-talk/internal/user"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Run() error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	sqlStore, err := db.NewSQLStore()
	if err != nil {
		fmt.Println("FAILED !!!")
		fmt.Println("Failed to create database client")
		return err
	}

	defaultValidator := common.NewDefaultValidator()

	router := echo.New()
	router.Use(middleware.CORS())

	router.Validator = defaultValidator
	router.HTTPErrorHandler = common.NewHttpErrorHandler().Handler

	v1Router := router.Group("/v1")
	authRouter := router.Group("")
	authRouter.Use(middlewares.JwtAuthMiddleware)

	userStore := user.NewStore(sqlStore)
	userHandler := user.NewHandler(userStore)

	authHandler := auth.NewHandler(userHandler)
	authHandler.RegisterRoutes(v1Router)

	hub := hub.New()
	go hub.Run()

	router.Use()

	wsHandler := socket_handler.NewSocketHandler(hub, defaultValidator)
	wsHandler.RegisterRoute(authRouter)

	return router.Start(":" + config.Env.Port)
}

func main() {
	fmt.Println("Starting...")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
