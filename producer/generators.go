package producer

import "math/rand"

const generator = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getRandomNumber() int {
	return rand.Intn(10000000)
}

func getRandomString(num int) string {
	str := []byte{}
	for range num {
		str = append(str, generator[rand.Intn(len(generator))])
	}
	return string(str)
}

func getUID() string {
	return getRandomString(15)
}
