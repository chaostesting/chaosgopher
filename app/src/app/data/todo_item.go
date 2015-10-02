package data

import "log"

type TodoItem struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
	Completed   bool   `db:"completed"`
}

func GetAllTodos(db Interface) ([]*TodoItem, error) {
	sql, result := `SELECT * FROM todo_items;`, []*TodoItem{}
	if e := db.Select(&result, sql); e != nil {
		log.Println(e.Error())
		return nil, e
	}

	return result, nil
}

func GetTodo(db Interface, id int64) (*TodoItem, error) {
	sql, result := `SELECT * FROM todo_items WHERE id = ?`, &TodoItem{}
	if e := db.Get(result, sql, id); e != nil {
		return nil, e
	}

	return result, nil
}

func PatchTodo(db Interface, id int64, updates *TodoItem) error {
	sql := `UPDATE todo_items SET completed = ? WHERE id = ?`
	if e := db.Exec(sql, updates.Completed, id); e != nil {
		return e
	}

	return nil
}
