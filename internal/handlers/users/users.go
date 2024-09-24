package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	users_repo "github.com/KylerJacobson/Go-Blog-API/internal/db/users"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	pgxv5 "github.com/jackc/pgx/v5"
)

type UsersApi interface {
	GetUserById(w http.ResponseWriter, r *http.Request)
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
			http.Error(w, "Post not found", http.StatusNotFound)
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
