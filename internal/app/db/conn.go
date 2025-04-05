package db

import "database/sql"

type Connector interface {
	UserDAO() UserDAO

	Close() error
}

type connector struct {
	db *sql.DB

	userDAO UserDAO
}

func NewConnector() (Connector, error) {
	db, err := sql.Open("postgres", "postgres://db?sslmode=disable")
	if err != nil {
		return nil, err
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
	return c.db.Close()
}
