package domain

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sivaprasadreddy/bookmarks/internal/config"
)

type BookmarkRepository interface {
	GetAll(ctx context.Context) ([]Bookmark, error)
	GetByID(ctx context.Context, id int) (*Bookmark, error)
	Create(ctx context.Context, b Bookmark) (*Bookmark, error)
	Update(ctx context.Context, b Bookmark) error
	Delete(ctx context.Context, id int) error
}

type bookmarkRepo struct {
	db     *pgx.Conn
	logger *config.Logger
}

func NewBookmarkRepository(db *pgx.Conn, logger *config.Logger) BookmarkRepository {
	return bookmarkRepo{
		db:     db,
		logger: logger,
	}
}

func (r bookmarkRepo) GetAll(ctx context.Context) ([]Bookmark, error) {
	query := `select id, title, url, created_at FROM bookmarks`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []Bookmark
	for rows.Next() {
		var bookmark = Bookmark{}
		err = rows.Scan(&bookmark.ID, &bookmark.Title, &bookmark.Url, &bookmark.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (r bookmarkRepo) GetByID(ctx context.Context, id int) (*Bookmark, error) {
	r.logger.Infof("Fetching bookmark with id=%d", id)
	var b = Bookmark{}
	query := "select id, title, url, created_at FROM bookmarks where id=$1"
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Url, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r bookmarkRepo) Create(ctx context.Context, b Bookmark) (*Bookmark, error) {
	query := "insert into bookmarks(title, url, created_at) values($1, $2, $3) RETURNING id"
	var lastInsertID int
	err := r.db.QueryRow(ctx, query, b.Title, b.Url, b.CreatedAt).Scan(&lastInsertID)
	if err != nil {
		r.logger.Errorf("Error while inserting url row: %v", err)
		return nil, err
	}
	b.ID = lastInsertID
	return &b, nil
}

func (r bookmarkRepo) Update(ctx context.Context, b Bookmark) error {
	query := "update bookmarks set title = $1, url=$2 where id=$3"
	_, err := r.db.Exec(ctx, query, b.Title, b.Url, b.ID)
	return err
}

func (r bookmarkRepo) Delete(ctx context.Context, id int) error {
	query := "delete from bookmarks where id=$1"
	_, err := r.db.Exec(ctx, query, id)
	return err
}
