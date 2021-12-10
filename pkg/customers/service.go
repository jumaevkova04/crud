package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound        = errors.New("item not found")
	ErrInternal        = errors.New("internal error")
	ErrNoSuchUser      = errors.New("no such user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrExpired         = errors.New("expired")
)

// Service ...
type Service struct {
	pool *pgxpool.Pool
}

// NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Customer ...
type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

// CustomersToken ...
type CustomersToken struct {
	Token      string    `json:"token"`
	CustomerID int64     `json:"customer_id"`
	Expire     time.Time `json:"expire"`
	Created    time.Time `json:"created"`
}

// All ...
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.pool.Query(ctx, `
	SELECT id, name, phone, active, created FROM customers
	`)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}

		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Println("ERROR", err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Println("ERROR", err)
		return nil, ErrInternal
	}

	log.Println("items:", items)

	return items, nil
}

// All ...
func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.pool.Query(ctx, `
	SELECT id, name, phone, active, created FROM customers WHERE active
	`)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}

		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Println("ERROR", err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Println("ERROR", err)
		return nil, ErrInternal
	}

	return items, nil
}

// ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
	SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println("ERROR", err)
		return nil, ErrInternal
	}

	log.Printf("item: %#v", item)

	return item, nil
}

// Save ...
func (s *Service) Save(ctx context.Context, item *Customer) (*Customer, error) {
	var customer = &Customer{}

	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
		// return nil, ErrInternal
	}

	password := string(hash)

	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
		INSERT INTO customers (name, phone, password) VALUES($1, $2, $3) RETURNING id, name, phone, active, created
		`, item.Name, item.Phone, password).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		if err != nil {
			log.Println("ERROR", err)
			return nil, ErrInternal
		}
	} else {
		err := s.pool.QueryRow(ctx, `
		UPDATE customers SET name = $1, phone = $2,  password = $3 WHERE id = $4 RETURNING id, name, phone, active, created
		`, item.Name, item.Phone, password, item.ID).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		if err != nil {
			log.Println("ERROR", err)
			return nil, ErrInternal
		}
	}

	return customer, nil
}

// RemoveByID ...
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	var customer = &Customer{}

	err := s.pool.QueryRow(ctx, `
	DELETE FROM customers WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println("ERROR", err)
		return nil, ErrInternal
	}

	return customer, nil
}

// BlockByID ...
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	var customer = &Customer{}

	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = $2 WHERE id = $1 RETURNING id, name, phone, active, created
		`, id, false).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println("ERROR", err)
		return nil, ErrInternal
	}

	return customer, nil
}

// UnblockByID ...
func (s *Service) UnblockByID(ctx context.Context, id int64) (*Customer, error) {
	var customer = &Customer{}

	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = $2 WHERE id = $1 RETURNING id, name, phone, active, created
		`, id, true).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println("ERROR", err)
		return nil, ErrInternal
	}

	return customer, nil
}

// TokenForCustomer ...
func (s *Service) TokenForCustomer(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64

	err = s.pool.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone = $1`, phone).Scan(&id, &hash)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoSuchUser
	}

	if err != nil {
		log.Println("ERROR", err)
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("ERROR", err)
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		log.Println("ERROR", err)
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO customers_tokens (token, customer_id) VALUES ($1, $2)`, token, id)
	if err != nil {
		log.Println("ERROR", err)
		return "", ErrInternal
	}

	return token, nil
}

// AuthenticateCustomer ...
func (s *Service) AuthenticateCustomer(ctx context.Context, token string) (id int64, err error) {
	var expire time.Time
	err = s.pool.QueryRow(ctx, `SELECT customer_id, expire FROM customers_tokens WHERE token = $1`, token).Scan(&id, &expire)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("ERROR", err)
		return 0, ErrNoSuchUser
	}

	if err != nil {
		log.Println("ERROR", err)
		return 0, ErrInternal
	}

	if IsTimePassed(time.Now().Add(-time.Hour), expire) {
		return 0, ErrExpired
	}

	return id, nil
}

// IsTimePassed - if time passed returns true
func IsTimePassed(check, date time.Time) bool {
	return check.After(date)
}
