package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
)

func main() {
	dbConn := db.NewDB()
	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()
	userRepository := repository.NewUserRepository(dbConn)
	taskRepository := repository.NewTaskRepository(dbConn)
	userUsecase := usecase.NewUserUseCase(userRepository, userValidator, nil)
	taskUsecase := usecase.NewTaskUseCase(taskRepository, taskValidator)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController, taskController)
	e.Logger.Fatal(e.Start(":8080"))
	db.CloseDB(dbConn)
}
