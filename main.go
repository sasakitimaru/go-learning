package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
)

func main() {
	dbConn := db.NewDB()
	userRepository := repository.NewUserRepository(dbConn)
	taskRepository := repository.NewTaskRepository(dbConn)
	userUsecase := usecase.NewUserUseCase(userRepository)
	taskUsecase := usecase.NewTaskUseCase(taskRepository)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController, taskController)
	e.Logger.Fatal(e.Start(":8080"))
	db.CloseDB(dbConn)
}
