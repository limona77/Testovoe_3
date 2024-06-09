package repository

import (
	custom_errors "Testovoe_3/custom-errors"
	"Testovoe_3/database"
	"Testovoe_3/graph/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	*database.DB
}
type IRepository interface {
	CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error)
	UpdatePost(ctx context.Context, input model.UpdatePostInput) (*model.Post, error)
	DeletePost(ctx context.Context, id string) (bool, error)
	CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error)
	DeleteComment(ctx context.Context, id string) (bool, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetPost(ctx context.Context, id string) (*model.Post, error)
	GetComments(ctx context.Context, postID string, cursor *int, limit *int) ([]*model.Comment, error)
}

func (db *Repository) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
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

func (db *Repository) UpdatePost(ctx context.Context, input model.UpdatePostInput) (*model.Post, error) {
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

func (db *Repository) DeletePost(ctx context.Context, id string) (bool, error) {
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

func (db *Repository) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {
	path := "database.database.CreateComment"
	var allow bool
	sql := `SELECT allow_comments 
					FROM posts 
					WHERE id = $1`

	err := db.Pool.QueryRow(ctx, sql, input.PostID).
		Scan(&allow)
	if err != nil {
		return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
	}
	if !allow {
		return nil, custom_errors.ErrNotAllowed
	}
	sql2 := `INSERT INTO comments (post_id, parent_id, author, content)
					VALUES ($1, $2, $3, $4)
					RETURNING id, post_id, parent_id, author, content`
	var comment model.Comment
	err = db.Pool.QueryRow(
		ctx,
		sql2,
		input.PostID,
		input.ParentID,
		input.Author,
		input.Content).
		Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Author,
			&comment.Content)
	if err != nil {
		return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
	}
	return &comment, nil
}

func (db *Repository) DeleteComment(ctx context.Context, id string) (bool, error) {
	path := "database.database.DeleteComment"
	sql := `DELETE FROM comments
					WHERE id = $1
				  RETURNING comments.id`
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

func (db *Repository) GetPosts(ctx context.Context) ([]*model.Post, error) {
	path := "database.database.GetPosts"
	sql := `SELECT id, title, content, allow_comments FROM posts`

	var posts []*model.Post
	rows, err := db.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf(path+".Query, error: {%w}", err)
	}
	sql2 := `SELECT id, post_id, parent_id, author, content FROM comments`
	var comments []*model.Comment
	rows2, err := db.Pool.Query(ctx, sql2)
	if err != nil {
		return nil, fmt.Errorf(path+".Query, error: {%w}", err)
	}
	for rows2.Next() {
		var comment model.Comment
		err = rows2.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Author,
			&comment.Content)
		if err != nil {
			return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
		}
		comments = append(comments, &comment)
	}
	for rows.Next() {
		var post model.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AllowComments)
		if err != nil {
			return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
		}
		posts = append(posts, &post)
	}
	for _, post := range posts {
		for _, comment := range comments {
			if comment.PostID == post.ID {
				post.Comments = append(post.Comments, comment)
			}
		}
	}
	return posts, nil
}

func (db *Repository) GetPost(ctx context.Context, id string) (*model.Post, error) {
	path := "database.database.GetPost"
	sql := `SELECT id, title, content, allow_comments FROM posts
					WHERE id = $1`

	sql2 := `SELECT id, post_id, parent_id, author, content FROM comments
					WHERE post_id = $1`
	var comments []*model.Comment
	rows2, err := db.Pool.Query(ctx, sql2, id)
	if err != nil {
		return nil, fmt.Errorf(path+".Query, error: {%w}", err)
	}
	for rows2.Next() {
		var comment model.Comment
		err = rows2.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Author,
			&comment.Content)
		if err != nil {
			return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
		}
		comments = append(comments, &comment)
	}
	var post model.Post
	err = db.Pool.QueryRow(ctx, sql, id).
		Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			return nil, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, custom_errors.ErrNotFound
		}
		return nil, fmt.Errorf(path+"Scan, error: {%w}", err)
	}
	for _, v := range comments {
		if v.PostID == post.ID {
			post.Comments = append(post.Comments, v)
		}
	}
	return &post, nil
}

func (db *Repository) GetComments(
	ctx context.Context,
	postID string, cursor *int,
	limit *int,
) ([]*model.Comment, error) {
	path := "database.database.GetComments"

	sql := `SELECT id, post_id, parent_id, author, content FROM comments
					WHERE post_id = $1`
	var comments []*model.Comment
	rows, err := db.Pool.Query(ctx, sql, postID)
	if err != nil {
		return nil, fmt.Errorf(path+".Query, error: {%w}", err)
	}
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Author,
			&comment.Content)
		if err != nil {
			return nil, fmt.Errorf(path+".Scan, error: {%w}", err)
		}
		comments = append(comments, &comment)
	}
	result := comments
	if limit != nil && cursor != nil {
		start := *cursor
		end := *limit + *cursor

		if end > len(comments) {
			end = len(comments)
		}

		return result[start:end], nil
	}
	return result, nil
}
