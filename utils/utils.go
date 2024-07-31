package utils

func Seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

func DayName(day int) string {
	days := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}
	return days[day-1]
}

func Add(a, b int) int {
	return a + b
}
