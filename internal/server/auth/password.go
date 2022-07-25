package auth

import (
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Returns password from environments if environment variable named by envKey is present.
// Otherwise, returns random-generated password.
func MaybeGenerateRandomPassword(envKey string) string {
	if envKey != "" {
		envPassword, present := os.LookupEnv(envKey)
		if present && len(envPassword) > 0 {
			return envPassword
		}
	}

	return generateRandomPassword()
}

const randomGeneratedPasswordLength = 8
const randomGeneratedPasswordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIZKLMNOPQSTUVWXYZ0123456789"

func generateRandomPassword() string {
	return generateRandomString(randomGeneratedPasswordCharset, randomGeneratedPasswordLength)
}

func generateRandomString(charset string, length int) string {
	randomString := make([]byte, length)
	lenCharset := len(charset)

	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(lenCharset)]
	}

	return string(randomString)
}
