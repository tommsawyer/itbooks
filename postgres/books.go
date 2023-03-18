package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var bookColumns = []string{
	"id",
	"isbn",
	"url",
	"title",
	"image",
	"description",
	"authors",
	"properties",
	"publisher",
	"created_at",
	"updated_at",
}

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
		&b.CreatedAt,
		&b.UpdatedAt,
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
//
// If row with the same ISBN already exists it will just update fields of existing row
// and returns id of old book.
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
    `,
	).Suffix("RETURNING id").ToSql()
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

// UpdateBook updates given fields on book.
func UpdateBook(ctx context.Context, id int64, fields Fields) error {
	query, params, err := psql.Update("books").
		SetMap(fields).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = getDB(ctx).Exec(ctx, query, params...)
	return err
}

// GetBook returns first found book by given filter.
//
// Use squirrel for filtering, e.g. postgres.GetBook(ctx, sq.Eq{"id": id}) to get book by id.
func GetBook(ctx context.Context, filter any) (*Book, error) {
	query, params, err := psql.Select(bookColumns...).From("books").ToSql()
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

// FindBooks returns books by given filter.
//
// Use squirrel for filtering, e.g. postgres.FindBooks(ctx, sq.Eq{"published": false}) to get books that aren't published yet.
func FindBooks(ctx context.Context, filter any) ([]*Book, error) {
	q := psql.Select(bookColumns...).From("books")
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
