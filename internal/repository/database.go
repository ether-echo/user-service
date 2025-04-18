package repository

import (
	"database/sql"
	"fmt"
	"github.com/ether-echo/user-service/internal/domain"
	"github.com/ether-echo/user-service/pkg/config"
	"github.com/ether-echo/user-service/pkg/logger"
	"time"

	_ "github.com/lib/pq"
)

var (
	log = logger.Logger().Named("repository").Sugar()
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cnf *config.Config) *PostgresDB {
	db, err := sql.Open(cnf.DBName, fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", cnf.DBUser, cnf.DBPassword, cnf.DBName, cnf.DBHost))
	if err != nil {
		log.Fatalf("couldn't connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	return &PostgresDB{db: db}
}

func (p *PostgresDB) IsUserRegistered(chatID int64) (bool, error) {
	var exists bool

	err := p.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE chat_id = $1)
	`, chatID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}

	return exists, nil
}

func (p *PostgresDB) RegisterUser(user *domain.User) error {

	_, err := p.db.Exec(`
		INSERT INTO users (chat_id, first_name, last_name, username, registered_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id) DO NOTHING
	`,
		user.ChatId,
		user.FirstName,
		user.LastName,
		user.Username,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error registering user: %v", err)
	}

	return nil
}

func (p *PostgresDB) SaveMessage(chatID int64, message string) error {
	_, err := p.db.Exec(`
		INSERT INTO messages (user_id, message, created_at)
		VALUES ($1, $2, $3)
	`,
		chatID,
		message,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("error inserting message into database: %v", err)
	}

	return nil
}

func (p *PostgresDB) Close() {
	err := p.db.Close()
	if err != nil {
		log.Errorf("unable to close the database connection: %v", err)
	}
}
