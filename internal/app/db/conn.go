package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.lvjp.me/demo-backend-go/internal/app/config"
)

type Connector interface {
	UserDAO() UserDAO

	Close() error
}

type connector struct {
	db *sql.DB

	userDAO UserDAO
}

func NewConnector(ctx context.Context, cfg config.Database) (Connector, error) {
	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("db connector creation: %w", err)
	}

	if pingErr := db.PingContext(ctx); pingErr != nil {
		retErr := fmt.Errorf("db connection ping: %w", pingErr)
		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(closeErr, fmt.Errorf("could not close db: %w", closeErr))
		}

		return nil, retErr
	}

	return &connector{
			db:      db,
			userDAO: &userDAOImpl{db: db},
		},
		nil
}

func (c *connector) UserDAO() UserDAO {
	return c.userDAO
}

func (c *connector) Close() error {
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("closing db connection: %w", err)
	}

	return nil
}
