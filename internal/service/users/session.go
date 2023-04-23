package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateSession() (string, error) {
	buf := make([]byte, 32)

	// Сгенерируйте случайные байты.
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	// Закодируйте байты в строку base64.
	return base64.StdEncoding.EncodeToString(buf), nil
}

func HashGenerate(password string) string {
	hashObject := sha256.New()
	hashObject.Write([]byte(password))
	hashedBytes := hashObject.Sum(nil)

	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}

func HashVerify(password, hash string) bool {
	hashedString := HashGenerate(password)
	return hashedString == hash
}
