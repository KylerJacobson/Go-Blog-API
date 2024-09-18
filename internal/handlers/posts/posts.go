package posts

import (
	"encoding/json"
	"net/http"

	"github.com/KylerJacobson/Go-Blog-API/internal/db/posts"
)

type PostsApi interface {
	GetRecentPosts(w http.ResponseWriter, r *http.Request)
}

type postsApi struct {
	postsRepository posts.PostsRepository
}

func New(postsRepo posts.PostsRepository) *postsApi {
	return &postsApi{
		postsRepository: postsRepo,
	}
}

func (postsApi *postsApi) GetRecentPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postsApi.postsRepository.GetRecentPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func (postsApi *postsApi) GetRecentPublicPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := postsApi.postsRepository.GetRecentPublicPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.Marshal(err)
		w.Write(b)
		return
	}
	b, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
