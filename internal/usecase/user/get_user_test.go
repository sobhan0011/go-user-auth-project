package user

import (
	"context"
	"testing"
	"time"

	userdomain "dekamond/internal/domain/user"

	"github.com/jackc/pgx/v5"
)

type mockUserRepository struct {
	users map[string]*userdomain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*userdomain.User),
	}
}

func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*userdomain.User, error) {
	for _, user := range m.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (m *mockUserRepository) Create(ctx context.Context, phone string) (*userdomain.User, error) {
	user := &userdomain.User{
		ID:        "new-user-id",
		Phone:     phone,
		CreatedAt: time.Now(),
	}
	m.users[phone] = user
	return user, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, pgx.ErrNoRows
}

func (m *mockUserRepository) List(ctx context.Context, phone string, limit, offset int) ([]userdomain.User, int, error) {
	users := make([]userdomain.User, 0)
	for _, user := range m.users {
		if phone == "" || user.Phone == phone {
			users = append(users, *user)
		}
	}
	
	start := offset
	end := offset + limit
	if start > len(users) {
		start = len(users)
	}
	if end > len(users) {
		end = len(users)
	}
	
	paginatedUsers := users[start:end]
	return paginatedUsers, len(users), nil
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		id        string
		setupRepo func(*mockUserRepository)
		wantError bool
		wantUser  *userdomain.User
	}{
		{
			name: "valid user ID",
			id:   "user-1",
			setupRepo: func(m *mockUserRepository) {
				m.users["user-1"] = &userdomain.User{
					ID:        "user-1",
					Phone:     "+15551234567",
					CreatedAt: time.Now(),
				}
			},
			wantError: false,
			wantUser: &userdomain.User{
				ID:        "user-1",
				Phone:     "+15551234567",
				CreatedAt: time.Now(),
			},
		},
		{
			name: "user not found",
			id:   "nonexistent",
			setupRepo: func(m *mockUserRepository) {
				m.users = make(map[string]*userdomain.User)
			},
			wantError: true,
			wantUser:  nil,
		},
		{
			name: "empty ID",
			id:   "",
			setupRepo: func(m *mockUserRepository) {
				m.users = make(map[string]*userdomain.User)
			},
			wantError: true,
			wantUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setupRepo(repo)

			uc := &UserUsecase{users: repo}
			user, err := uc.GetByID(ctx, tt.id)

			if tt.wantError {
				if err == nil {
					t.Errorf("GetByID() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetByID() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Errorf("GetByID() user is nil")
				return
			}

			if user.ID != tt.wantUser.ID {
				t.Errorf("GetByID() user ID = %s, want %s", user.ID, tt.wantUser.ID)
			}

			if user.Phone != tt.wantUser.Phone {
				t.Errorf("GetByID() user phone = %s, want %s", user.Phone, tt.wantUser.Phone)
			}
		})
	}
}
