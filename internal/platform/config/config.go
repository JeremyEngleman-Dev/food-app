package config

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type SecretKeys struct {
	EmailEncryptionKey cipher.Block
	EmailHashKey       []byte
	TokenHashKey       []byte
}

type Config struct {
	DBString string
	Secrets  SecretKeys
}

func LoadConfig() Config {
	logger := slog.Default()

	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using environment variables")
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
		logger.Error("Missing value: EMAIL_ENCRYPTION_KEY")
		os.Exit(1)
	}
	emailEncryptionKey, err := GetCipher(rawEmailEncryptionKey)
	if err != nil {
		logger.Error("Failed to get cipher of email encryption key", "error", err.Error())
		os.Exit(1)
	}

	emailHashKey := os.Getenv("EMAIL_HASH_KEY")

	tokenHashKey := os.Getenv("TOKEN_HASH_KEY")

	// Return Config
	return Config{
		DBString: dbString,
		Secrets: SecretKeys{
			EmailEncryptionKey: emailEncryptionKey,
			EmailHashKey:       []byte(emailHashKey),
			TokenHashKey:       []byte(tokenHashKey),
		},
	}
}

func GetCipher(data string) (cipher.Block, error) {
	dataByte := []byte(data)

	if len(dataByte) != 16 && len(dataByte) != 24 && len(dataByte) != 32 {
		return nil, fmt.Errorf("Invalid key size: %d", len(dataByte))
	}

	block, err := aes.NewCipher([]byte(data))
	if err != nil {
		return nil, err
	}

	return block, nil
}
