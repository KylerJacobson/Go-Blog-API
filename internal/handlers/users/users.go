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
	var userRequest users.UserCreate
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("Error decoding the user request body: %v", err)
		httperr.Write(w, httperr.BadRequest("Invalid request body", err.Error()))
		return
	}

	err = validateCreateUserRequest(userRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error validating user create request", err)
		httperr.Write(w, httperr.BadRequest("Invalid request body", err.Error()))
		return
	}

	userId, err := usersApi.usersRepository.CreateUser(userRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error creating user for %s %s : %v", userRequest.FirstName, userRequest.LastName, err)
		httperr.Write(w, httperr.Internal("failed to create user", err.Error()))
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
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
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
	users, err := usersApi.usersRepository.GetAllUsers()
	if err != nil {
		usersApi.logger.Sugar().Errorf("error listing users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, err := json.Marshal(users)
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
	var userUpdateRequest users.UserUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&userUpdateRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("Error decoding the user request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	// validate user update request
	err = usersApi.validateUpdateUserRequest(r, userUpdateRequest)
	if err != nil {
		usersApi.logger.Sugar().Errorf("error validating user update request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	err = usersApi.usersRepository.UpdateUser(userUpdateRequest.User)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			usersApi.logger.Sugar().Infof("User with id: %s does not exist in the database", userUpdateRequest.User.Id)
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

func (usersApi *usersApi) validateUpdateUserRequest(r *http.Request, userUpdateRequest users.UserUpdateRequest) error {
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
	if claims.Sub != userID {
		usersApi.logger.Sugar().Errorf("user %d attempted to update user %d", claims.Sub, userID)
		errors = append(errors, err)
	}

	// Validate fields
	if userUpdateRequest.User.FirstName == "" {
		errors = append(errors, fmt.Errorf("first name is required"))
	}
	if userUpdateRequest.User.LastName == "" {
		errors = append(errors, fmt.Errorf("last name is required"))
	}
	if userUpdateRequest.User.Email == "" {
		errors = append(errors, fmt.Errorf("email is required"))
	}
	if userUpdateRequest.Role == "" {
		errors = append(errors, fmt.Errorf("logged in user role is required"))
	}
	if userUpdateRequest.User.Role == "" {
		errors = append(errors, fmt.Errorf("new role is required"))
	}
	// Using strings here because the nil value of an int is a valid role
	newRoleRequest, err := strconv.Atoi(userUpdateRequest.User.Role)
	if err != nil {
		usersApi.logger.Sugar().Errorf("new user role is not an integer: %v", err)
		errors = append(errors, err)
	}
	currentUserRole, err := strconv.Atoi(userUpdateRequest.Role)
	if err != nil {
		usersApi.logger.Sugar().Errorf("current user role is not an integer: %v", err)
		errors = append(errors, err)
	}
	if currentUserRole == 1 && newRoleRequest != 1 {
		usersApi.logger.Sugar().Errorf("user is already an admin, cannot decrease permission: %v", err)
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return fmt.Errorf("validation failed")
	}
	return nil
}
