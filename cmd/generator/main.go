package main

import (
	"fmt"
	"math/rand"
)

func RandomString(n int) string {
	letters := []rune("0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func main() {
	for {
		fmt.Println(RandomString(9))
	}
}
