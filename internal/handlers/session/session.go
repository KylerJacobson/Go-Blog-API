package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/KylerJacobson/Go-Blog-API/internal/api/types/users"
	users_repo "github.com/KylerJacobson/Go-Blog-API/internal/db/users"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	"github.com/alexedwards/scs/v2"
	"github.com/golang-jwt/jwt/v5"
)

var Manager *scs.SessionManager

type UserClaim struct {
	Sub  int `json:"sub"`
	Role int `json:"role"`
	jwt.RegisteredClaims
}

type SessionApi interface {
	CreateSession(w http.ResponseWriter, r *http.Request)
	DeleteSession(w http.ResponseWriter, r *http.Request)
}
type sessionApi struct {
	usersRepository users_repo.UsersRepository
}

func New(usersRepo users_repo.UsersRepository) *sessionApi {
	return &sessionApi{
		usersRepository: usersRepo,
	}
}

func Init() {
	Manager = scs.New()
	Manager.Lifetime = 3 * time.Hour
	// Manager.IdleTimeout = 20 * time.Minute
	// Manager.Cookie.Domain = "kylerjacobson.dev"
	// Manager.Cookie.HttpOnly = true
	// Manager.Cookie.Persist = true
	// Manager.Cookie.SameSite = http.SameSiteStrictMode
	// Manager.Cookie.Secure = true
}

func (sessionApi *sessionApi) CreateSession(w http.ResponseWriter, r *http.Request) {
	var userLoginFormRequest users.UserLoginForm
	err := json.NewDecoder(r.Body).Decode(&userLoginFormRequest)
	if err != nil {
		logger.Sugar.Errorf("Error decoding the user request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	user, err := sessionApi.usersRepository.LoginUser(userLoginFormRequest.FormData)
	if err != nil {
		logger.Sugar.Errorf("error logging in user for %s : %v", userLoginFormRequest.FormData.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	iId, _ := strconv.Atoi(user.Id)

	claims := UserClaim{
		iId,
		user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "kylerjacobson.dev",
		},
	}

	// Sign and get the complete encoded token as a string using the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	fmt.Println(ss, err)
	Manager.Put(r.Context(), "session_token", ss)
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(ss)
	w.Write(b)
	return
}

func (sessionApi *sessionApi) DeleteSession(w http.ResponseWriter, r *http.Request) {
	Manager.Put(r.Context(), "session_token", "")
	w.WriteHeader(http.StatusOK)
	return
}
