package users

import (
	"context"

	user_models "github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	"github.com/jackc/pgx/v5"
)

type UsersRepository interface {
	GetUserById(id int) (*user_models.User, error)
}

type usersRepository struct {
	conn pgx.Conn
}

func New(conn pgx.Conn) *usersRepository {
	return &usersRepository{
		conn: conn,
	}
}

func (repository *usersRepository) GetUserById(id int) (*user_models.User, error) {
	logger.Sugar.Infof("getting user from the database")

	rows, err := repository.conn.Query(
		context.TODO(), `SELECT id, first_name, last_name, email, password, created_at, updated_at, role, email_notification FROM users WHERE id = $1;`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_models.User])
	if err != nil {
		logger.Sugar.Errorf("Error getting user: %v", err)
		return nil, err
	}
	if len(users) < 1 {
		return nil, nil
	}

	return &users[0], nil
}
