package products

import (
	"context"
	"time"
)

// Category entity
type Category struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ParentID    *int64     `json:"parent_id,omitempty"`
	EmpresaID   int64      `json:"empresa_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Brand entity
type Brand struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	EmpresaID   int64      `json:"empresa_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Product entity
type Product struct {
	ID           int64      `json:"id"`
	SKU          string     `json:"sku"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	CategoryID   *int64     `json:"category_id,omitempty"`
	BrandID      *int64     `json:"brand_id,omitempty"`
	CostPrice    float64    `json:"cost_price"`
	SalePrice    float64    `json:"sale_price"`
	TaxRate      float64    `json:"tax_rate"`
	MinStock     int        `json:"min_stock"`
	CurrentStock int        `json:"current_stock"`
	ImageURL     string     `json:"image_url"`
	Status       string     `json:"status"`
	EmpresaID    int64      `json:"empresa_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// Repository interfaces
type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context, empresaID int64) ([]Category, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id int64) error
}

type BrandRepository interface {
	Create(ctx context.Context, b *Brand) error
	GetByID(ctx context.Context, id int64) (*Brand, error)
	List(ctx context.Context, empresaID int64) ([]Brand, error)
	Update(ctx context.Context, b *Brand) error
	Delete(ctx context.Context, id int64) error
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetBySKU(ctx context.Context, sku string, empresaID int64) (*Product, error)
	List(ctx context.Context, empresaID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
}

// Service interfaces
type CategoryService interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id int64) (*Category, error)
	List(ctx context.Context, empresaID int64) ([]Category, error)
	Update(ctx context.Context, id int64, c *Category) error
	Delete(ctx context.Context, id int64) error
}

type BrandService interface {
	Create(ctx context.Context, b *Brand) error
	GetByID(ctx context.Context, id int64) (*Brand, error)
	List(ctx context.Context, empresaID int64) ([]Brand, error)
	Update(ctx context.Context, id int64, b *Brand) error
	Delete(ctx context.Context, id int64) error
}

type ProductService interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, empresaID int64, filters ListFilters) ([]Product, error)
	Update(ctx context.Context, id int64, p *Product) error
	Delete(ctx context.Context, id int64) error
}

// List filters
type ListFilters struct {
	Page       int
	Limit      int
	Search     string
	CategoryID *int64
	BrandID    *int64
}
