package customers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrNotFound = errors.New("item not found")
	ErrInternal = errors.New("internal error")
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
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
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

	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
		INSERT INTO customers (name, phone) VALUES($1, $2) RETURNING id, name, phone, active, created
		`, item.Name, item.Phone).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		if err != nil {
			log.Println("ERROR", err)
			return nil, ErrInternal
		}
	} else {
		err := s.pool.QueryRow(ctx, `
		UPDATE customers SET name = $1, phone = $2 WHERE id = $3 RETURNING id, name, phone, active, created
		`, item.Name, item.Phone, item.ID).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

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
