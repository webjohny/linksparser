package mysql

import (
	"math/rand"
	"strings"
	"time"
)

func ArrayRand(arr []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(arr)
	return strings.Trim(arr[n], " ")
}
