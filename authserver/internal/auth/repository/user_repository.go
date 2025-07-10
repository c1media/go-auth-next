package repository

import (
	"context"
	"fmt"

	"github.com/simple-auth-roles/internal/types"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *types.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*types.User, error) {
	var user types.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *types.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*types.User, error) {
	var users []*types.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) FindByRole(ctx context.Context, role string) ([]*types.User, error) {
	var users []*types.User
	if err := r.db.WithContext(ctx).Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}
	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&types.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetByID is an alias for FindByID (for WebAuthn service)
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*types.User, error) {
	return r.FindByID(ctx, id)
}

// WebAuthn credential methods
func (r *UserRepository) CreateWebAuthnCredential(ctx context.Context, credential *types.WebAuthnCredential) error {
	if err := r.db.WithContext(ctx).Create(credential).Error; err != nil {
		return fmt.Errorf("failed to create webauthn credential: %w", err)
	}
	return nil
}

func (r *UserRepository) LoadWebAuthnCredentials(ctx context.Context, user *types.User) error {
	// Load credentials manually to avoid GORM relationship conflicts
	var credentials []types.WebAuthnCredential
	if err := r.db.WithContext(ctx).Where("user_id = ?", user.ID).Find(&credentials).Error; err != nil {
		return fmt.Errorf("failed to load webauthn credentials: %w", err)
	}
	user.WebAuthnCredentialsData = credentials
	return nil
}

func (r *UserRepository) UpdateWebAuthnCredentialCounter(ctx context.Context, credentialID []byte, counter uint32) error {
	if err := r.db.WithContext(ctx).Model(&types.WebAuthnCredential{}).
		Where("credential_id = ?", credentialID).
		Update("counter", counter).Error; err != nil {
		return fmt.Errorf("failed to update credential counter: %w", err)
	}
	return nil
}

func (r *UserRepository) CountWebAuthnCredentials(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&types.WebAuthnCredential{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count webauthn credentials: %w", err)
	}
	return count, nil
}

func (r *UserRepository) DeleteWebAuthnCredential(ctx context.Context, userID uint, credentialID []byte) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND credential_id = ?", userID, credentialID).
		Delete(&types.WebAuthnCredential{}).Error; err != nil {
		return fmt.Errorf("failed to delete webauthn credential: %w", err)
	}
	return nil
}

func (r *UserRepository) ListWebAuthnCredentials(ctx context.Context, userID uint) ([]types.WebAuthnCredential, error) {
	var creds []types.WebAuthnCredential
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&creds).Error; err != nil {
		return nil, fmt.Errorf("failed to list webauthn credentials: %w", err)
	}
	return creds, nil
}
