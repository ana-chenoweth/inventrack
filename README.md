# Inventrack

API REST de gestión de inventario escrita en **Go** con **PostgreSQL**. Permite administrar productos y registrar movimientos de entrada/salida de existencias, con validaciones y transacciones SQL.

## Demo
https://inventrack-t0tr.onrender.com/api/products

> Plan gratis de Render: la primera petición después de un rato sin uso puede tardar ~1 minuto.

## Tecnologías

- Go (`net/http`, enrutador estándar con métodos y parámetros de ruta)
- PostgreSQL 16 + [pgx](https://github.com/jackc/pgx) (pool de conexiones)
- Docker Compose (base de datos local)

## Cómo ejecutarlo

Requisitos: Go 1.22+ y Docker.

```bash
docker compose up -d          # levanta PostgreSQL y crea las tablas
go run ./cmd/api              # API en http://localhost:8080
```

Variables de entorno opcionales: `DATABASE_URL`, `PORT`.

## Endpoints

| Método | Ruta                  | Descripción                                   |
|--------|-----------------------|-----------------------------------------------|
| GET    | `/api/products`       | Listar productos                              |
| POST   | `/api/products`       | Crear producto                                |
| GET    | `/api/products/{id}`  | Consultar producto                            |
| PUT    | `/api/products/{id}`  | Actualizar nombre, SKU y descripción          |
| DELETE | `/api/products/{id}`  | Eliminar producto                             |
| POST   | `/api/movements`      | Registrar entrada (`IN`) o salida (`OUT`)     |
| GET    | `/api/movements`      | Listar movimientos (`?product_id=1` opcional) |

### Ejemplos

```bash
# Crear producto
curl -X POST localhost:8080/api/products \
  -d '{"name":"Tornillo M4","sku":"TOR-M4","description":"Caja de 100","quantity":10}'

# Registrar una salida de 4 unidades
curl -X POST localhost:8080/api/movements \
  -d '{"product_id":1,"type":"OUT","quantity":4}'

# Salida mayor al stock -> 409 Conflict
curl -i -X POST localhost:8080/api/movements \
  -d '{"product_id":1,"type":"OUT","quantity":999}'
```

### Códigos de respuesta

| Código | Cuándo                                              |
|--------|-----------------------------------------------------|
| 200/201/204 | Operación exitosa                              |
| 400    | JSON inválido, id inválido                          |
| 404    | Producto no existe                                  |
| 409    | SKU duplicado o stock insuficiente                  |
| 422    | Validación fallida (campo vacío, cantidad <= 0...)  |
| 500    | Error interno                                       |

## Decisiones de diseño

- **Capas**: `handlers` (HTTP) → `repository` (SQL) → `models`. Los handlers no escriben SQL.
- **Stock consistente**: el stock solo cambia mediante movimientos. Cada movimiento actualiza el producto y guarda el registro en **una sola transacción**.
- **Sin stock negativo, incluso con concurrencia**: la regla se aplica dentro del `UPDATE` (`AND quantity + delta >= 0`) y además existe un `CHECK (quantity >= 0)` en la base de datos.
- **Errores tipados**: el repositorio devuelve errores como `ErrNotFound` o `ErrInsufficientStock` y los handlers los convierten en códigos HTTP.

## Pruebas

```bash
go test ./...
```

## Posibles mejoras

Paginación, autenticación, tests de integración con una base de datos real, Dockerfile para la API.
