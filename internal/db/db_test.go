package db

import (
	"path/filepath"
	"testing"
)

func TestMigrateAndTodos(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	// Running migrations twice must be a no-op.
	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate (second run): %v", err)
	}

	if err := CreateTodo(database, "write tests"); err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}
	todos, err := ListTodos(database)
	if err != nil {
		t.Fatalf("ListTodos: %v", err)
	}
	if len(todos) != 1 || todos[0].Title != "write tests" || todos[0].Done {
		t.Fatalf("unexpected todos: %+v", todos)
	}

	if err := ToggleTodo(database, todos[0].ID); err != nil {
		t.Fatalf("ToggleTodo: %v", err)
	}
	todos, _ = ListTodos(database)
	if !todos[0].Done {
		t.Fatal("todo should be done after toggle")
	}

	if err := DeleteTodo(database, todos[0].ID); err != nil {
		t.Fatalf("DeleteTodo: %v", err)
	}
	todos, _ = ListTodos(database)
	if len(todos) != 0 {
		t.Fatalf("expected empty list after delete, got %+v", todos)
	}
}
