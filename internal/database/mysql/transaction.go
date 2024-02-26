package mysql

import (
	"context"
	"database/sql"
	"errors"
)

type UserStore struct {
	*Queries
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		Queries: New(db),
		db:      db,
	}
}

// AddUserIfNotExists is a return active registration users
func (us *UserStore) AddUserIfNotExists(ctx context.Context, param AddUserParams) error {
	tx, err := us.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := us.WithTx(tx)
	_, err = qtx.FindRegistrationUser(ctx, param.Email)
	if errors.Is(err, sql.ErrNoRows) {
		if err := qtx.AddUser(ctx, param); err != nil {
			return err
		}
		return tx.Commit()
	}
	return err
}

type UserFindStore struct {
	*Queries
	db *sql.DB
}

func NewUserFindStore(db *sql.DB) *UserFindStore {
	return &UserFindStore{
		Queries: New(db),
		db:      db,
	}
}
