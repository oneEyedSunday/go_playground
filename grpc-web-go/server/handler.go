package server

import (
	"context"
	"github.com/satori/go.uuid"
	"log"
)

// grpc server
// Server is the handler implementation
type Server struct {
	Todos []*TodoItem
}

func (s *Server) AddTodo(ctx context.Context, newTodo *AddTodoRequest) (*TodoItem, error) {
	log.Printf("Received new task %s\n", newTodo.Task)

	todoItem := &TodoItem{
		Id:   uuid.NewV1().String(),
		Task: newTodo.Task,
	}

	s.Todos = append(s.Todos, todoItem)

	return todoItem, nil
}

func (s *Server) GetTodos(ctx context.Context, todoRequest *GetTodoRequest) (*TodoResponse, error) {
	log.Printf("Get todos")

	return &TodoResponse{Todos: s.Todos}, nil
}

func (s *Server) DeleteTodo(ctx context.Context, deleteRequest *DeleteTodoRequest) (*DeleteTodoResponse, error) {
	log.Printf("Deleting todo %s\n", deleteRequest.Id)

	updatedTodos := make([]*TodoItem, len(s.Todos)-1)

	for index, todo := range s.Todos {
		if todo.Id == deleteRequest.Id {
			updatedTodos = append(s.Todos[:index], s.Todos[index+1:]...)
			break
		}
	}

	s.Todos = updatedTodos

	return &DeleteTodoResponse{Message: "Deleted Successfully"}, nil

}
