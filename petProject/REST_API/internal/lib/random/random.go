package random

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// NewRandomString генерирует случайную строку указанной длины
func NewRandomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // Инициализируем генератор случайных чисел
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))] // Генерируем случайный символ
	}
	return string(result)
}
