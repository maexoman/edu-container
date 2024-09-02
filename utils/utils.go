package utils

import (
	"math/rand"
	"time"
)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewContainerId() string {
	rand.Seed(time.Now().Unix())

	result := make([]rune, 10)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
