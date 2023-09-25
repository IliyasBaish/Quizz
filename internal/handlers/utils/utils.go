package handlerUtils

import "math/rand"

func GetCode() int {
	min := 10000
	max := 99999
	return rand.Intn(max-min) + min
}
