package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	"net/http"
	"strconv"
	"strings"

	"github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	"github.com/KylerJacobson/Go-Blog-API/internal/authorization"
	users_repo "github.com/KylerJacobson/Go-Blog-API/internal/db/users"
	"github.com/KylerJacobson/Go-Blog-API/internal/handlers/session"
	"github.com/KylerJacobson/Go-Blog-API/internal/httperr"
	pgxv5 "github.com/jackc/pgx/v5"
)

type UsersApi interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	GetUserFromSession(w http.ResponseWriter, r *http.Request)
}

type usersApi struct {
	usersRepository users_repo.UsersRepository
	logger          logger.Logger
}

func New(usersRepo users_repo.UsersRepository, logger logger.Logger) *usersApi {
	return &usersApi{
		usersRepository: usersRepo,
		logger:          logger,
	}
}

func (usersApi *usersApi) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	val, err := strconv.Atoi(id)
	if err != nil {
		usersApi.logger.Sugar().Errorf("GetPostId parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	user, err := usersApi.usersRepository.GetUserById(val)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			usersApi.logger.Sugar().Infof("User with id: %d does not exist in the database", val)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	if user == nil {
		usersApi.logger.Sugar().Infof("user with id %d not found", val)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error marshalling user : %v", err)
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
		usersApi.logger.Sugar().Errorf("Delete user parameter was not an integer: %v", err)
		http.Error(w, "postId must be an integer", http.StatusBadRequest)
		return
	}
	err = usersApi.usersRepository.DeleteUserById(val)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			usersApi.logger.Sugar().Infof("User with id: %d does not exist in the database", val)
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
	var accountCreationRequest users.AccountCreationRequest
	err := json.NewDecoder(r.Body).Decode(&accountCreationRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("Error decoding the user request body: %v", err)
		httperr.Write(w, httperr.BadRequest("Invalid request body", err.Error()))
		return
	}

	err = validateCreateUserRequest(accountCreationRequest.User)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error validating user create request", err)
		httperr.Write(w, httperr.BadRequest("Invalid request body", err.Error()))
		return
	}

	userId, err := usersApi.usersRepository.CreateUser(accountCreationRequest.User)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error creating user for %s %s : %v", accountCreationRequest.User.FirstName, accountCreationRequest.User.LastName, err)
		httperr.Write(w, httperr.New(http.StatusInternalServerError, "failed to create user", ""))
		return
	}

	w.WriteHeader(http.StatusOK)
	// Set content header to application/json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": userId,
	})
}

// create a function to validate the request body on create user

func validateCreateUserRequest(userRequest users.UserCreate) error {
	var errors []string
	if userRequest.FirstName == "" {
		errors = append(errors, "first name is required")
	}
	if userRequest.LastName == "" {
		errors = append(errors, "last name is required")
	}
	if userRequest.Email == "" {
		errors = append(errors, "email is required")
	}
	if userRequest.Password == "" {
		errors = append(errors, "password is required")
	}
	if len(userRequest.Password) < 8 {
		errors = append(errors, "password must be at least 8 characters long")
	}
	if strings.Contains(userRequest.Email, "@") == false {
		errors = append(errors, "invalid email format")
	}
	if userRequest.AccessRequest != -1 && userRequest.AccessRequest != 0 && userRequest.AccessRequest != 2 {
		errors = append(errors, "access request must be -1, 0 or 2")
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}

func (usersApi *usersApi) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest users.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("Error decoding the user request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if userLoginRequest.Email == "" || userLoginRequest.Password == "" {
		usersApi.logger.Sugar().Errorf("error validating user login request: %v", err)
		http.Error(w, "Invalid request body - email or password is missing", http.StatusBadRequest)
		return
	}
	user, err := usersApi.usersRepository.LoginUser(userLoginRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error logging in user for %s : %v", userLoginRequest.Email, err)
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
		usersApi.logger.Sugar().Errorf("error marshalling the login user response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) ListUsers(w http.ResponseWriter, r *http.Request) {

	// Check if user is logged in and has admin role
	allUsers, err := usersApi.usersRepository.GetAllUsers()
	if err != nil {
		usersApi.logger.Sugar().Errorf("error listing users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, err := json.Marshal(allUsers)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error marshalling the users list: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userUpdate users.UserUpdate
	err := json.NewDecoder(r.Body).Decode(&userUpdate)
	if err != nil {
		usersApi.logger.Sugar().Errorf("Error decoding the user request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	// validate user update request
	err = usersApi.validateUpdateUserRequest(r, userUpdate)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error validating user update request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}

	err = usersApi.usersRepository.UpdateUser(userUpdate)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			usersApi.logger.Sugar().Infof("User with id: %s does not exist in the database", userUpdate.Id)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (usersApi *usersApi) GetUserFromSession(w http.ResponseWriter, r *http.Request) {

	loggedIn := session.Manager.Exists(r.Context(), "session_token")
	if !loggedIn {
		//usersApi.logger.Sugar().Errorf("user not logged in")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	token := session.Manager.GetString(r.Context(), "session_token")
	fmt.Println(token)
	claims := authorization.DecodeToken(token)
	fmt.Println(claims)
	user, err := usersApi.usersRepository.GetUserById(claims.Sub)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			usersApi.logger.Sugar().Infof("User with id: %d does not exist in the database", claims.Sub)
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	if user == nil {
		usersApi.logger.Sugar().Infof("user with id %d not found", claims.Sub)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error marshalling user : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (usersApi *usersApi) validateUpdateUserRequest(r *http.Request, userUpdate users.UserUpdate) error {
	// Check path value matches current userID
	errors := []error{}
	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		usersApi.logger.Sugar().Errorf("Update user parameter was not an integer: %v", err)
		errors = append(errors, err)
	}
	token := session.Manager.GetString(r.Context(), "session_token")
	claims := authorization.DecodeToken(token)
	if claims.Sub != userID && claims.Role != 1 {
		usersApi.logger.Sugar().Errorf("user %d attempted to update user %d", claims.Sub, userID)
		errors = append(errors, err)
	}

	// Validate permission escalation / deescalation
	if claims.Sub == 1 && userID == 1 && userUpdate.Role != 1 {
		usersApi.logger.Sugar().Errorf("user is already an admin, cannot decrease permission: %v", err)
		errors = append(errors, err)
	}
	if claims.Role != 1 && userUpdate.Role == 1 {
		usersApi.logger.Sugar().Errorf("access denied: %v", err)
		errors = append(errors, err)
	}

	// Validate fields
	if userUpdate.FirstName == "" {
		errors = append(errors, fmt.Errorf("first name is required"))
	}
	if userUpdate.LastName == "" {
		errors = append(errors, fmt.Errorf("last name is required"))
	}
	if userUpdate.Email == "" {
		errors = append(errors, fmt.Errorf("email is required"))
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed")
	}
	return nil
}
