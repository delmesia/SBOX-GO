package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	// This will initialize an instance of Snippet struct with zero value.
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {

	// SQL Statement that will be executed
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing
	// the result of the query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns.
	// this defer statement should come *after* the error check from the Query() method.
	// otherwise, if Query() returns an error, it'll panic trying to close a nil resultset.

	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// rows.Next() is used to iterate through the rows in the resultset.
	// this prepares the first (and then each subsequent) row to be acted on by
	// the row.Scan() method. if iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {

		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}

		// row.Scan() will copy the values from each field in a the row to the
		// new Snippet object created. The arguments to row.Scan() must be pointers to the place we want
		// to copy the data into. The number of arguments must be exactly the same as the number
		// of columns returned by the statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append each Snippet got from row.Scan after putting it to S.
		snippets = append(snippets, s)
	}

	// when the rows.Next() loop has finished, we call rows.Err() to retrieve any error
	// that was encountered during the iteration. It's important to call this,
	// don't assume a successful iteration was completed over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}
	//if everyting went okay, return the Snippets slice and nil to indicate that it's successful
	return snippets, nil

}
