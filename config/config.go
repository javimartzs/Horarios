package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuramos la base de datos
var (
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	JWTSecret string

	AdminUser string
	AdminPass string
	AdminRole string
	AdminName string
)

// Init carga las variables de entorno
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")

	JWTSecret = os.Getenv("JWT_SECRET")

	AdminUser = os.Getenv("ADMIN_USER")
	AdminPass = os.Getenv("ADMIN_PASS")
	AdminRole = os.Getenv("ADMIN_ROLE")
	AdminName = os.Getenv("ADMIN_NAME")

}
