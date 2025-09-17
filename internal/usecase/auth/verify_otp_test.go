package auth

import (
	"context"
	"testing"
	"time"

	userdomain "dekamond/internal/domain/user"

	"github.com/jackc/pgx/v5"
)

func TestVerifyOTPAndIssueToken(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name           string
		phone          string
		code           string
		setupCache     func(*mockCacheStore)
		setupRepo      func(*mockUserRepositoryWithStorage)
		wantError      bool
		wantErrorMsg   string
		expectUser     bool
	}{
		{
			name:  "valid OTP for existing user",
			phone: "+15551234567",
			code:  "123456",
			setupCache: func(m *mockCacheStore) {
				m.store["otp:+15551234567"] = "123456"
			},
			setupRepo: func(m *mockUserRepositoryWithStorage) {
				m.users = map[string]*userdomain.User{
					"+15551234567": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
				}
			},
			wantError:  false,
			expectUser: true,
		},
		{
			name:  "valid OTP for new user",
			phone: "+15551234568",
			code:  "654321",
			setupCache: func(m *mockCacheStore) {
				m.store["otp:+15551234568"] = "654321"
			},
			setupRepo: func(m *mockUserRepositoryWithStorage) {
				m.users = make(map[string]*userdomain.User)
			},
			wantError:  false,
			expectUser: true,
		},
		{
			name:  "invalid OTP",
			phone: "+15551234567",
			code:  "wrong",
			setupCache: func(m *mockCacheStore) {
				m.store["otp:+15551234567"] = "123456"
			},
			setupRepo: func(m *mockUserRepositoryWithStorage) {
				m.users = make(map[string]*userdomain.User)
			},
			wantError:    true,
			wantErrorMsg: "invalid_or_expired_otp",
			expectUser:   false,
		},
		{
			name:  "expired OTP",
			phone: "+15551234567",
			code:  "123456",
			setupCache: func(m *mockCacheStore) {
			},
			setupRepo: func(m *mockUserRepositoryWithStorage) {
				m.users = make(map[string]*userdomain.User)
			},
			wantError:    true,
			wantErrorMsg: "invalid_or_expired_otp",
			expectUser:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := newMockCacheStore()
			repo := &mockUserRepositoryWithStorage{
				users: make(map[string]*userdomain.User),
			}
			
			tt.setupCache(cache)
			tt.setupRepo(repo)

			auc := &AuthUsecase{
				users:     repo,
				cache:     cache,
				jwtSecret: []byte("test-secret"),
				tokenTTL:  24 * time.Hour,
			}

			token, user, err := auc.VerifyOTPAndIssueToken(ctx, tt.phone, tt.code)

			if tt.wantError {
				if err == nil {
					t.Errorf("VerifyOTPAndIssueToken() expected error, got nil")
					return
				}
				if tt.wantErrorMsg != "" && err.Error() != tt.wantErrorMsg {
					t.Errorf("VerifyOTPAndIssueToken() error = %v, want %s", err, tt.wantErrorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("VerifyOTPAndIssueToken() unexpected error: %v", err)
				return
			}

			if token == "" {
				t.Errorf("VerifyOTPAndIssueToken() token is empty")
			}

			if tt.expectUser && user == nil {
				t.Errorf("VerifyOTPAndIssueToken() expected user, got nil")
			}

			if tt.expectUser && user != nil && user.Phone != tt.phone {
				t.Errorf("VerifyOTPAndIssueToken() user phone = %s, want %s", user.Phone, tt.phone)
			}

			key := otpKey(tt.phone)
			_, err = cache.Get(ctx, key)
			if err == nil {
				t.Errorf("VerifyOTPAndIssueToken() OTP should be deleted after verification")
			}
		})
	}
}

func TestValidateOTP(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name       string
		phone      string
		code       string
		setupCache func(*mockCacheStore)
		wantError  bool
	}{
		{
			name:  "valid OTP",
			phone: "+15551234567",
			code:  "123456",
			setupCache: func(m *mockCacheStore) {
				m.store["otp:+15551234567"] = "123456"
			},
			wantError: false,
		},
		{
			name:  "invalid OTP",
			phone: "+15551234567",
			code:  "wrong",
			setupCache: func(m *mockCacheStore) {
				m.store["otp:+15551234567"] = "123456"
			},
			wantError: true,
		},
		{
			name:  "expired OTP",
			phone: "+15551234567",
			code:  "123456",
			setupCache: func(m *mockCacheStore) {
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := newMockCacheStore()
			tt.setupCache(cache)

			auc := &AuthUsecase{
				cache: cache,
			}

			err := auc.validateOTP(ctx, tt.phone, tt.code)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateOTP() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("validateOTP() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetOrCreateUser(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name      string
		phone     string
		setupRepo func(*mockUserRepositoryWithStorage)
		wantError bool
		expectNew bool
	}{
		{
			name:  "existing user",
			phone: "+15551234567",
			setupRepo: func(m *mockUserRepositoryWithStorage) {
				m.users = map[string]*userdomain.User{
					"+15551234567": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
				}
			},
			wantError:  false,
			expectNew:  false,
		},
		{
			name:  "new user",
			phone: "+15551234568",
			setupRepo: func(m *mockUserRepositoryWithStorage) {
				m.users = make(map[string]*userdomain.User)
			},
			wantError:  false,
			expectNew:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockUserRepositoryWithStorage{
				users: make(map[string]*userdomain.User),
			}
			tt.setupRepo(repo)

			auc := &AuthUsecase{
				users: repo,
			}

			user, err := auc.getOrCreateUser(ctx, tt.phone)

			if tt.wantError {
				if err == nil {
					t.Errorf("getOrCreateUser() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("getOrCreateUser() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("getOrCreateUser() user is nil")
				return
			}

			if user.Phone != tt.phone {
				t.Errorf("getOrCreateUser() user phone = %s, want %s", user.Phone, tt.phone)
			}

			if tt.expectNew {
				if user.ID == "" {
					t.Errorf("getOrCreateUser() new user should have ID")
				}
			}
		})
	}
}

type mockUserRepositoryWithStorage struct {
	users map[string]*userdomain.User
}

func (m *mockUserRepositoryWithStorage) GetByPhone(ctx context.Context, phone string) (*userdomain.User, error) {
	user, exists := m.users[phone]
	if !exists {
		return nil, pgx.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepositoryWithStorage) Create(ctx context.Context, phone string) (*userdomain.User, error) {
	user := &userdomain.User{
		ID:        "new-user-id",
		Phone:     phone,
		CreatedAt: time.Now(),
	}
	m.users[phone] = user
	return user, nil
}

func (m *mockUserRepositoryWithStorage) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (m *mockUserRepositoryWithStorage) List(ctx context.Context, search string, limit, offset int) ([]userdomain.User, int, error) {
	users := make([]userdomain.User, 0)
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, len(users), nil
}
