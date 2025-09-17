package auth

import (
	"context"
	"testing"
	"time"

	userdomain "dekamond/internal/domain/user"
)

type mockUserRepository struct{}

func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*userdomain.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Create(ctx context.Context, phone string) (*userdomain.User, error) {
	return &userdomain.User{
		ID:        "test-id",
		Phone:     phone,
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	return nil, nil
}

func (m *mockUserRepository) List(ctx context.Context, search string, limit, offset int) ([]userdomain.User, int, error) {
	return nil, 0, nil
}

type mockCacheStore struct {
	store map[string]string
}

func newMockCacheStore() *mockCacheStore {
	return &mockCacheStore{
		store: make(map[string]string),
	}
}

func (m *mockCacheStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	m.store[key] = value
	return nil
}

func (m *mockCacheStore) Get(ctx context.Context, key string) (string, error) {
	value, exists := m.store[key]
	if !exists {
		return "", userdomain.ErrNotFound
	}
	return value, nil
}

func (m *mockCacheStore) Delete(ctx context.Context, key string) error {
	delete(m.store, key)
	return nil
}

func (m *mockCacheStore) Increment(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *mockCacheStore) SetExpiry(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func TestRequestOTP(t *testing.T) {
	ctx := context.Background()
	auc := &AuthUsecase{
		users:     &mockUserRepository{},
		cache:     newMockCacheStore(),
		jwtSecret: []byte("test-secret"),
		tokenTTL:  24 * time.Hour,
	}

	tests := []struct {
		name      string
		phone     string
		wantError bool
	}{
		{
			name:      "valid phone number",
			phone:     "+15551234567",
			wantError: false,
		},
		{
			name:      "another valid phone number",
			phone:     "+1234567890",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := auc.RequestOTP(ctx, tt.phone)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("RequestOTP() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("RequestOTP() unexpected error: %v", err)
				return
			}

			if len(otp) != 6 {
				t.Errorf("RequestOTP() OTP length = %d, want 6", len(otp))
			}

			for _, char := range otp {
				if char < '0' || char > '9' {
					t.Errorf("RequestOTP() OTP contains non-numeric character: %c", char)
					break
				}
			}

			key := otpKey(tt.phone)
			storedOTP, err := auc.cache.Get(ctx, key)
			if err != nil {
				t.Errorf("RequestOTP() OTP not stored in cache: %v", err)
			}
			if storedOTP != otp {
				t.Errorf("RequestOTP() stored OTP = %s, want %s", storedOTP, otp)
			}
		})
	}
}

func TestGenerateNumericOTP(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "valid length 6",
			length:  6,
			wantErr: false,
		},
		{
			name:    "valid length 4",
			length:  4,
			wantErr: false,
		},
		{
			name:    "invalid length 0",
			length:  0,
			wantErr: true,
		},
		{
			name:    "invalid negative length",
			length:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := generateNumericOTP(tt.length)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("generateNumericOTP() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("generateNumericOTP() unexpected error: %v", err)
				return
			}

			if len(otp) != tt.length {
				t.Errorf("generateNumericOTP() length = %d, want %d", len(otp), tt.length)
			}

			for _, char := range otp {
				if char < '0' || char > '9' {
					t.Errorf("generateNumericOTP() contains non-numeric character: %c", char)
					break
				}
			}
		})
	}
}

func TestOTPKey(t *testing.T) {
	phone := "+15551234567"
	expected := "otp:+15551234567"
	result := otpKey(phone)
	
	if result != expected {
		t.Errorf("otpKey() = %s, want %s", result, expected)
	}
}
