package main

import "time"
import "strings"
import "log"
import "app/config"
import "app/data"
import "app/client"

var api = &client.APIClient{Host: config.Get().APIEndpoint}

func main() {
	_ = time.Sleep

	for {
		listTodos()
		completeTodo(2)
		completeTodo(4)
		listTodos()
		restartTodo(2)
		restartTodo(4)
		// time.Sleep(1 * time.Second)

		listTodos()
		completeTodo(1)
		completeTodo(3)
		listTodos()
		restartTodo(1)
		restartTodo(3)
		// time.Sleep(2 * time.Second)
	}
}

func listTodos() {
	log.Println("List of TODOs:")
	todos, e := api.GetAllTodos()
	noError(e)

	for _, todo := range todos {
		desc := todo.Description
		if todo.Completed {
			desc = strings.Repeat("-", len(desc))
		}

		log.Printf(" %d. %s", todo.ID, desc)
	}
}

func completeTodo(id int64) {
	log.Println("Complete #", id)
	_, e := api.PatchTodo(id, &data.TodoItem{Completed: true})
	noError(e)
}

func restartTodo(id int64) {
	log.Println("Restart #", id)
	_, e := api.PatchTodo(id, &data.TodoItem{Completed: false})
	noError(e)
}

func noError(e error) {
	if e != nil {
		log.Println("error: ", e.Error())
	}
}
