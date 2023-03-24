package main

import (
	"Explore1/controller"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
)

func main() {

	m := martini.Classic()
	// 1 admin
	// 2 users
	m.Group("users", func(r martini.Router) {
		r.Get("", controller.Authenticate(controller.GetAllUsers, 1))
		r.Post("", controller.Authenticate(controller.InsertUser, 2))
		r.Put("/:user_id", controller.UpdateUser)
		r.Delete("/:user_id", controller.DeleteUser)
	})

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"localhost:8181"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowCredentials: true,
	}))
	// bisa juga digunakan seperti dibawah ini (untuk masing-masing endpoint)
	allowCORSHandler := cors.Allow(&cors.Options{
		AllowOrigins: []string{"localhost:8181"},
		AllowMethods: []string{"PUT"},
		AllowHeaders: []string{"Origin"},
	})
	m.Put("userscors/:user_id", allowCORSHandler, controller.UpdateUser)

	m.Run()
	m.RunOnAddr(":8181")
}
