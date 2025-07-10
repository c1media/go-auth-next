package types

import (
	"strconv"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// User represents a user entity with role-based access
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name"`
	Company   string    `json:"company"`
	Role      string    `json:"role" gorm:"not null;default:'user'"` // admin, moderator, user
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// WebAuthn credentials - loaded manually to avoid GORM relationship conflicts
	WebAuthnCredentialsData []WebAuthnCredential `json:"webauthn_credentials" gorm:"-"`
}

// Role constants
const (
	RoleAdmin     = "admin"     // Full system access
	RoleModerator = "moderator" // Content management
	RoleUser      = "user"      // Basic access
)

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Role    string `json:"role"`
}

// AuthResponse represents the response from authentication
type AuthResponse struct {
	User    *User  `json:"user"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// HasPermission checks if the user has the required permission based on role
func (u *User) HasPermission(permission string) bool {
	switch u.Role {
	case RoleAdmin:
		return true // Admin has all permissions
	case RoleModerator:
		return permission == "read" || permission == "write" || permission == "moderate"
	case RoleUser:
		return permission == "read"
	default:
		return false
	}
}

// IsAdmin checks if user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsModerator checks if user is a moderator or admin
func (u *User) IsModerator() bool {
	return u.Role == RoleModerator || u.Role == RoleAdmin
}

// ValidateRole checks if a role is valid
func ValidateRole(role string) bool {
	return role == RoleAdmin || role == RoleModerator || role == RoleUser
}

// WebAuthn interface implementation
func (u *User) WebAuthnID() []byte {
	return []byte(strconv.Itoa(int(u.ID)))
}

func (u *User) WebAuthnName() string {
	return u.Email
}

func (u *User) WebAuthnDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}

func (u *User) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns the user's WebAuthn credentials as []webauthn.Credential
// This method is required by the go-webauthn library
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	credentials := make([]webauthn.Credential, len(u.WebAuthnCredentialsData))
	for i, cred := range u.WebAuthnCredentialsData {
		credentials[i] = cred.ToWebAuthnCredential()
	}
	return credentials
}

func (u *User) WebAuthnCredentialExcludeList() []protocol.CredentialDescriptor {
	excludeList := make([]protocol.CredentialDescriptor, len(u.WebAuthnCredentialsData))
	for i, cred := range u.WebAuthnCredentialsData {
		excludeList[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: []byte(cred.CredentialID),
		}
	}
	return excludeList
}
