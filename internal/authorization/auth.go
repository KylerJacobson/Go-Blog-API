package authorization

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	Sub int `json:"sub"`
	Role int `json:"role"`
	jwt.RegisteredClaims
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
        token = strings.TrimPrefix(token, "Bearer ")
    }
	key := os.Getenv("JWT_SECRET")
	parsedToken, err := jwt.ParseWithClaims(token, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		log.Fatal(err)
	} else if claims, ok := parsedToken.Claims.(*UserClaim); ok {
		fmt.Println(claims.Sub, claims.RegisteredClaims.Issuer)
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
}