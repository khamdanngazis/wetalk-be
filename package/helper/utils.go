package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CalculateTokenExpiration() time.Time {
	// Set token expiration to 15 minutes
	return time.Now().Add(time.Minute * 10)
}

// GenerateRandomString generates a random string of a specified length.
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}

func GenerateRandomNumber(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}

func GenerateRandomStringCaps(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}

func GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	b := make([]byte, length)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[idx.Int64()]
	}
	return string(b), nil
}

func HashPassword(password string) (string, error) {
	// Generate a hashed representation of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashAndPassword(storedHash string, password string) error {
	// Convert the stored hash from string to []byte
	hashBytes := []byte(storedHash)

	// Compare the stored hash with the password
	return bcrypt.CompareHashAndPassword(hashBytes, []byte(password))
}

func GenerateBcryptSalt() (string, error) {
	// Generate 16 random bytes for the salt
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to base64 to get a string representation
	salt := base64.StdEncoding.EncodeToString(saltBytes)

	// Prefix with the bcrypt identifier and cost factor
	return fmt.Sprintf("$2a$12$%s", salt), nil
}

func Decrypt(cipherText string, key string) (string, error) {
	// Convert the key and cipher text to byte slices
	keyBytes := []byte(key)
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %v", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %v", err)
	}

	// Create a GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM cipher: %v", err)
	}

	// Extract the nonce and the encrypted text
	nonceSize := gcm.NonceSize()
	if len(cipherTextBytes) < nonceSize {
		return "", fmt.Errorf("cipher text too short")
	}
	nonce, cipherTextBytes := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]

	// Decrypt the cipher text
	plainTextBytes, err := gcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting text: %v", err)
	}

	return string(plainTextBytes), nil
}

func Encrypt(plainText string, key string) (string, error) {
	// Convert the key to a byte slice
	keyBytes := []byte(key)
	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %v", err)
	}

	// Create a GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM cipher: %v", err)
	}

	// Generate a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %v", err)
	}

	// Encrypt the plain text
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	// Encode to base64 for easy storage
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func FormatMessageTime(t time.Time) string {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04") // 24-hour format
	} else if t.Year() == yesterday.Year() && t.YearDay() == yesterday.YearDay() {
		return "Yesterday"
	} else {
		return t.Format("2006-01-02")
	}
}
