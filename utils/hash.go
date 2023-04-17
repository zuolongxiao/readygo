package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword using bcrypt to hash a string
func HashPassword(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	byteHash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(byteHash), nil
}

// VerifyPassword verify a bcrypt hashed string
func VerifyPassword(hashedPwd, plainPwd string) bool {
	byteHashed := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	return err == nil
}
