package posts

import (
	"context"

	post_models "github.com/KylerJacobson/Go-Blog-API/internal/api/types/posts"
	"github.com/KylerJacobson/Go-Blog-API/logger"
	"github.com/jackc/pgx/v5"
)

type PostsRepository interface {
	GetRecentPosts() ([]post_models.Post, error)
	GetRecentPublicPosts() ([]post_models.Post, error)
	GetPostById(postId int) (*post_models.Post, error)
	// DeletePostById(postId int) (post_models.Post, error)
	// CreatePost(title, content string, restricted bool, userId int) error
}

type postsRepository struct {
	conn pgx.Conn
}

func New(conn pgx.Conn) *postsRepository {
	return &postsRepository{
		conn: conn,
	}
}

func (repository *postsRepository) GetRecentPosts() ([]post_models.Post, error) {
	logger.Sugar.Infof("getting posts from the database")

	rows, err := repository.conn.Query(
		context.TODO(), `SELECT * FROM posts ORDER BY created_at DESC LIMIT 10;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[post_models.Post])
	if err != nil {
		logger.Sugar.Errorf("Error getting recent posts from the database: %v", err)
		return nil, err
	}
	if len(posts) > 0 {
		logger.Sugar.Infof("postId: %d postBlob %s ", posts[0].PostId, posts[0].Title)
	}
	return posts, nil
}

func (repository *postsRepository) GetRecentPublicPosts() ([]post_models.Post, error) {
	logger.Sugar.Info("getting public posts from the database")

	rows, err := repository.conn.Query(
		context.TODO(), `SELECT * FROM posts WHERE restricted = false ORDER BY created_at DESC LIMIT 10`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[post_models.Post])
	if err != nil {
		logger.Sugar.Errorf("Error getting recent public posts from the database: %v ", err)
		return nil, err
	}
	if len(posts) > 0 {
		logger.Sugar.Infof("postId: %d postBlob %s ", posts[0].PostId, posts[0].Title)
	}
	return posts, nil
}

func (repository *postsRepository) GetPostById(postId int) (*post_models.Post, error) {

	rows, err := repository.conn.Query(
		context.TODO(), `SELECT * FROM posts WHERE post_id = $1`, postId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	post, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[post_models.Post])
	//post, err := pgx.CollectRows(rows, pgx.RowToStructByName[post_models.Post])
	if err != nil {
		logger.Sugar.Errorf("Error getting post %v: %v ", err)
		return nil, err
	}
	return &post, nil

}
