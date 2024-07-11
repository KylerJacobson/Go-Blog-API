package posts

import (
	"context"
	"fmt"

	post_models "github.com/KylerJacobson/Go-Blog-API/internal/api/types/posts"
	"github.com/jackc/pgx/v5"
)

type PostsRepoistory interface {
	Get() (*post_models.Post, error)
}

type postsRepository struct {
	conn pgx.Conn
}

func New(conn pgx.Conn) *postsRepository {
	return &postsRepository{
		conn: conn,
	}
}

func (repository *postsRepository) Get() (*post_models.Post, error) {
	fmt.Println("getting posts from the database")

	rows, err := repository.conn.Query(
		context.TODO(), `SELECT * FROM posts ORDER BY created_at DESC LIMIT 10;`,
	)
	if err != nil {
		return nil,nil
	}
	defer rows.Close()

	posts, _ := pgx.CollectRows(rows, pgx.RowToStructByName[post_models.Post])
	fmt.Printf("postId: %d postBlob %s ", posts[0].PostId, posts[0].BlobName)
	return nil,nil
}