package service

import (
	"todo-app/pkg/repository"
	"todo-app/types"
)

type TodoItemsService struct {
	repo     repository.TodoItems
	listRepo repository.TodoList
}

func NewTodoItemsService(repo repository.TodoItems, listRepo repository.TodoList) *TodoItemsService {
	return &TodoItemsService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemsService) Create(userId, listId int, item types.TodoItems) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemsService) GetAllItemsList(userId, listId int) ([]types.TodoItems, error) {
	return s.repo.GetAllItemsList(userId, listId)
}

func (s *TodoItemsService) GetAllItem() ([]types.TodoItems, error) {
	return s.repo.GetAllItem()
}

func (s *TodoItemsService) GetById(userId, itemId int) (types.TodoItems, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemsService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemsService) Update(userId, listId int, input types.UpdadeItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
