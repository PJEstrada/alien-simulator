package aliemsim

import (
	"math/rand"
)

func getRandomItem(arr []string) string {
	return arr[rand.Intn(len(arr))]
}
