package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataSlice := strings.Split(data, ",")

	if len(dataSlice) != 3 {
		return 0, "", 0, fmt.Errorf("incorrect amount of data params, want: 3, but got: %d", len(dataSlice))
	}

	stepsStr := dataSlice[0]
	activityType := dataSlice[1]
	durationStr := dataSlice[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error durnig convert string to int: %w", err)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("error during parse time: %w", err)
	}

	if steps <= 0 || duration <= 0 {
		return 0, "", 0, fmt.Errorf("steps or duration should be more than 0, got steps: %d, got duration: %d", steps, duration)
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distanceInKm := (float64(steps) * stepLength) / mInKm

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	durationInHours := duration.Hours()

	return distance / durationInHours
}

// TrainingInfo возвращает информацию о тренировке "0,Ходьба,1h30m"
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("error during parse data: %w", err)
	}

	var commonDistance, speed, calories float64
	durationInHours := duration.Hours()

	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("error during calculate calories: %w", err)
		}

		commonDistance = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("error during calculate calories: %w", err)
		}

		commonDistance = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, durationInHours, commonDistance, speed, calories,
	)

	return result, nil
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("param should be more than 0, got: steps %d, weight %.2f, height %.2f, duration %s",
			steps, weight, height, duration.String(),
		)
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("param should be more than 0, got: steps %d, weight %.2f, height %.2f, duration %s",
			steps, weight, height, duration.String(),
		)
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH * walkingCaloriesCoefficient

	return calories, nil
}
