package database

import (
	custom_errors "Testovoe_3/custom-errors"
	"Testovoe_3/graph/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(
		ctx context.Context,
		tableName pgx.Identifier,
		columnNames []string,
		rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}
type DB struct {
	Pool PgxPool
}

var url = "postgres://postgres:postgres@localhost:5432/testovoe_3?sslmode=disable"

func Connect() *DB {
	db := &DB{}
	var err error
	db.Pool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		panic("can't connect to Postgres")
	}
	_, err = pgx.Connect(context.Background(), url)
	fmt.Println(url)
	if err != nil {
		panic("can't connect to Postgres")
	}
	return db
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *DB) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	path := "database.database.CreatePost"
	var post model.Post
	sql := `INSERT INTO posts (title, content, allow_comments) 
					VALUES ($1, $2, $3) 
					RETURNING id, title, content, allow_comments`

	err := db.Pool.QueryRow(
		ctx,
		sql,
		input.Title, input.Content, input.AllowComments).
		Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments)
	if err != nil {
		return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
	}
	return &post, nil
}

func (db *DB) UpdatePost(ctx context.Context, input model.UpdatePostInput) (*model.Post, error) {
	path := "database.database.UpdatePost"
	sql := `UPDATE posts 
					SET title = $1, content = $2, allow_comments = $3 
					WHERE id = $4
					RETURNING id, title, content, allow_comments`
	var post model.Post
	err := db.Pool.QueryRow(
		ctx,
		sql,
		input.Title,
		input.Content,
		input.AllowComments,
		input.ID).
		Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments)
	if err != nil {
		return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
	}
	return &post, nil
}

func (db *DB) DeletePost(ctx context.Context, id string) (bool, error) {
	path := "database.database.DeletePost"
	sql := `DELETE FROM posts
					WHERE id = $1
				  RETURNING posts.id`
	var deletedId int
	err := db.Pool.QueryRow(ctx, sql, id).
		Scan(&deletedId)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			return false, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return false, custom_errors.ErrNotFound
		}
		return false, fmt.Errorf(path+".QueryRow, error: {%w}", err)
	}
	if deletedId == 0 {
		return false, fmt.Errorf(path + custom_errors.ErrNotFound.Error())
	}
	return true, nil
}

func (db *DB) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	path := "database.database.CreateComment"
	var comment model.Comment
	sql := `INSERT INTO comments (post_id, parent_id, author, content)
					VALUES ($1, $2, $3, $4)
					RETURNING id, post_id, parent_id, author, content`
	err := db.Pool.QueryRow(ctx, sql, input.PostID, input.ParentID, input.Author, input.Content).
		Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Author, &comment.Content)
	if err != nil {
		return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
	}
	return &comment, nil
}
