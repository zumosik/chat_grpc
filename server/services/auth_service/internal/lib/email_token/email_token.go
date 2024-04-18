package email_token

import "math/rand"

func GetRndEmailToken(length int) string {
	// create 6 digit string

	chars := "0123456789"
	var result string

	for i := 0; i < length; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}

	return result
}
