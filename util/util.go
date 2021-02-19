package util

import "github.com/gogoods/x/random"

const RANDCHAR = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateRandString(n int) string {
	nonceStr := ""
	for i := 0; i < n; i++ {
		index := random.Int(len(RANDCHAR))
		nonceStr += string(RANDCHAR[index])
	}
	return nonceStr
}
