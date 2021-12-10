package security

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Service ...
type Service struct {
	pool *pgxpool.Pool
}

// NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Manager ...
type Manager struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	Salary     int64     `json:"salary"`
	Plan       int64     `json:"plan"`
	BossID     int64     `json:"boss_id"`
	Department string    `json:"department"`
	Active     bool      `json:"active"`
	Created    time.Time `json:"created"`
}

// Auth ...
func (s *Service) Auth(login, password string) (ok bool) {
	manager := &Manager{}

	err := s.pool.QueryRow(context.Background(), `
	SELECT id FROM managers WHERE login = $1 AND password = $2
	`, login, password).Scan(&manager.ID)
	if err != nil {
		log.Println("ERROR", err)
		return false
	}

	return true
}
