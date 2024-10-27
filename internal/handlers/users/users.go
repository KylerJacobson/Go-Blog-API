package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	"github.com/KylerJacobson/Go-Blog-API/internal/authorization"
	users_repo "github.com/KylerJacobson/Go-Blog-API/internal/db/users"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/session"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	pgxv5 "github.com/jackc/pgx/v5"
)

type UsersApi interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	GetUserFromSession(w http.ResponseWriter, r *http.Request)
}

type usersApi struct {
	usersRepository users_repo.UsersRepository
}

func New(usersRepo users_repo.UsersRepository) *usersApi {
	return &usersApi{
		usersRepository: usersRepo,
	}
}

func (usersApi *usersApi) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	val, err := strconv.Atoi(id)
	if err != nil {
		logger.Sugar.Errorf("GetPostId parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	user, err := usersApi.usersRepository.GetUserById(val)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			logger.Sugar.Infof("User with id: %d does not exist in the database", val)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	if user == nil {
		logger.Sugar.Infof("user with id %d not found", val)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		logger.Sugar.Errorf("error marshalling user : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	val, err := strconv.Atoi(id)
	if err != nil {
		logger.Sugar.Errorf("Delete user parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	err = usersApi.usersRepository.DeleteUserById(val)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			logger.Sugar.Infof("User with id: %d does not exist in the database", val)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (usersApi *usersApi) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest users.UserCreate
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		logger.Sugar.Errorf("Error decoding the user request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	userId, err := usersApi.usersRepository.CreateUser(userRequest)
	if err != nil {
		logger.Sugar.Errorf("error creating user for %s %s : %v", userRequest.FirstName, userRequest.LastName, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, err := json.Marshal(userId)
	if err != nil {
		logger.Sugar.Errorf("error marshalling the create user response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest users.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
	if err != nil {
		logger.Sugar.Errorf("Error decoding the user request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	user, err := usersApi.usersRepository.LoginUser(userLoginRequest)
	if err != nil {
		logger.Sugar.Errorf("error logging in user for %s : %v", userLoginRequest.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		logger.Sugar.Errorf("error marshalling the login user response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := usersApi.usersRepository.GetAllUsers()
	if err != nil {
		logger.Sugar.Errorf("error listing users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, err := json.Marshal(users)
	if err != nil {
		logger.Sugar.Errorf("error marshalling the users list: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) GetUserFromSession(w http.ResponseWriter, r *http.Request) {

	token := session.Manager.GetString(r.Context(), "session_token")

	fmt.Println(token)
	claims := authorization.DecodeToken(token)
	fmt.Println(claims)
	user, err := usersApi.usersRepository.GetUserById(claims.Sub)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			logger.Sugar.Infof("User with id: %d does not exist in the database", claims.Sub)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	if user == nil {
		logger.Sugar.Infof("user with id %d not found", claims.Sub)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		logger.Sugar.Errorf("error marshalling user : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
