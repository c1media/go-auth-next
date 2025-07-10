package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const csrfSecret = "your-csrf-secret-change-in-production" // Note: Move to config in production

// GenerateToken creates a new CSRF token
func GenerateToken() (string, error) {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	randomData := hex.EncodeToString(randomBytes)
	data := fmt.Sprintf("%d.%s", timestamp, randomData)

	h := hmac.New(sha256.New, []byte(csrfSecret))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))

	token := fmt.Sprintf("%s.%s", data, signature)
	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}

// ValidateToken validates a CSRF token
func ValidateToken(token string) bool {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	parts := strings.Split(string(decoded), ".")
	if len(parts) != 3 {
		return false
	}

	timestamp, randomData, signature := parts[0], parts[1], parts[2]

	// Check if token is expired (1 hour)
	tokenTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix()-tokenTime > 3600 {
		return false
	}

	// Verify signature
	data := fmt.Sprintf("%s.%s", timestamp, randomData)
	h := hmac.New(sha256.New, []byte(csrfSecret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
