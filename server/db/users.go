package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"

	errs "server/errors"
)

type UserEntity struct {
	UserId    string
	Email     string
	FirstName string
	LastName  string
	Age       int32
	Aliases   []string
}

func GetUsers(db *sql.DB, ctx context.Context) ([]UserEntity, error) {
	var users []UserEntity

	rows, err := db.QueryContext(
		ctx,
		`SELECT
			user_id,
			email,
			first_name,
			last_name,
			age,
			aliases
		FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserEntity
		if err := rows.Scan(
			&user.UserId, &user.Email, &user.FirstName, &user.LastName, &user.Age, pq.Array(&user.Aliases),
		); err != nil {
			return nil, errs.NewHTTPError("Failed to retrieve users", http.StatusInternalServerError, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, errs.NewHTTPError("Failed to retrieve users", http.StatusInternalServerError, err)
	}
	return users, nil
}

func GetUserById(db *sql.DB, ctx context.Context, userId string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`SELECT
			user_id,
			email,
			first_name,
			last_name,
			age,
			aliases
		FROM users
		WHERE user_id = $1`,
		userId,
	).Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		pq.Array(&user.Aliases),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userIdNotFoundError(userId)
		}

		return nil, errs.NewHTTPError("Failed to retrieve user", http.StatusInternalServerError, err)
	}

	return &user, nil
}

type UpdateUserPayload struct {
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Age       int32     `json:"age"`
	Aliases   *[]string `json:"aliases"`
}

func CreateUser(db *sql.DB, ctx context.Context, data UpdateUserPayload) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, first_name, last_name, age, aliases)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			user_id,
			email,
			first_name,
			last_name,
			age,
			aliases`,
		data.Email,
		data.FirstName,
		data.LastName,
		data.Age,
		pq.Array(data.Aliases),
	).Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		pq.Array(&user.Aliases),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Duplicate key error
				return nil, errs.NewHTTPError(
					fmt.Sprintf("Duplicate key violation for constraint `%s`", pgErr.ConstraintName),
					http.StatusConflict,
					nil,
				)
			}
		}

		return nil, errs.NewHTTPError("Failed to save user to db", http.StatusInternalServerError, err)
	}

	return &user, nil
}

func UpdateUser(db *sql.DB, ctx context.Context, userId string, data UpdateUserPayload) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`UPDATE users
		SET email = $1,
			first_name = $2,
			last_name = $3,
			age = $4,
			aliases = $5
		WHERE user_id = $6
		RETURNING
			user_id,
			email,
			first_name,
			last_name,
			age,
			aliases`,
		data.Email,
		data.FirstName,
		data.LastName,
		data.Age,
		pq.Array(data.Aliases),
		userId,
	).Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		pq.Array(&user.Aliases),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, userIdNotFoundError(userId)
		}

		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				// Duplicate key error
				return nil, errs.NewHTTPError(
					fmt.Sprintf("Duplicate key violation for constraint `%s`", pgError.ConstraintName),
					http.StatusConflict,
					nil,
				)
			}
		}

		return nil, errs.NewHTTPError("Could not update user", http.StatusInternalServerError, err)
	}

	return &user, nil
}

func DeleteUser(db *sql.DB, ctx context.Context, userId string) error {
	_, err := db.ExecContext(
		ctx,
		`DELETE FROM users WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return errs.NewHTTPError(
			fmt.Sprintf("Could not delete user %s", userId),
			http.StatusInternalServerError,
			err,
		)
	}

	return nil
}

func userIdNotFoundError(userId string) error {
	return errs.NewHTTPError(fmt.Sprintf("No user with id %s", userId), http.StatusNotFound, nil)
}
