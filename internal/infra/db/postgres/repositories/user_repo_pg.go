package postgresrepositories

import (
	"context"

	userdomain "dekamond/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) GetByPhone(ctx context.Context, phone string) (*userdomain.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, phone, created_at FROM users WHERE phone=$1`, phone)
	var u userdomain.User
	if err := row.Scan(&u.ID, &u.Phone, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, phone string) (*userdomain.User, error) {
	id := uuid.New().String()
	row := r.pool.QueryRow(ctx, `INSERT INTO users(id, phone) VALUES($1,$2) RETURNING id, phone, created_at`, id, phone)
	var u userdomain.User
	if err := row.Scan(&u.ID, &u.Phone, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, phone, created_at FROM users WHERE id=$1`, id)
	var u userdomain.User
	if err := row.Scan(&u.ID, &u.Phone, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, phone string, limit, offset int) ([]userdomain.User, int, error) {
	q := `SELECT id, phone, created_at FROM users WHERE ($1 = '' OR phone = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, phone, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	users := make([]userdomain.User, 0)
	for rows.Next() {
		var u userdomain.User
		if err := rows.Scan(&u.ID, &u.Phone, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE ($1 = '' OR phone = $1)`, phone).Scan(&total); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}