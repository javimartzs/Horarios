package database

import (
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

	// Creamos la conexion a la base de datos de Postgres
	dsn := "host=localhost user=postgres password=gagll1i1 dbname=horarios port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Migramos las struct de models como tablas de postgres
	err = DB.AutoMigrate(&models.User{}, &models.Worker{}, &models.Week{}, &models.ScheduleEntry{}, &models.WorkerTotal{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	// Creamos el usuario admin si no existe
	createInitialUsers()

	// Llenamos la tabla de semanas si está vacia
	var count int64
	DB.Model(&models.Week{}).Count(&count)
	if count == 0 {
		populateWeeks()
	}
}

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
			continue // El usuario ya existe
		} else if err != gorm.ErrRecordNotFound {
			log.Fatal("Failed to query user:", err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		u.Password = string(hashedPassword)

		if err := DB.Create(&u).Error; err != nil {
			log.Fatal("Failed to create initial user:", err)
		}

		log.Printf("Initial user created: username=%s, password=%s\n", u.Username, u.Password)
	}
}

// Funcion para generar los datos de la tabla Week
func populateWeeks() {
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

func getWeekStartEnd(year int, week int) (string, string) {
	firstDayOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	firstMonday := firstDayOfYear.AddDate(0, 0, (8-int(firstDayOfYear.Weekday()))%7)
	weekStart := firstMonday.AddDate(0, 0, (week-1)*7)
	weekEnd := weekStart.AddDate(0, 0, 6)
	return weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")
}
