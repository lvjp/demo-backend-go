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
		return nil, fmt.Errorf("db connector creation: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil,
			errors.Join(
				fmt.Errorf("db connection ping: %v", err),
				(&connector{db: db}).Close(),
			)
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
		return fmt.Errorf("closing db connection: %v", err)
	}

	return nil
}
