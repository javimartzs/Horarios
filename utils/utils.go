package utils

import (
	"fmt"
	"time"
)

// Función para formatear la fecha en español
func DateFormat(dateStr string) string {
	t, _ := time.Parse("2006-01-02", dateStr)
	dayNames := []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
	monthNames := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	return fmt.Sprintf("%s %d de %s", dayNames[t.Weekday()], t.Day(), monthNames[t.Month()-1])
}

// Función para generar una secuencia de horas
func Seq(start, end, step float64) []float64 {
	var seq []float64
	for i := start; i <= end; i += step {
		seq = append(seq, i)
	}
	return seq
}

// Función para convertir float64 a string sin punto decimal si es entero
func FormatHour(hour float64) string {
	if hour == float64(int(hour)) {
		return fmt.Sprintf("%02d:00", int(hour))
	}
	return fmt.Sprintf("%02d:30", int(hour))
}

// Función para dividir
func Div(a, b int) int {
	return a / b
}

// Función para obtener el módulo
func Mod(a, b int) int {
	return a % b
}
