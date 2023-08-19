package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
)

type ITaskUseCase interface {
	GetAllTasks(userID uint) ([]model.TaskResponse, error)
	GetTaskByID(userId uint, taskid uint) (model.TaskResponse, error)
	CreateTask(task model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUseCase struct {
	tr repository.ITaskRepository
	tv validator.ITaskValidator
}

func NewTaskUseCase(tr repository.ITaskRepository, tv validator.ITaskValidator) ITaskUseCase {
	return &taskUseCase{tr, tv}
}

func (tu *taskUseCase) GetAllTasks(userId uint) ([]model.TaskResponse, error) {
	tasks := []model.Task{}
	taskResponses := []model.TaskResponse{}
	if err := tu.tr.GetAllTasks(&tasks, userId); err != nil {
		return nil, err
	}
	for _, task := range tasks {
		t := model.TaskResponse{
			ID:        task.ID,
			Title:     task.Title,
			CreatedAt: task.CreatedAt,
			UpdateAt:  task.UpdateAt,
		}
		taskResponses = append(taskResponses, t)

	}
	return taskResponses, nil
}

func (tu *taskUseCase) GetTaskByID(userId uint, taskId uint) (model.TaskResponse, error) {
	task := model.Task{}
	if err := tu.tr.GetTaskByID(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}
	taskResponse := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdateAt:  task.UpdateAt,
	}
	return taskResponse, nil
}

func (tu *taskUseCase) CreateTask(task model.Task) (model.TaskResponse, error) {
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}
	if err := tu.tr.CreateTask(&task); err != nil {
		return model.TaskResponse{}, err
	}
	taskResponse := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdateAt:  task.UpdateAt,
	}
	return taskResponse, nil
}

func (tu *taskUseCase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}
	if err := tu.tr.UpdateTask(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}
	taskResponse := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdateAt:  task.UpdateAt,
	}
	return taskResponse, nil
}

func (tu *taskUseCase) DeleteTask(userId uint, taskId uint) error {
	if err := tu.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	return nil
}
