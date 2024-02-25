// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
)

const addContact = `-- name: AddContact :one
;

INSERT OR IGNORE INTO contacts
(first_name, last_name, phone, email)
VALUES (?, ?, ?, ?) 
RETURNING id, last_name, first_name, phone, email
`

type AddContactParams struct {
	FirstName sql.NullString
	LastName  sql.NullString
	Phone     sql.NullString
	Email     sql.NullString
}

func (q *Queries) AddContact(ctx context.Context, arg AddContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, addContact,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
		arg.Email,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.LastName,
		&i.FirstName,
		&i.Phone,
		&i.Email,
	)
	return i, err
}

const deleteContact = `-- name: DeleteContact :exec
;

DELETE FROM contacts WHERE id=?
`

func (q *Queries) DeleteContact(ctx context.Context, id sql.NullInt64) error {
	_, err := q.db.ExecContext(ctx, deleteContact, id)
	return err
}

const getContact = `-- name: GetContact :one
SELECT id, last_name, first_name, phone, email
FROM contacts
WHERE id=?
`

func (q *Queries) GetContact(ctx context.Context, id sql.NullInt64) (Contact, error) {
	row := q.db.QueryRowContext(ctx, getContact, id)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.LastName,
		&i.FirstName,
		&i.Phone,
		&i.Email,
	)
	return i, err
}

const getContacts = `-- name: GetContacts :many
SELECT id, last_name, first_name, phone, email
FROM contacts 
ORDER BY last_name, first_name
LIMIT ?
`

func (q *Queries) GetContacts(ctx context.Context, limit int64) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, getContacts, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.LastName,
			&i.FirstName,
			&i.Phone,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEmail = `-- name: GetEmail :one
;

SELECT id FROM contacts WHERE email=?
`

func (q *Queries) GetEmail(ctx context.Context, email sql.NullString) (sql.NullInt64, error) {
	row := q.db.QueryRowContext(ctx, getEmail, email)
	var id sql.NullInt64
	err := row.Scan(&id)
	return id, err
}

const searchContacts = `-- name: SearchContacts :many
;

SELECT id, last_name, first_name, phone, email
FROM contacts
WHERE 0
OR first_name LIKE ?
OR last_name LIKE ?
OR phone LIKE ?
OR email LIKE ?
`

type SearchContactsParams struct {
	FirstName sql.NullString
	LastName  sql.NullString
	Phone     sql.NullString
	Email     sql.NullString
}

func (q *Queries) SearchContacts(ctx context.Context, arg SearchContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, searchContacts,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
		arg.Email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contact
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.LastName,
			&i.FirstName,
			&i.Phone,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateContact = `-- name: UpdateContact :one
;

UPDATE contacts
SET first_name=?, last_name=?, phone=?, email=?
WHERE id=? 
RETURNING id, last_name, first_name, phone, email
`

type UpdateContactParams struct {
	FirstName sql.NullString
	LastName  sql.NullString
	Phone     sql.NullString
	Email     sql.NullString
	ID        sql.NullInt64
}

func (q *Queries) UpdateContact(ctx context.Context, arg UpdateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, updateContact,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
		arg.Email,
		arg.ID,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.LastName,
		&i.FirstName,
		&i.Phone,
		&i.Email,
	)
	return i, err
}

// vim: foldmethod=indent
