package models

import (
	"errors"
	"strings"
	"time"
)

const (
	MovementIn  = "IN"
	MovementOut = "OUT"
)

// Product es una fila de la tabla products.
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
}

// Movement es una fila de la tabla movements.
type Movement struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductInput es el cuerpo JSON de POST/PUT /api/products.
// Quantity solo se usa al crear: después el stock cambia únicamente con movimientos.
type ProductInput struct {
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// Validate limpia espacios y revisa las reglas de negocio.
func (in *ProductInput) Validate() error {
	in.Name = strings.TrimSpace(in.Name)
	in.SKU = strings.TrimSpace(in.SKU)
	in.Description = strings.TrimSpace(in.Description)

	switch {
	case in.Name == "":
		return errors.New("name is required")
	case in.SKU == "":
		return errors.New("sku is required")
	case in.Quantity < 0:
		return errors.New("quantity cannot be negative")
	}
	return nil
}

// MovementInput es el cuerpo JSON de POST /api/movements.
type MovementInput struct {
	ProductID int64  `json:"product_id"`
	Type      string `json:"type"`
	Quantity  int    `json:"quantity"`
}

// Validate normaliza el tipo a mayúsculas y revisa los campos.
func (in *MovementInput) Validate() error {
	in.Type = strings.ToUpper(strings.TrimSpace(in.Type))

	switch {
	case in.ProductID <= 0:
		return errors.New("product_id is required")
	case in.Type != MovementIn && in.Type != MovementOut:
		return errors.New("type must be IN or OUT")
	case in.Quantity <= 0:
		return errors.New("quantity must be greater than 0")
	}
	return nil
}
