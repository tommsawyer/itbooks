package postgres

import (
	"testing"

	"github.com/matryer/is"
)

func TestUpsertBookInsertUnexistingBook(t *testing.T) {
	ctx, is, rollback := runAndRollback(t)
	defer rollback()

	params := UpsertBookParams{
		ISBN:        "isbn",
		URL:         "url",
		Title:       "title",
		Image:       "image",
		Description: "description",
		Authors:     []string{"author"},
		Publisher:   "publisher",
		Properties: map[string]string{
			"test": "test",
		},
	}
	id, err := UpsertBook(ctx, params)
	is.NoErr(err) // we can create book

	book, err := GetBook(ctx, id)
	is.NoErr(err) // we should be able to get created book by id

	assertBookFieldsMatch(is, id, book, params)
}

func TestUpsertBookWithTheSameISBNUpdatesBook(t *testing.T) {
	ctx, is, rollback := runAndRollback(t)
	defer rollback()

	params := UpsertBookParams{
		ISBN:        "isbn",
		URL:         "url",
		Title:       "title",
		Image:       "image",
		Description: "description",
		Authors:     []string{"author"},
		Publisher:   "publisher",
		Properties: map[string]string{
			"test": "test",
		},
	}
	oldID, err := UpsertBook(ctx, params)
	is.NoErr(err) // we can create book

	updatedParams := UpsertBookParams{
		ISBN:        "isbn",
		URL:         "url2",
		Title:       "title2",
		Image:       "image2",
		Description: "description2",
		Authors:     []string{"author2"},
		Publisher:   "publisher2",
		Properties: map[string]string{
			"test": "test",
		},
	}
	newID, err := UpsertBook(ctx, updatedParams)
	is.NoErr(err)

	is.Equal(oldID, newID) // book should be updated, not created

	book, err := GetBook(ctx, newID)
	is.NoErr(err) // we should be able to get updated book by id

	assertBookFieldsMatch(is, newID, book, updatedParams)
}

func TestFindUnpublishedBooks(t *testing.T) {
	ctx, is, rollback := runAndRollback(t)
	defer rollback()

	publishedBook := UpsertBookParams{
		ISBN:  "isbn",
		Title: "published",
	}
	publishedID, err := UpsertBook(ctx, publishedBook)
	is.NoErr(err)

	is.NoErr(SetBookPublished(ctx, publishedID, true))

	unpublishedBook := UpsertBookParams{
		ISBN:  "isbn2",
		Title: "unpublished",
	}
	id, err := UpsertBook(ctx, unpublishedBook)
	is.NoErr(err)

	books, err := FindUnpublishedBooks(ctx)
	is.NoErr(err)

	is.Equal(len(books), 1)
	assertBookFieldsMatch(is, id, books[0], unpublishedBook)
}

func assertBookFieldsMatch(is *is.I, id int64, book *Book, params UpsertBookParams) {
	is.Helper()

	is.Equal(book.ID, id)
	is.Equal(book.ISBN.String, params.ISBN)
	is.Equal(book.URL.String, params.URL)
	is.Equal(book.Title.String, params.Title)
	is.Equal(book.Image.String, params.Image)
	is.Equal(book.Description.String, params.Description)
	is.Equal(len(book.Authors.Elements), len(params.Authors))
	for i, e := range book.Authors.Elements {
		is.Equal(e.String, params.Authors[i])
	}
	is.Equal(book.Publisher.String, params.Publisher)
	is.Equal(book.Properties, params.Properties)
}
