package user

import (
	"context"
	"testing"
	"time"

	userdomain "dekamond/internal/domain/user"
)

func TestList(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		query    ListQuery
		setupRepo func(*mockUserRepository)
		wantPage Page[userdomain.User]
		wantErr  bool
	}{
		{
			name: "list all users",
			query: ListQuery{
				Phone: "",
				Page:  1,
				Limit: 10,
			},
			setupRepo: func(m *mockUserRepository) {
				m.users = map[string]*userdomain.User{
					"user-1": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
					"user-2": {
						ID:        "user-2",
						Phone:     "+15551234568",
						CreatedAt: time.Now(),
					},
				}
			},
			wantPage: Page[userdomain.User]{
				Items: []userdomain.User{
					{ID: "user-1", Phone: "+15551234567"},
					{ID: "user-2", Phone: "+15551234568"},
				},
				Total: 2,
				Page:  1,
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "search by phone",
			query: ListQuery{
				Phone: "+15551234567",
				Page:  1,
				Limit: 10,
			},
			setupRepo: func(m *mockUserRepository) {
				m.users = map[string]*userdomain.User{
					"user-1": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
					"user-2": {
						ID:        "user-2",
						Phone:     "+1234567890",
						CreatedAt: time.Now(),
					},
				}
			},
			wantPage: Page[userdomain.User]{
				Items: []userdomain.User{
					{ID: "user-1", Phone: "+15551234567"},
				},
				Total: 1,
				Page:  1,
				Limit: 10,
			},
			wantErr: false,
		},
		{
			name: "pagination",
			query: ListQuery{
				Phone: "",
				Page: 2,
				Limit: 1,
			},
			setupRepo: func(m *mockUserRepository) {
				m.users = map[string]*userdomain.User{
					"user-1": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
					"user-2": {
						ID:        "user-2",
						Phone:     "+15551234568",
						CreatedAt: time.Now(),
					},
				}
			},
			wantPage: Page[userdomain.User]{
				Items: []userdomain.User{
					{ID: "user-2", Phone: "+15551234568"},
				},
				Total: 2,
				Page:  2,
				Limit:  1,
			},
			wantErr: false,
		},
		{
			name: "default page and size",
			query: ListQuery{
				Phone: "",
				Page: 0,
				Limit: 0,
			},
			setupRepo: func(m *mockUserRepository) {
				m.users = map[string]*userdomain.User{
					"user-1": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
				}
			},
			wantPage: Page[userdomain.User]{
				Items: []userdomain.User{
					{ID: "user-1", Phone: "+15551234567"},
				},
				Total: 1,
				Page:  1,
				Limit: 20,
			},
			wantErr: false,
		},
		{
			name: "size limit",
			query: ListQuery{
				Phone: "",
				Page: 1,
				Limit: 150,
			},
			setupRepo: func(m *mockUserRepository) {
				m.users = map[string]*userdomain.User{
					"user-1": {
						ID:        "user-1",
						Phone:     "+15551234567",
						CreatedAt: time.Now(),
					},
				}
			},
			wantPage: Page[userdomain.User]{
				Items: []userdomain.User{
					{ID: "user-1", Phone: "+15551234567"},
				},
				Total: 1,
				Page:  1,
				Limit: 20,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockUserRepository()
			tt.setupRepo(repo)

			uc := &UserUsecase{users: repo}
			page, err := uc.List(ctx, tt.query)

			if tt.wantErr {
				if err == nil {
					t.Errorf("List() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("List() unexpected error: %v", err)
				return
			}

			if page.Page != tt.wantPage.Page {
				t.Errorf("List() page = %d, want %d", page.Page, tt.wantPage.Page)
			}

			if page.Limit != tt.wantPage.Limit {
				t.Errorf("List() limit = %d, want %d", page.Limit, tt.wantPage.Limit)
			}

			if page.Total != tt.wantPage.Total {
				t.Errorf("List() total = %d, want %d", page.Total, tt.wantPage.Total)
			}

			if len(page.Items) != len(tt.wantPage.Items) {
				t.Errorf("List() items length = %d, want %d", len(page.Items), len(tt.wantPage.Items))
			}
		})
	}
}
