package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ITaskController interface {
	GetAllTasks(c echo.Context) error
	GetAllTasksById(c echo.Context) error
	CreateTask(c echo.Context) error
	UpdateTask(c echo.Context) error
	DeleteTask(c echo.Context) error
}

type taskController struct {
	tu usecase.ITaskUseCase
}

func NewTaskController(tu usecase.ITaskUseCase) ITaskController {
	return &taskController{tu}
}

func (tc *taskController) GetAllTasks(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))
	tasks, err := tc.tu.GetAllTasks(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tasks)
}

func (tc *taskController) GetAllTasksById(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	task, err := tc.tu.GetTaskByID(userId, uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, task)
}

func (tc *taskController) CreateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint64(claims["user_id"].(float64))
	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	task.UserID = userId
	taskResponse, err := tc.tu.CreateTask(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskResponse)
}

func (tc *taskController) UpdateTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint64(claims["user_id"].(float64))
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	task := model.Task{}
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	task.UserID = userId
	taskResponse, err := tc.tu.UpdateTask(task, uint(userId), uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, taskResponse)
}

func (tc *taskController) DeleteTask(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = tc.tu.DeleteTask(userId, uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Task deleted")
}
