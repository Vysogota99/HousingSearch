package postgresstore

import "golang.org/x/crypto/bcrypt"

// HashPassword - ...
func HashPassword(value string) chan string {
	resultChan := make(chan string)
	go func(ch chan string) {
		result, err := bcrypt.GenerateFromPassword([]byte(value), 8)
		if err != nil {
			ch <- ""
		} else {
			ch <- string(result)
		}
	}(resultChan)

	return resultChan
}

// ComparePasswords - ...
func ComparePasswords(hash, password string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return result == nil
}
