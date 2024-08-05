package main

import (
	"horariosapp/config"
	"horariosapp/controllers"
	"horariosapp/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Importamos las variables de entorno
	config.Init()
	// Iniciamos la base de datos
	database.ConnectDB()
	// Creamos el router de Gin
	r := gin.Default()
	// Importamos las plantillas HTML
	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "./css")

	// Rutas y controladores del login
	r.GET("", controllers.ShowLoginPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
	r.POST("/logout", controllers.Logout)

	// Rutas del Panel de Administrador (Group)
	admin := r.Group("/admin")
	{
		admin.GET("", controllers.ShowPageAdmin)

		// Rutas de trabajadores del grupo Admin
		admin.GET("/workers", controllers.WorkersPage)
		admin.POST("/workers/create", controllers.CreateWorker)
		admin.POST("/workers/update/:id", controllers.UpdateWorker)
		admin.POST("/workers/delete/:id", controllers.DeleteWorker)

		// Rutas de tiendas del grupo Admin
		admin.GET("/stores", controllers.StoresPage)
		admin.POST("/stores/create", controllers.CreateStore)
		admin.POST("/stores/update/:id", controllers.UpdateStore)
		admin.POST("/stores/delete/:id", controllers.DeleteStore)

		// Rutas de calendario del grupo Admin
		admin.GET("/calendar", controllers.WeeksPage)
		admin.GET("/calendar/:weekID", controllers.ShowWeekPage)
		admin.POST("/calendar/:weekID/save", controllers.SaveSchedule)
	}

	r.Run(":8080")

}
