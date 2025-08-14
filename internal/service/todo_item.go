package service

import (
	"fmt"
	"todo-app/internal/repository"
	"todo-app/internal/repository/mysql/models"

	"github.com/sirupsen/logrus"
)

type TodoItemsService struct {
	repo     repository.TodoItems
	listRepo repository.TodoList
}

func NewTodoItemsService(repo repository.TodoItems, lR repository.TodoList) *TodoItemsService {
	return &TodoItemsService{repo: repo, listRepo: lR}
}

func (s *TodoItemsService) Create(userId, listId int, item models.TodoItems) (int, error) {
	logrus.Debug("todo_item/Create. user_id: ", userId, " list_id: ", listId)

	if s == nil || s.listRepo == nil {
		logrus.Error("todo items service is not initialized")
		return 0, fmt.Errorf("todo items service is not initialized")
	}
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		logrus.Fatal("error get item by id. ", err.Error())
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemsService) GetAllItemsList(userId, listId int) ([]models.TodoItems, error) {
	return s.repo.GetAllItemsList(userId, listId)
}

func (s *TodoItemsService) GetAllItem() ([]models.TodoItems, error) {
	return s.repo.GetAllItem()
}

func (s *TodoItemsService) GetById(userId, itemId int) (models.TodoItems, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemsService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemsService) Update(userId, listId int, input models.UpdadeItemInput) error {
	logrus.Debug("todo_item/Update. user id:", userId, "list id:", listId)

	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
