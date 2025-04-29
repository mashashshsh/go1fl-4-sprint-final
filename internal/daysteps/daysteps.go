package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go1fl-4-sprint-final/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) != 2 {
		return 0, 0, fmt.Errorf("incorrect amount of data params, want: 2, but got: %d", len(dataSlice))
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error durnig convert string to int: %w", err)
	}

	duration, err := time.ParseDuration(dataSlice[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error during parse time: %w", err)
	}

	if steps <= 0 || duration <= 0 {
		return 0, 0, fmt.Errorf("steps or duration should be more than 0, got steps: %d, got duration: %s", steps, duration.String())
	}

	return steps, duration, nil
}

// DayActionInfo возвращает данные о пройденной дистанции в километрах и количество потраченных калорий
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceInM := float64(steps) * stepLength
	distanceInKm := distanceInM / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)

	return result
}
