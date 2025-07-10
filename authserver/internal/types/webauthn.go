package types

import (
	"encoding/base64"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthnCredential represents a stored WebAuthn credential
type WebAuthnCredential struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	CredentialID []byte    `json:"credential_id" gorm:"uniqueIndex;not null"`
	PublicKey    []byte    `json:"public_key" gorm:"not null"`
	Counter      uint32    `json:"counter" gorm:"not null;default:0"`
	Name         string    `json:"name" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	// Remove or comment out the User field to avoid recursion/conflict
	// User User `json:"user" gorm:"foreignKey:UserID"`
}

// CredentialResponse represents the credential for JSON responses
type CredentialResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	CredentialID string    `json:"credential_id"` // base64url encoded
	PublicKey    string    `json:"public_key"`    // base64 encoded
	Counter      uint32    `json:"counter"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse converts WebAuthnCredential to CredentialResponse
func (c *WebAuthnCredential) ToResponse() CredentialResponse {
	return CredentialResponse{
		ID:           c.ID,
		UserID:       c.UserID,
		CredentialID: base64.RawURLEncoding.EncodeToString(c.CredentialID), // No padding
		PublicKey:    base64.StdEncoding.EncodeToString(c.PublicKey),
		Counter:      c.Counter,
		Name:         c.Name,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

// Convert to webauthn.Credential
func (c *WebAuthnCredential) ToWebAuthnCredential() webauthn.Credential {
	return webauthn.Credential{
		ID:              c.CredentialID,
		PublicKey:       c.PublicKey,
		AttestationType: "none",
		Authenticator: webauthn.Authenticator{
			AAGUID:    make([]byte, 16),
			SignCount: c.Counter,
		},
	}
}

// WebAuthnChallenge represents a temporary challenge for WebAuthn
type WebAuthnChallenge struct {
	UserID    uint      `json:"user_id"`
	Challenge string    `json:"challenge"`
	ExpiresAt time.Time `json:"expires_at"`
}
