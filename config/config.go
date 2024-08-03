package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// Variables de inicio de postgres
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	// Clave para firmar los tokens jwt
	JWTSecret string

	// Variables del primer usuario
	AdminUser string
	PassUser  string
	RoleUser  string
	NameUser  string
)

func Init() {
	// Importamos el fichero .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Extraemos las variables del .env a variables de Go
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")

	// Clave para firmar los tokens jwt
	JWTSecret = os.Getenv("JWTSecret")

	// Variables del primer usuario
	AdminUser = os.Getenv("ADMIN_USER")
	PassUser = os.Getenv("ADMIN_PASS")
	RoleUser = os.Getenv("ADMIN_ROLE")
	NameUser = os.Getenv("ADMIN_NAME")
}
