package users

import (
	"context"

	user_models "github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	"github.com/jackc/pgx/v5"
)

type UsersRepository interface {
	GetUserById(id int) (*user_models.User, error)
	DeleteUserById(id int) error
	CreateUser(user user_models.UserCreate) (string, error)
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

func (repository *usersRepository) DeleteUserById(id int) error {
	logger.Sugar.Infof("deleting user from the database")

	rows, err := repository.conn.Query(
		context.TODO(), `DELETE FROM users WHERE id = $1;`, id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (repository *usersRepository) CreateUser(user user_models.UserCreate) (string, error) {

	rows, err := repository.conn.Query(context.TODO(), `INSERT INTO users (first_name, last_name, email, password, role, email_notification) VALUES ($1, $2, $3, crypt($4, gen_salt('bf', 8)), $5, $6) RETURNING *`, user.FirstName, user.LastName, user.Email, user.Password, user.AccessRequest, user.EmailNotification)
	if err != nil {
		logger.Sugar.Errorf("Error creating user %s %s : %v", user.FirstName, user.FirstName, err)
		return "", err
	}
	defer rows.Close()
	createdUser, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_models.User])
	if err != nil {
		logger.Sugar.Errorf("Error returning user %s %s from database: %v", user.FirstName, user.FirstName, err)
		return "", err
	}

	if createdUser[0].Id == "" {
		logger.Sugar.Errorf("Error returning user %s %s from database: %v", user.FirstName, user.FirstName, err)
		return "", err
	}

	return createdUser[0].Id, nil
}
