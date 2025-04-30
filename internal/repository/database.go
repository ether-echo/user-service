package repository

import (
	"context"
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

func (p *PostgresDB) SaveMessage(ctx context.Context, chatID int64, message string) error {
	_, err := p.db.ExecContext(ctx, `
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

	log.Info("successfully inserted message into database")

	return nil
}

func (p *PostgresDB) GetTaro(ctx context.Context, chatID int64) (bool, error) {
	var taro bool

	err := p.db.QueryRowContext(ctx, `SELECT got_taro FROM users WHERE chat_id = $1`, chatID).Scan(&taro)
	if err != nil {
		return false, fmt.Errorf("error getting taro: %v", err)
	}

	return taro, nil
}

func (p *PostgresDB) ChangeAccessTaro(ctx context.Context, chatID int64) error {
	_, err := p.db.ExecContext(ctx, `UPDATE users SET got_taro = TRUE WHERE chat_id = $1`, chatID)
	if err != nil {
		return fmt.Errorf("error updating access taro user: %v", err)
	}

	log.Info("successfully updated access taro user")

	return nil
}

func (p *PostgresDB) ResetFlags(ctx context.Context) error {
	_, err := p.db.ExecContext(ctx, `
		UPDATE users
		SET
    		got_taro = FALSE,
    		got_numerology = FALSE
		WHERE
    		got_taro = TRUE OR got_numerology = TRUE;
	`)
	if err != nil {
		return fmt.Errorf("error resetting flags: %v", err)
	}

	return nil
}

func (p *PostgresDB) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT chat_id, first_name, last_name, username FROM users")
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %v", err)
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ChatId, &u.FirstName, &u.LastName, &u.Username); err != nil {
			return nil, fmt.Errorf("error scanning all users: %v", err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (p *PostgresDB) GetAllChatId(ctx context.Context) ([]int64, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT chat_id FROM users")
	if err != nil {
		return nil, fmt.Errorf("error getting all chatId: %v", err)
	}
	defer rows.Close()

	var chatIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning all chatId: %v", err)
		}
		chatIDs = append(chatIDs, id)
	}
	return chatIDs, nil
}

func (p *PostgresDB) Close() {
	err := p.db.Close()
	if err != nil {
		log.Errorf("unable to close the database connection: %v", err)
	}
}
