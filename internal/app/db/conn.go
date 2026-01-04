package db

import (
	"context"
	"fmt"

	"go.lvjp.me/demo-backend-go/internal/app/config"

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

func NewConnector(ctx context.Context, appCfg config.Database) (Connector, error) {
	dbCfg, err := loadDatabaseConfig(appCfg)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, dbCfg)
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

func loadDatabaseConfig(appCfg config.Database) (*pgx.ConnConfig, error) {
	dbCfg, err := pgx.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("could not parse database config: %w", err)
	}

	if appCfg.Host != nil {
		dbCfg.Host = *appCfg.Host
	}
	if appCfg.Port != nil {
		dbCfg.Port = *appCfg.Port
	}
	if appCfg.Database != nil {
		dbCfg.Database = *appCfg.Database
	}
	if appCfg.User != nil {
		dbCfg.User = *appCfg.User
	}
	if appCfg.Password != nil {
		dbCfg.Password = *appCfg.Password
	}

	return dbCfg, nil
}
