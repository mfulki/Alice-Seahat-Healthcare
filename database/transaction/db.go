package transaction

import (
	"context"
	"database/sql"

	"Alice-Seahat-Healthcare/seahat-be/constant"
)

type DBTransaction interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	PrepareContext(context.Context, string, ...any) (*sql.Stmt, error)
}

type dbTransaction struct {
	conn *sql.DB
}

func NewDBTransaction(db *sql.DB) *dbTransaction {
	return &dbTransaction{
		conn: db,
	}
}

func (dt *dbTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, ok := ctx.Value(constant.TxContext).(*sql.Tx)
	if !ok {
		return dt.conn.ExecContext(ctx, query, args...)
	}

	return tx.ExecContext(ctx, query, args...)
}

func (dt *dbTransaction) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	tx, ok := ctx.Value(constant.TxContext).(*sql.Tx)
	if !ok {
		return dt.conn.QueryContext(ctx, query, args...)
	}

	return tx.QueryContext(ctx, query, args...)
}

func (dt *dbTransaction) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	tx, ok := ctx.Value(constant.TxContext).(*sql.Tx)
	if !ok {
		return dt.conn.QueryRowContext(ctx, query, args...)
	}

	return tx.QueryRowContext(ctx, query, args...)
}
func (dt *dbTransaction) PrepareContext(ctx context.Context, query string, args ...any) (*sql.Stmt, error) {
	tx, ok := ctx.Value(constant.TxContext).(*sql.Tx)
	if !ok {
		return dt.conn.PrepareContext(ctx, query)
	}
	return tx.PrepareContext(ctx, query)

}
