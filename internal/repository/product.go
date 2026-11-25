package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"inventrack/internal/models"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrDuplicateSKU      = errors.New("sku already exists")
	ErrInsufficientStock = errors.New("insufficient stock")
)

// pgUniqueViolation es el código de error de PostgreSQL para UNIQUE.
const pgUniqueViolation = "23505"

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) List(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, sku, description, quantity, created_at
		   FROM products
		  ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{} // slice vacío (no nil) para que el JSON sea [] y no null
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Description, &p.Quantity, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepo) Get(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(ctx,
		`SELECT id, name, sku, description, quantity, created_at
		   FROM products
		  WHERE id = $1`, id).
		Scan(&p.ID, &p.Name, &p.SKU, &p.Description, &p.Quantity, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) Create(ctx context.Context, p *models.Product) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO products (name, sku, description, quantity)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		p.Name, p.SKU, p.Description, p.Quantity).
		Scan(&p.ID, &p.CreatedAt)
	return mapWriteErr(err)
}

// Update cambia nombre, sku y descripción. El stock NO se toca aquí:
// solo cambia con movimientos.
func (r *ProductRepo) Update(ctx context.Context, p *models.Product) error {
	err := r.db.QueryRow(ctx,
		`UPDATE products
		    SET name = $1, sku = $2, description = $3
		  WHERE id = $4
		 RETURNING quantity, created_at`,
		p.Name, p.SKU, p.Description, p.ID).
		Scan(&p.Quantity, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return mapWriteErr(err)
}

func (r *ProductRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func mapWriteErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return ErrDuplicateSKU
	}
	return err
}
