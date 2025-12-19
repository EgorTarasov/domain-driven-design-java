package persistance

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// DB represents a database connection interface
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// MinimalRepository provides basic CRUD operations using Squirrel
type MinimalRepository struct {
	db      DB
	builder sq.StatementBuilderType
}

// NewMinimalRepository creates a new repository instance
func NewMinimalRepository(db DB) *MinimalRepository {
	return &MinimalRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// WithTransaction returns a new repository using the given transaction
func (r *MinimalRepository) WithTransaction(tx DB) *MinimalRepository {
	return &MinimalRepository{
		db:      tx,
		builder: r.builder,
	}
}

// Insert executes an insert query built with Squirrel
func (r *MinimalRepository) Insert(ctx context.Context, builder sq.InsertBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	return r.db.ExecContext(ctx, query, args...)
}

// Update executes an update query built with Squirrel
func (r *MinimalRepository) Update(ctx context.Context, builder sq.UpdateBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	return r.db.ExecContext(ctx, query, args...)
}

// Delete executes a delete query built with Squirrel
func (r *MinimalRepository) Delete(ctx context.Context, builder sq.DeleteBuilder) (sql.Result, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	return r.db.ExecContext(ctx, query, args...)
}

// QueryRow executes a select query and returns a single row
func (r *MinimalRepository) QueryRow(ctx context.Context, builder sq.SelectBuilder) sq.RowScanner {
	query, args, err := builder.ToSql()
	if err != nil {
		return &errorScanner{err: err}
	}
	return r.db.QueryRowContext(ctx, query, args...)
}

// Query executes a select query and returns multiple rows
func (r *MinimalRepository) Query(ctx context.Context, builder sq.SelectBuilder) (*sql.Rows, error) {
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	return r.db.QueryContext(ctx, query, args...)
}

// GetBuilder returns the Squirrel statement builder
func (r *MinimalRepository) GetBuilder() sq.StatementBuilderType {
	return r.builder
}

// errorScanner implements RowScanner for error cases
type errorScanner struct {
	err error
}

func (e *errorScanner) Scan(dest ...interface{}) error {
	return e.err
}
