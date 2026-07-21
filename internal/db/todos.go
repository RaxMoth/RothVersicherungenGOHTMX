package db

import "database/sql"

// Todo backs the HTMX demo. Replace with your own models per project.
type Todo struct {
	ID    int64
	Title string
	Done  bool
}

func ListTodos(database *sql.DB) ([]Todo, error) {
	rows, err := database.Query(`SELECT id, title, done FROM todos ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func CreateTodo(database *sql.DB, title string) error {
	_, err := database.Exec(`INSERT INTO todos (title) VALUES (?)`, title)
	return err
}

func ToggleTodo(database *sql.DB, id int64) error {
	_, err := database.Exec(`UPDATE todos SET done = NOT done WHERE id = ?`, id)
	return err
}

func DeleteTodo(database *sql.DB, id int64) error {
	_, err := database.Exec(`DELETE FROM todos WHERE id = ?`, id)
	return err
}
