package repository

import (
	"context"

	"github.com/arsyad-id-99/task-manager-go/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User, passwordHash string) error {
	query := `
        INSERT INTO users (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, user.Name, user.Email, passwordHash).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, string, error) {
	user := &model.User{}
	var passwordHash string
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt)
	return user, passwordHash, err
}
