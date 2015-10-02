package main

import "strconv"
import "app/data"
import "github.com/gin-gonic/gin"

type Handlers struct{ *Services }

func (h *Handlers) GetAllTodos(c *gin.Context) {
	todos, e := data.GetAllTodos(h.DB)

	switch {
	case e != nil:
		c.JSON(500, data.ErrorResult(e))
	default:
		c.JSON(200, todos)
	}
}

func (h *Handlers) PatchTodo(c *gin.Context) {
	id, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		c.JSON(400, data.ErrorResult(e))
		return
	}

	todo, e := data.GetTodo(h.DB, id)
	if e != nil {
		c.JSON(500, data.ErrorResult(e))
		return
	} else if todo == nil {
		c.JSON(404, data.NotFoundResult())
		return
	}

	if e := c.BindJSON(todo); e != nil {
		c.JSON(400, data.ErrorResult(e))
		return
	}

	if e := data.PatchTodo(h.DB, id, todo); e != nil {
		c.JSON(500, data.ErrorResult(e))
		return
	}

	c.JSON(200, data.SuccessResult())
}
