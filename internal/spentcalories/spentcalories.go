package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	dataSlice := strings.Split(data, ",")
	if len(dataSlice) != 3 {
		return 0, "", 0, errors.New("incorrect data format")
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("steps are <=0")
	}

	duration, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("duration is <=0")
	}

	activity := dataSlice[1]

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height
	dist := stepLength * float64(steps)

	return dist / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	calories := 0.0

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		break
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		break
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity, duration.Minutes()/minInH, distance(steps, height),
		meanSpeed(steps, height, duration), calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0.0, errors.New("incorrect steps param")
	}

	if duration <= 0 {
		return 0.0, errors.New("incorrect duration param")
	}

	if weight <= 0 {
		return 0.0, errors.New("incorrect weight param")
	}

	if height <= 0 {
		return 0.0, errors.New("incorrect height param")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0, fmt.Errorf("incorrect walking params")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH * walkingCaloriesCoefficient, nil
}
