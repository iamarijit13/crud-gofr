package main

import (
	"employee-crud/handler"
	"employee-crud/store"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	k := gofr.New()

	s := store.New()
	h := handler.New(s)

	k.GET("/employee", h.Find)
	k.GET("/employee/{id}", h.FindByID)
	k.POST("/employee", h.Create)
	k.PUT("/employee/{id}", h.Update)
	k.DELETE("/employee/{id}", h.Delete)

	k.Start()
}
