package main

import "log"
import "time"
import "runtime"
import "github.com/gin-gonic/gin"
import _ "github.com/hyperworks/mysql"

var ms = &runtime.MemStats{}

func printStats(logger *log.Logger) {
	for {
		runtime.GC()
		runtime.ReadMemStats(ms)
		logger.Printf("MemStats:\n"+
			"   alloc: %d\n"+
			" t alloc: %d\n"+
			"     sys: %d\n",
			ms.Alloc,
			ms.TotalAlloc,
			ms.Sys,
		)

		time.Sleep(1 * time.Second)
	}
}

func main() {
	srv, e := DefaultServices()
	if e != nil {
		panic(e)
	}

	srv.Logger.Println("starting up...")
	go printStats(srv.Logger)

	r := gin.Default()
	r.GET("/todos", srv.Handlers.GetAllTodos)
	r.PATCH("/todos/:id", srv.Handlers.PatchTodo)

	r.Run(":8080")
}
