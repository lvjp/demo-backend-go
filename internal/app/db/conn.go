package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Connector interface {
	UserDAO() UserDAO

	Close(context.Context) error
}

type connector struct {
	conn *pgx.Conn

	userDAO UserDAO
}

func NewConnector(ctx context.Context) (Connector, error) {
	config, err := pgx.ParseConfig("postgres://db?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("could not parse database config: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &connector{
			conn:    conn,
			userDAO: &userDAOImpl{conn: conn},
		},
		nil
}

func (c *connector) UserDAO() UserDAO {
	return c.userDAO
}

func (c *connector) Close(ctx context.Context) error {
	if err := c.conn.Close(ctx); err != nil {
		return fmt.Errorf("could not close database connection: %w", err)
	}

	return nil
}
