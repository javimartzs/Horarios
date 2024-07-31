package main

import (
	"horariosapp/controllers"
	"horariosapp/database"
	"horariosapp/utils"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// Inicamos la base de datos
	database.ConnectDB()

	// Configuramos GIN
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatHour": utils.FormatHour,
		"dateFormat": utils.DateFormat,
		"seq":        utils.Seq,
		"div":        utils.Div,
		"mod":        utils.Mod,
	})
	r.LoadHTMLGlob("templates/*")

	// Rutas para la autenticacion
	r.GET("", controllers.ShowLoginPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.GET("/logout", controllers.Logout)

	// Ruta protegida para administradores (workers)
	r.GET("/admin", controllers.ShowAdminPage)
	r.GET("/admin/workers", controllers.ShowWorkersPage)
	r.POST("/admin/workers/create", controllers.CreateWorker)
	r.POST("/admin/workers/update/:id", controllers.UpdateWorker)
	r.POST("/admin/workers/delete/:id", controllers.DeleteWorker)

	// Ruta protegida para administradores (weeks)
	r.GET("/admin/horarios", controllers.ShowWeeksPage)

	// Iniciamos el servidor
	r.Run(":8080")
	log.Print("Server running on port :8080")

}
