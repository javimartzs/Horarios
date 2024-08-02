package database

import (
	"fmt"
	"horariosapp/config"
	"horariosapp/models"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	// Creamos la conexion a la base de datos de postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBPort)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Migramos las structs de los models como tablas de la base de datos
	err = DB.AutoMigrate(
		&models.User{},
		&models.Worker{},
		&models.Week{},
		&models.ScheduleEntry{},
		&models.WorkerHours{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	// Creamos el usuario admin si no existe
	createInitialUsers()

	// Llenamos la tabla weeks con los datos generados
	var count int64
	DB.Model(&models.Week{}).Count(&count)
	if count == 0 {
		generateWeeks()
	}

}

// Funcion para crear el primer usuario admin
func createInitialUsers() {
	users := []models.User{
		{
			Username: "admin",
			Password: "pass",
			Role:     "admin",
			Name:     "Ana",
		},
	}

	for _, u := range users {
		var user models.User
		if err := DB.First(&user, "username = ?", u.Username).Error; err == nil {
			continue // el usuario ya existe
		} else if err != gorm.ErrRecordNotFound {
			log.Fatal("Failed to query user", err)
		}

		hashPasswoord, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password", err)
		}

		u.Password = string(hashPasswoord)

		if err := DB.Create(&u).Error; err != nil {
			log.Fatal("Failed to create initial users", err)
		}
	}
}

// Funcion para generar los datos de la tabla week de postgres
func generateWeeks() {
	startYear := 2024
	endYear := 2034

	for year := startYear; year <= endYear; year++ {
		for week := 1; week <= 52; week++ {

			startDate, endDate := getWeekStartEnd(year, week)
			weekEntry := models.Week{
				Year:   year,
				Week:   week,
				Start:  startDate,
				End:    endDate,
				WeekID: strconv.Itoa(year) + strconv.Itoa(week),
			}
			DB.Create(&weekEntry)
		}
	}
}

// Funcion que genera primer y ultimo dia de cada semana
func getWeekStartEnd(year int, week int) (string, string) {
	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	firstMonday := firstDayOfYear.AddDate(0, 0, (8-int(firstDayOfYear.Weekday()))%7)
	weekStart := firstMonday.AddDate(0, 0, (week-1)*7)
	weekEnd := weekStart.AddDate(0, 0, 6)
	return weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")
}
