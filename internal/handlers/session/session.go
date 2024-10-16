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
	"github.com/golang-jwt/jwt/v5"
)
type UserClaim struct {
	Sub int `json:"sub"`
	Role int `json:"role"`
	jwt.RegisteredClaims
}

type sessionApi struct {
	usersRepository users_repo.UsersRepository
}
func New(usersRepo users_repo.UsersRepository) *sessionApi {
	return &sessionApi{
		usersRepository: usersRepo,
	}
}
func (sessionApi *sessionApi) CreateSession(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest users.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
	if err != nil {
		logger.Sugar.Errorf("Error decoding the user request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	user, err := sessionApi.usersRepository.LoginUser(userLoginRequest)
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

}