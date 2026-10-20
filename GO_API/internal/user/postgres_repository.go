package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const userColumns = `id, name, email, birthday::text, cpf, inserted_at, updated_at`

type PostgresRepository struct {
	database *pgxpool.Pool
}

func NewPostgresRepository(database *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{database: database}
}

func (repository *PostgresRepository) List(ctx context.Context) ([]User, error) {
	query := `SELECT ` + userColumns + ` FROM users ORDER BY id`
	rows, err := repository.database.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repository *PostgresRepository) Get(ctx context.Context, id int64) (User, error) {
	query := `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	user, err := scanUser(repository.database.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (repository *PostgresRepository) Create(ctx context.Context, input Input) (User, error) {
	query := `INSERT INTO users (name, email, birthday, cpf)
		VALUES ($1, $2, $3, $4)
		RETURNING ` + userColumns
	return scanUser(repository.database.QueryRow(
		ctx,
		query,
		input.Name,
		input.Email,
		input.Birthday,
		input.CPF,
	))
}

func (repository *PostgresRepository) Update(
	ctx context.Context,
	id int64,
	input Input,
) (User, error) {
	query := `UPDATE users
		SET name = $1, email = $2, birthday = $3, cpf = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING ` + userColumns
	user, err := scanUser(repository.database.QueryRow(
		ctx,
		query,
		input.Name,
		input.Email,
		input.Birthday,
		input.CPF,
		id,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (repository *PostgresRepository) Delete(ctx context.Context, id int64) (bool, error) {
	result, err := repository.database.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

type rowScanner interface {
	Scan(...any) error
}

func scanUser(row rowScanner) (User, error) {
	var user User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Birthday,
		&user.CPF,
		&user.InsertedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func IsConflict(err error) bool {
	var databaseError *pgconn.PgError
	return errors.As(err, &databaseError) && databaseError.Code == "23505"
}

func IsInvalidDate(err error) bool {
	var databaseError *pgconn.PgError
	if !errors.As(err, &databaseError) {
		return false
	}
	return databaseError.Code == "22007" || databaseError.Code == "22008"
}
