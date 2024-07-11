package posts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
)

type PostsApi interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type postsApi struct {
	postsRepository posts.PostsRepoistory
}

func New(postsRepo posts.PostsRepoistory) *postsApi {
	return &postsApi{
		postsRepository: postsRepo,
	}
}

func (postsApi *postsApi) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting posts")
	post, err := postsApi.postsRepository.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	
}
