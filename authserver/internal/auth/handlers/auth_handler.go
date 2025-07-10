package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/simple-auth-roles/internal/auth/service"
	"github.com/simple-auth-roles/internal/types"
	"github.com/simple-auth-roles/pkg/clientdetection"
	"github.com/simple-auth-roles/pkg/csrf"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.With("handler", "auth"),
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/send-code", h.SendLoginCode)
		auth.POST("/verify-code", h.VerifyLoginCode)
		auth.POST("/check-user", h.CheckUser)
		auth.POST("/create-user", h.CreateUser)       // Admin only
		auth.PUT("/users/:id/role", h.UpdateUserRole) // Admin only
	}
	webauthn := router.Group("/webauthn")
	{
		webauthn.POST("/begin-registration", h.BeginWebAuthnRegistration)
		webauthn.POST("/finish-registration", h.FinishWebAuthnRegistration)
		webauthn.POST("/begin-login", h.BeginWebAuthnLogin)
		webauthn.POST("/finish-login", h.FinishWebAuthnLogin)
		webauthn.POST("/list-credentials", h.ListWebAuthnCredentials)
		webauthn.POST("/delete-credential", h.DeleteWebAuthnCredential)
	}
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name,omitempty"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// --- WebAuthn Handlers ---
type BeginRegistrationRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

type FinishRegistrationRequest struct {
	UserID   uint        `json:"user_id" binding:"required"`
	Response interface{} `json:"credential" binding:"required"`
}

type BeginLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type FinishLoginRequest struct {
	UserID   uint        `json:"user_id" binding:"required"`
	Response interface{} `json:"assertion" binding:"required"`
}

func (h *AuthHandler) BeginWebAuthnRegistration(c *gin.Context) {
	var req BeginRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Note: Add authentication check for user registration in production
	options, err := h.authService.WebAuthnService().BeginRegistration(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, options.Response)
}

func (h *AuthHandler) FinishWebAuthnRegistration(c *gin.Context) {
	var req FinishRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Parse req.Response to protocol.ParsedCredentialCreationData
	credBytes, err := json.Marshal(req.Response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential data"})
		return
	}

	// Log the credential data to see what's being sent
	h.logger.Info("Received credential data", "userID", req.UserID, "credentialData", string(credBytes))

	parsed, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credBytes))
	if err != nil {
		h.logger.Error("Failed to parse credential", "error", err, "userID", req.UserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse credential"})
		return
	}

	// Log the parsed challenge
	h.logger.Info("Parsed credential challenge", "userID", req.UserID, "challenge", parsed.Response.CollectedClientData.Challenge)

	err = h.authService.WebAuthnService().FinishRegistration(c.Request.Context(), req.UserID, parsed)
	if err != nil {
		h.logger.Error("Failed to finish registration", "error", err, "userID", req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AuthHandler) BeginWebAuthnLogin(c *gin.Context) {
	var req BeginLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	options, err := h.authService.WebAuthnService().BeginLogin(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, options)
}

func (h *AuthHandler) FinishWebAuthnLogin(c *gin.Context) {
	var req FinishLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Parse req.Response to protocol.ParsedCredentialAssertionData
	assertionBytes, err := json.Marshal(req.Response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assertion data"})
		return
	}
	parsed, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(assertionBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse assertion"})
		return
	}
	user, err := h.authService.WebAuthnService().FinishLogin(c.Request.Context(), req.UserID, parsed)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token for the authenticated user
	token, err := h.authService.GenerateJWT(user)
	if err != nil {
		h.logger.Error("Failed to generate JWT token", "error", err, "userID", user.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	// Detect client type
	clientInfo := clientdetection.DetectClient(c)

	// Build response based on client type
	responseData := gin.H{
		"success": true,
		"message": "Login successful",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"role":       user.Role,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt,
		},
		"sessionToken": token,
		"clientType":   string(clientInfo.Type),
	}

	c.JSON(http.StatusOK, responseData)
}

// SendLoginCode sends a magic link code to the user's email
func (h *AuthHandler) SendLoginCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	if err := h.authService.SendLoginCode(c.Request.Context(), req.Email, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Detect client type
	clientInfo := clientdetection.DetectClient(c)

	// Build response based on client type
	responseData := gin.H{
		"success":    true,
		"message":    "Login code sent to your email",
		"clientType": string(clientInfo.Type),
	}

	// Only include CSRF token for clients that require it
	if clientInfo.RequiresCSRF() {
		csrfToken, err := csrf.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}
		responseData["csrfToken"] = csrfToken
	}

	c.JSON(http.StatusOK, responseData)
}

// VerifyLoginCode verifies the magic link code and returns a JWT token
func (h *AuthHandler) VerifyLoginCode(c *gin.Context) {
	// Detect client type
	clientInfo := clientdetection.DetectClient(c)

	// Validate CSRF token only for clients that require it
	if clientInfo.RequiresCSRF() {
		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSRF token required"})
			return
		}

		if !csrf.ValidateToken(csrfToken) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			return
		}
	}

	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := h.authService.VerifyLoginCode(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired code"})
		return
	}

	// Build response based on client type
	responseData := gin.H{
		"success": true,
		"message": "Login successful",
		"user": gin.H{
			"id":         response.User.ID,
			"email":      response.User.Email,
			"name":       response.User.Name,
			"role":       response.User.Role,
			"is_active":  response.User.IsActive,
			"created_at": response.User.CreatedAt,
		},
		"sessionToken": response.Token,
		"clientType":   string(clientInfo.Type),
	}

	// Only include CSRF token for clients that require it
	if clientInfo.RequiresCSRF() {
		newCSRFToken, err := csrf.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
			return
		}
		responseData["csrfToken"] = newCSRFToken
	}

	c.JSON(http.StatusOK, responseData)
}

