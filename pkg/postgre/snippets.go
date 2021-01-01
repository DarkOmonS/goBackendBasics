package postgre

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx"
	"module1/pkg/models"
)

type SnippetModel struct {
	Conn *pgx.Conn
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
			 VALUES($1, $2, current_timestamp, current_timestamp+$3*interval'1 day') returning id`
	res := m.Conn.QueryRow(stmt, title, content, expires)
	var id int
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > current_timestamp AND id = $1`
	s := &models.Snippet{}
	err := m.Conn.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`
	rows, err := m.Conn.Query(stmt)
	defer rows.Close()
	var snippets []*models.Snippet
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
