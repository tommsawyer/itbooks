package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Book represents book table in postgres.
type Book struct {
	ID          int64                     `db:"id"`
	ISBN        pgtype.Text               `db:"isbn"`
	URL         pgtype.Text               `db:"url"`
	Title       pgtype.Text               `db:"title"`
	Image       pgtype.Text               `db:"image"`
	Description pgtype.Text               `db:"description"`
	Authors     pgtype.Array[pgtype.Text] `db:"authors"`
	Publisher   pgtype.Text               `db:"publisher"`
	Properties  map[string]string         `db:"properties"`
	CreatedAt   pgtype.Timestamp          `db:"created_at"`
	UpdatedAt   pgtype.Timestamp          `db:"updated_at"`
}

func (b *Book) scan(row pgx.Row) error {
	return row.Scan(
		&b.ID,
		&b.ISBN,
		&b.URL,
		&b.Title,
		&b.Image,
		&b.Description,
		&b.Authors,
		&b.Properties,
		&b.Publisher,
	)
}

// UpsertBookParams is parameters required for inserting book.
type UpsertBookParams struct {
	ISBN        string
	URL         string
	Title       string
	Image       string
	Description string
	Authors     []string
	Publisher   string
	Properties  map[string]string
}

// UpsertBook creates book in postgres and returns ID.
func UpsertBook(ctx context.Context, params UpsertBookParams) (int64, error) {
	query, args, err := psql.Insert("books").Columns(
		"isbn", "url", "title", "image",
		"description", "authors", "properties", "publisher",
	).Values(
		params.ISBN, params.URL, params.Title, params.Image,
		params.Description, params.Authors, params.Properties, params.Publisher,
	).Suffix(`
      ON CONFLICT(isbn) DO UPDATE 
      SET 
        title=EXCLUDED.title, 
        url=EXCLUDED.url, 
        image=EXCLUDED.image,
        authors=EXCLUDED.authors,
        properties=EXCLUDED.properties,
        publisher=EXCLUDED.publisher,
        description=EXCLUDED.description
    `).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return 0, err
	}

	var id int64

	err = getDB(ctx).QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("cannot create book: %w", err)
	}

	return id, nil
}

// GetBook returns book by id.
func GetBook(ctx context.Context, id int64) (*Book, error) {
	query, params, err := psql.Select(
		"id",
		"isbn",
		"url",
		"title",
		"image",
		"description",
		"authors",
		"properties",
		"publisher",
	).From("books").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var book Book
	row := getDB(ctx).QueryRow(ctx, query, params...)
	err = book.scan(row)
	if err != nil {
		return nil, fmt.Errorf("cannot get book: %w", err)
	}

	return &book, nil
}

// GetBookByISBN returns book by isbn.
func GetBookByISBN(ctx context.Context, isbn string) (*Book, error) {
	query, params, err := psql.Select(
		"id",
		"isbn",
		"url",
		"title",
		"image",
		"description",
		"authors",
		"properties",
		"publisher",
	).From("books").
		Where(sq.Eq{"isbn": isbn}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var book Book
	row := getDB(ctx).QueryRow(ctx, query, params...)
	err = book.scan(row)
	if err != nil {
		return nil, fmt.Errorf("cannot get book: %w", err)
	}

	return &book, nil
}

// SetBookPublished sets published flag on book.
func SetBookPublished(ctx context.Context, id int64, published bool) error {
	query, params, err := psql.Update("books").
		Set("published", published).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = getDB(ctx).Exec(ctx, query, params...)
	return err
}

// FindUnpublishedBooks returns unpublished books.
func FindUnpublishedBooks(ctx context.Context) ([]*Book, error) {
	return findBooks(ctx, sq.Eq{"published": false})
}

// FindBooks returns all books.
func FindBooks(ctx context.Context) ([]*Book, error) {
	return findBooks(ctx, nil)
}

func findBooks(ctx context.Context, filter any) ([]*Book, error) {
	q := psql.Select(
		"id",
		"isbn",
		"url",
		"title",
		"image",
		"description",
		"authors",
		"properties",
		"publisher",
	).From("books")

	if filter != nil {
		q = q.Where(filter)
	}

	q = q.OrderBy("created_at")

	query, params, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var books []*Book
	rows, err := getDB(ctx).Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("cannot find books: %w", err)
	}

	for rows.Next() {
		var book Book
		if err := book.scan(rows); err != nil {
			return nil, fmt.Errorf("cannot scan book: %w", err)
		}

		books = append(books, &book)
	}

	return books, nil
}
