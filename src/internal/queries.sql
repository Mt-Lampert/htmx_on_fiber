
-- name: GetAllContacts :many
SELECT *
FROM contacts 
ORDER BY last_name, first_name;

-- name: GetContact :one
SELECT *
FROM contacts
WHERE id=? ;

-- name: SearchContacts :many
SELECT *
FROM contacts
WHERE 0
OR first_name LIKE ?
OR last_name LIKE ?
OR phone LIKE ?
OR email LIKE ? ;

-- name: AddContact :one
INSERT OR IGNORE INTO contacts
(first_name, last_name, phone, email)
VALUES (?, ?, ?, ?) 
RETURNING * ;

-- name: UpdateContact :one
UPDATE contacts
SET first_name=?, last_name=?, phone=?, email=?
WHERE id=? 
RETURNING * ;

-- name: DeleteContact :exec
DELETE FROM contacts WHERE id=? ;

