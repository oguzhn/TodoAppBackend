package main

import (
	"fmt"
	"sync"
)

type TodoService interface {
	GetAll() ([]Todo, error)
	Save(todo *Todo) error
}

type MockTodoService struct {
	m      sync.Mutex
	nextId int
	todos  []*Todo
}

func NewMockTodoService() *MockTodoService {
	t := new(MockTodoService)
	t.m.Lock()
	t.todos = make([]*Todo, 0)
	t.nextId = 1
	t.m.Unlock()
	return t
}

func (t *MockTodoService) GetAll() ([]*Todo, error) {
	t.m.Lock()
	defer t.m.Unlock()
	return t.todos, nil
}

func (t *MockTodoService) Save(todo *Todo) error {
	t.m.Lock()
	defer t.m.Unlock()
	if todo.Id == 0 { // Insert
		todo.Id = t.nextId
		t.nextId++
		t.todos = append(t.todos, todo)
		return nil
	}

	// Update existing
	for i, value := range t.todos {
		if value.Id == todo.Id {
			t.todos[i] = todo
			return nil
		}
	}

	return fmt.Errorf("Not Found")
}
