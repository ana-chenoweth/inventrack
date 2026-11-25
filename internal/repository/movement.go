package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"inventrack/internal/models"
)

type MovementRepo struct {
	db *pgxpool.Pool
}

func NewMovementRepo(db *pgxpool.Pool) *MovementRepo {
	return &MovementRepo{db: db}
}

// Create registra un movimiento y ajusta el stock en UNA transacción.
//
// La regla "no sacar más de lo que hay" se aplica dentro del propio UPDATE
// (AND quantity + delta >= 0). Así, dos salidas simultáneas no pueden dejar
// el stock en negativo: la base de datos serializa el acceso a la fila.
func (r *MovementRepo) Create(ctx context.Context, m *models.Movement) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op si ya se hizo Commit

	delta := m.Quantity
	if m.Type == models.MovementOut {
		delta = -m.Quantity
	}

	tag, err := tx.Exec(ctx,
		`UPDATE products
		    SET quantity = quantity + $1::int
		  WHERE id = $2 AND quantity + $1::int >= 0`,
		delta, m.ProductID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		// O el producto no existe, o no hay stock suficiente.
		var exists bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM products WHERE id = $1)`, m.ProductID).
			Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return ErrNotFound
		}
		return ErrInsufficientStock
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO movements (product_id, type, quantity)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		m.ProductID, m.Type, m.Quantity).
		Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// List devuelve los movimientos, del más reciente al más antiguo.
// Si productID > 0 filtra por ese producto.
func (r *MovementRepo) List(ctx context.Context, productID int64) ([]models.Movement, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, product_id, type, quantity, created_at
		   FROM movements
		  WHERE ($1::bigint = 0 OR product_id = $1::bigint)
		  ORDER BY id DESC`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.Movement{}
	for rows.Next() {
		var m models.Movement
		if err := rows.Scan(&m.ID, &m.ProductID, &m.Type, &m.Quantity, &m.CreatedAt); err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}
	return movements, rows.Err()
}
