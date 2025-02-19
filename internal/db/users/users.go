package users

import (
	"context"
	"errors"

	user_models "github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	"github.com/jackc/pgx/v5"
	pgxv5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	CreateUser(user user_models.UserCreate) (string, error)
	UpdateUser(user user_models.UserUpdate) error
	GetUserById(id int) (*user_models.User, error)
	GetUserByEmail(email string) (*user_models.User, error)
	GetAllUsers() (*[]user_models.FrontendUser, error)
	DeleteUserById(id int) error
	LoginUser(user user_models.UserLogin) (*user_models.User, error)
}

type usersRepository struct {
	conn   *pgxpool.Pool
	logger logger.Logger
}

func New(conn *pgxpool.Pool, logger logger.Logger) *usersRepository {
	return &usersRepository{
		conn:   conn,
		logger: logger,
	}
}

func (repository *usersRepository) GetUserById(id int) (*user_models.User, error) {
	repository.logger.Sugar().Infof("getting user from the database")

	rows, err := repository.conn.Query(
		context.TODO(), `SELECT id, first_name, last_name, email, password, created_at, updated_at, role, email_notification FROM users WHERE id = $1;`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_models.User])
	if err != nil {
		repository.logger.Sugar().Errorf("Error getting user: %v", err)
		return nil, err
	}
	if len(users) < 1 {
		return nil, nil
	}

	return &users[0], nil
}

func (repository *usersRepository) DeleteUserById(id int) error {
	repository.logger.Sugar().Infof("deleting user from the database")

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
		repository.logger.Sugar().Errorf("Error creating user %s %s : %v", user.FirstName, user.FirstName, err)
		return "", err
	}
	defer rows.Close()
	createdUser, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_models.User])
	if err != nil {
		repository.logger.Sugar().Errorf("Error returning user %s %s from database: %v", user.FirstName, user.FirstName, err)
		return "", err
	}

	if createdUser[0].Id == "" {
		repository.logger.Sugar().Errorf("Error returning user %s %s from database: %v", user.FirstName, user.FirstName, err)
		return "", err
	}

	return createdUser[0].Id, nil
}

func (repository *usersRepository) UpdateUser(user user_models.UserUpdate) error {
	rows, err := repository.conn.Query(context.TODO(), `UPDATE users SET first_name = $1, last_name = $2, email = $3, role = $4, email_notification = $5 WHERE id = $6 `, user.FirstName, user.LastName, user.Email, user.Role, user.EmailNotification, user.Id)
	if err != nil {
		repository.logger.Sugar().Errorf("Error updating user %s %s : %v", user.FirstName, user.FirstName, err)
		return err
	}
	defer rows.Close()
	return nil
}

func (repository *usersRepository) GetUserByEmail(email string) (*user_models.User, error) {
	rows, err := repository.conn.Query(context.TODO(), `SELECT id, first_name, last_name, email, password, created_at, updated_at, role, email_notification FROM users WHERE email = $1`, email)
	if err != nil {
		repository.logger.Sugar().Errorf("Error retrieving user (%s) from the database: %v", email, err)
		return nil, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_models.User])
	if err != nil {
		repository.logger.Sugar().Errorf("Error getting user: %v", err)
		return nil, err
	}
	if len(users) < 1 {
		repository.logger.Sugar().Errorf("User %s not found: %v", email, err)
		return nil, errors.New("User not found")
	}
	return &users[0], nil
}

func (repository *usersRepository) LoginUser(user user_models.UserLogin) (*user_models.User, error) {
	var match bool
	err := repository.conn.QueryRow(
		context.TODO(), `SELECT (password = crypt($1, password)) AS isMatch FROM users WHERE email = $2`, user.Password, user.Email,
	).Scan(&match)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			repository.logger.Sugar().Infof("User with id: %s does not exist in the database", user.Email)
			return nil, nil
		}
		repository.logger.Sugar().Errorf("Error retrieving user (%s) from the database: %v", user.Email, err)
		return nil, err
	}
	if match {
		user, err := repository.GetUserByEmail(user.Email)
		if err != nil {
			repository.logger.Sugar().Errorf("Error getting user: %v", err)
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

func (repository *usersRepository) GetAllUsers() (*[]user_models.FrontendUser, error) {
	rows, err := repository.conn.Query(context.TODO(), `SELECT id, first_name, last_name, email, role, email_notification, created_at FROM users ORDER BY created_at ASC`)
	if err != nil {
		repository.logger.Sugar().Errorf("Error retrieving users from the database: %v", err)
		return nil, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_models.FrontendUser])
	if err != nil {
		repository.logger.Sugar().Errorf("Error getting user: %v", err)
		return nil, err
	}
	return &users, nil
}
