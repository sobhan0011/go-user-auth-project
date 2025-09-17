package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	userdomain "dekamond/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

func (auc *AuthUsecase) VerifyOTPAndIssueToken(ctx context.Context, phone, code string) (string, *userdomain.User, error) {
	if err := auc.validateOTP(ctx, phone, code); err != nil {
		return "", nil, err
	}
	
	auc.cache.Delete(ctx, fmt.Sprintf("otp:%s", phone))

	user, err := auc.getOrCreateUser(ctx, phone)
	if err != nil {
		return "", nil, err
	}

	token, err := auc.generateJWT(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (auc *AuthUsecase) validateOTP(ctx context.Context, phone, code string) error {
	key := fmt.Sprintf("otp:%s", phone)
	storedOTP, err := auc.cache.Get(ctx, key)
	if err != nil {
		return errors.New("invalid_or_expired_otp")
	}
	if storedOTP != code {
		return errors.New("invalid_or_expired_otp")
	}
	return nil
}

func (auc *AuthUsecase) getOrCreateUser(ctx context.Context, phone string) (*userdomain.User, error) {
	user, err := auc.users.GetByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = auc.users.Create(ctx, phone)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			return user, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (auc *AuthUsecase) generateJWT(user *userdomain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"phone": user.Phone,
		"exp": time.Now().Add(auc.tokenTTL).Unix(),
		"iat": time.Now().Unix(),
	})
	signed, err := token.SignedString(auc.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