// CreateUser creates a new user (admin only)
func (h *AuthHandler) CreateUser(c *gin.Context) {
	// Note: Add admin authentication middleware in production

	var req types.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "User created successfully",
	})
}

// UpdateUserRole updates a user's role (admin only)
func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	// Note: Add admin authentication middleware in production

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.UpdateUserRole(c.Request.Context(), uint(userID), req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
	})
}

type ListCredentialsRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

func (h *AuthHandler) ListWebAuthnCredentials(c *gin.Context) {
	var req ListCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	creds, err := h.authService.WebAuthnService().ListCredentials(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"credentials": creds})
}

type DeleteCredentialRequest struct {
	UserID       uint   `json:"user_id" binding:"required"`
	CredentialID string `json:"credential_id" binding:"required"`
}

func (h *AuthHandler) DeleteWebAuthnCredential(c *gin.Context) {
	var req DeleteCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Decode credential_id from base64url to []byte
	credID, err := base64.RawURLEncoding.DecodeString(req.CredentialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credential_id encoding"})
		return
	}
	if err := h.authService.WebAuthnService().DeleteCredential(c.Request.Context(), req.UserID, credID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type CheckUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type CheckUserResponse struct {
	UserExists  bool `json:"user_exists"`
	HasPasskeys bool `json:"has_passkeys"`
	UserID      uint `json:"user_id,omitempty"`
}

func (h *AuthHandler) CheckUser(c *gin.Context) {
	var req CheckUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// Check if user exists
	user, err := h.authService.UserRepository().FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Error("Failed to find user", "error", err, "email", req.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, CheckUserResponse{
			UserExists:  false,
			HasPasskeys: false,
		})
		return
	}

	// Check if user has passkeys
	hasPasskeys, err := h.authService.WebAuthnService().HasWebAuthnCredentials(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Error("Failed to check passkeys", "error", err, "email", req.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check passkeys"})
		return
	}

	c.JSON(http.StatusOK, CheckUserResponse{
		UserExists:  true,
		HasPasskeys: hasPasskeys,
		UserID:      user.ID,
	})
}
