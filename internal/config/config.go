package config

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	EmailKey cipher.Block
}

func LoadConfig() Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}

	rawEmailKey := os.Getenv("EMAIL_ENCRYPTION_KEY")
	if rawEmailKey == "" {
		log.Fatal("Missing value: EMAIL_ENCRYPTION_KEY")
	}
	emailKey, err := GetCypher(rawEmailKey)
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		EmailKey: emailKey,
	}
}

func GetCypher(data string) (cipher.Block, error) {
	dataByte := []byte(data)

	if len(dataByte) != 16 && len(dataByte) != 24 && len(dataByte) != 32 {
		return nil, fmt.Errorf("invalid key size")
	}

	block, err := aes.NewCipher([]byte(data))
	if err != nil {
		log.Fatal(err)
	}

	return block, nil
}
