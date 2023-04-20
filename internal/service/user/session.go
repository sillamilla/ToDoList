package user

import (
	"crypto/rand"
	"encoding/base64"
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
