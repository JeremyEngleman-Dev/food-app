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
	DBString           string
	EmailEncryptionKey cipher.Block
	EmailHashKey       []byte
}

func LoadConfig() Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}

	// Database
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPW := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	dbString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPW,
		dbName,
		dbSSLMode,
	)

	// Keys
	rawEmailEncryptionKey := os.Getenv("EMAIL_ENCRYPTION_KEY")
	if rawEmailEncryptionKey == "" {
		log.Fatal("Missing value: EMAIL_ENCRYPTION_KEY")
	}
	emailEncryptionKey, err := GetCypher(rawEmailEncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	emailHashKey := os.Getenv("EMAIL_HASH_KEY")

	// Return Config
	return Config{
		DBString:           dbString,
		EmailEncryptionKey: emailEncryptionKey,
		EmailHashKey:       []byte(emailHashKey),
	}
}

func GetCypher(data string) (cipher.Block, error) {
	dataByte := []byte(data)

	if len(dataByte) != 16 && len(dataByte) != 24 && len(dataByte) != 32 {
		return nil, fmt.Errorf("Invalid key size: %t", len(dataByte))
	}

	block, err := aes.NewCipher([]byte(data))
	if err != nil {
		log.Fatal(err)
	}

	return block, nil
}
