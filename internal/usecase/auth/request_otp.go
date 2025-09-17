package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func (auc *AuthUsecase) RequestOTP(ctx context.Context, phone string) (string, error) {
	otp, err := generateNumericOTP(6)
	if err != nil {
		return "", err
	}

	key := otpKey(phone)
	if err := auc.cache.Set(ctx, key, otp, 2*time.Minute); err != nil {
		return "", err
	}
	return otp, nil
}

func generateNumericOTP(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("otp length must be > 0")
	}
	const base = 10
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		d, err := rand.Int(rand.Reader, big.NewInt(base))
		if err != nil {
			return "", err
		}
		b[i] = byte('0' + d.Int64())
	}
	return string(b), nil
}

func otpKey(phone string) string {
	return fmt.Sprintf("otp:%s", phone)
}

