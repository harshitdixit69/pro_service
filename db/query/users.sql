-- name: Createusers :exec
INSERT INTO users (username, first_name, last_name, pre_password, password, encrypted_passwrod, address,
phone_no, aadhaar_no, status, is_active) 
VALUES(? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,?);

-- name: Getusers :one
SELECT * FROM users where id = ? OR username = ? OR phone_no = ?;

-- name: GetListusers :many
SELECT * FROM users LIMIT ? OFFSET ?;

-- name: Updateusers :exec
UPDATE users SET username = ? , first_name = ? , last_name = ? , pre_password = ? , 
password = ?,encrypted_passwrod = ? , phone_no = ? WHERE username =? ;

-- name: GetUpdatedAccount :one
SELECT * FROM users WHERE id = ?;

-- name: DeleteusersById :exec 
Delete FROM users where id = ?;

