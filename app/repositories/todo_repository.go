package repositories

import (
	"app/models"

	"gorm.io/gorm"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
	GetAllTodos(todos *[]models.Todo, userID int) error
	GetTodoByID(todo *models.Todo, id int, userID int) error
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(todo *models.Todo) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db}
}

func (tr *todoRepository) CreateTodo(todo *models.Todo) error {
	if err := tr.db.Create(&todo).Error; err != nil {
		return err
	}

	return nil
}

func (tr *todoRepository) GetAllTodos(todos *[]models.Todo, userID int) error {
	if err := tr.db.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return err
	}

	return nil
}

func (tr *todoRepository) GetTodoByID(todo *models.Todo, id int, userID int) error {
	if err := tr.db.Where("user_id = ?", userID).First(&todo, id).Error; err != nil {
		return err
	}

	return nil
}

func (tr *todoRepository) UpdateTodo(todo *models.Todo) error {
	err := tr.db.Model(&todo).Updates(map[string]interface{}{
		"title":   todo.Title,
		"content": todo.Content,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (tr *todoRepository) DeleteTodo(todo *models.Todo) error {
	if err := tr.db.Delete(&todo).Error; err != nil {
		return err
	}

	return nil
}
