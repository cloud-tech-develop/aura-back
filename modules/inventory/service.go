package inventory

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type repository struct {
	db db.Querier
}

func NewRepository(db db.Querier) Repository {
	return &repository{db: db}
}

// Inventory operations

func (r *repository) CreateInventory(ctx context.Context, inv *Inventory) error {
	query := `
		INSERT INTO inventory (product_id, branch_id, quantity, reserved_quantity, min_stock, max_stock, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		inv.ProductID, inv.BranchID, inv.Quantity, inv.ReservedQuantity,
		inv.MinStock, inv.MaxStock, inv.Location,
	).Scan(&inv.ID, &inv.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

func (r *repository) GetInventory(ctx context.Context, productID, branchID int64) (*Inventory, error) {
	inv := &Inventory{}
	query := `
		SELECT id, product_id, branch_id, quantity, reserved_quantity, min_stock, max_stock, location, created_at, updated_at
		FROM inventory WHERE product_id = $1 AND branch_id = $2`

	err := r.db.QueryRowContext(ctx, query, productID, branchID).Scan(
		&inv.ID, &inv.ProductID, &inv.BranchID, &inv.Quantity, &inv.ReservedQuantity,
		&inv.MinStock, &inv.MaxStock, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return inv, nil
}

func (r *repository) UpdateInventory(ctx context.Context, inv *Inventory) error {
	query := `
		UPDATE inventory SET
			quantity = $1, reserved_quantity = $2, min_stock = $3, max_stock = $4, location = $5
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		inv.Quantity, inv.ReservedQuantity, inv.MinStock, inv.MaxStock, inv.Location, inv.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) ListInventory(ctx context.Context, enterpriseID int64, filters InventoryFilters) ([]Inventory, error) {
	query := `
		SELECT i.id, i.product_id, i.branch_id, i.quantity, i.reserved_quantity, i.min_stock, i.max_stock, i.location, i.created_at, i.updated_at
		FROM inventory i
		JOIN product p ON p.id = i.product_id
		WHERE p.enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND i.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	if filters.ProductID != nil {
		query += fmt.Sprintf(" AND i.product_id = $%d", argPos)
		args = append(args, *filters.ProductID)
		argPos++
	}

	if filters.LowStock {
		query += fmt.Sprintf(" AND i.quantity <= i.min_stock")
	}

	query += " ORDER BY i.product_id"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
		argPos++

		offset := (filters.Page - 1) * filters.Limit
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	defer rows.Close()

	var inventories []Inventory
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.BranchID, &inv.Quantity, &inv.ReservedQuantity,
			&inv.MinStock, &inv.MaxStock, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}
	return inventories, nil
}

func (r *repository) GetLowStock(ctx context.Context, enterpriseID int64) ([]Inventory, error) {
	query := `
		SELECT i.id, i.product_id, i.branch_id, i.quantity, i.reserved_quantity, i.min_stock, i.max_stock, i.location, i.created_at, i.updated_at
		FROM inventory i
		JOIN product p ON p.id = i.product_id
		WHERE p.enterprise_id = $1 AND i.quantity <= i.min_stock
		ORDER BY (i.quantity - i.min_stock) ASC`

	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock: %w", err)
	}
	defer rows.Close()

	var inventories []Inventory
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.BranchID, &inv.Quantity, &inv.ReservedQuantity,
			&inv.MinStock, &inv.MaxStock, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}
	return inventories, nil
}

func (r *repository) GetInventoryByProduct(ctx context.Context, productID int64) ([]Inventory, error) {
	query := `
		SELECT id, product_id, branch_id, quantity, reserved_quantity, min_stock, max_stock, location, created_at, updated_at
		FROM inventory WHERE product_id = $1
		ORDER BY branch_id`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by product: %w", err)
	}
	defer rows.Close()

	var inventories []Inventory
	for rows.Next() {
		var inv Inventory
		if err := rows.Scan(
			&inv.ID, &inv.ProductID, &inv.BranchID, &inv.Quantity, &inv.ReservedQuantity,
			&inv.MinStock, &inv.MaxStock, &inv.Location, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}
	return inventories, nil
}

// Movement operations

func (r *repository) CreateMovement(ctx context.Context, mov *InventoryMovement) error {
	query := `
		INSERT INTO inventory_movement (
			inventory_id, movement_type, movement_reason, quantity, previous_balance, new_balance,
			reference_id, reference_type, batch_number, serial_number, expiration_date, notes,
			user_id, branch_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		mov.InventoryID, mov.MovementType, mov.MovementReason, mov.Quantity,
		mov.PreviousBalance, mov.NewBalance, mov.ReferenceID, mov.ReferenceType,
		mov.BatchNumber, mov.SerialNumber, mov.ExpirationDate, mov.Notes,
		mov.UserID, mov.BranchID,
	).Scan(&mov.ID, &mov.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create movement: %w", err)
	}
	return nil
}

func (r *repository) GetMovement(ctx context.Context, id int64) (*InventoryMovement, error) {
	mov := &InventoryMovement{}
	query := `
		SELECT id, inventory_id, movement_type, movement_reason, quantity, previous_balance, new_balance,
			reference_id, reference_type, batch_number, serial_number, expiration_date, notes,
			user_id, branch_id, created_at
		FROM inventory_movement WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&mov.ID, &mov.InventoryID, &mov.MovementType, &mov.MovementReason, &mov.Quantity,
		&mov.PreviousBalance, &mov.NewBalance, &mov.ReferenceID, &mov.ReferenceType,
		&mov.BatchNumber, &mov.SerialNumber, &mov.ExpirationDate, &mov.Notes,
		&mov.UserID, &mov.BranchID, &mov.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get movement: %w", err)
	}
	return mov, nil
}

func (r *repository) ListMovements(ctx context.Context, enterpriseID int64, filters MovementFilters) ([]InventoryMovement, error) {
	query := `
		SELECT m.id, m.inventory_id, m.movement_type, m.movement_reason, m.quantity, m.previous_balance, m.new_balance,
			m.reference_id, m.reference_type, m.batch_number, m.serial_number, m.expiration_date, m.notes,
			m.user_id, m.branch_id, m.created_at
		FROM inventory_movement m
		JOIN inventory i ON i.id = m.inventory_id
		JOIN product p ON p.id = i.product_id
		WHERE p.enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.InventoryID != nil {
		query += fmt.Sprintf(" AND m.inventory_id = $%d", argPos)
		args = append(args, *filters.InventoryID)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND m.branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	if filters.ProductID != nil {
		query += fmt.Sprintf(" AND i.product_id = $%d", argPos)
		args = append(args, *filters.ProductID)
		argPos++
	}

	if filters.MovementType != "" {
		query += fmt.Sprintf(" AND m.movement_type = $%d", argPos)
		args = append(args, filters.MovementType)
		argPos++
	}

	if filters.MovementReason != "" {
		query += fmt.Sprintf(" AND m.movement_reason = $%d", argPos)
		args = append(args, filters.MovementReason)
		argPos++
	}

	query += " ORDER BY m.created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
		argPos++

		offset := (filters.Page - 1) * filters.Limit
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list movements: %w", err)
	}
	defer rows.Close()

	var movements []InventoryMovement
	for rows.Next() {
		var mov InventoryMovement
		if err := rows.Scan(
			&mov.ID, &mov.InventoryID, &mov.MovementType, &mov.MovementReason, &mov.Quantity,
			&mov.PreviousBalance, &mov.NewBalance, &mov.ReferenceID, &mov.ReferenceType,
			&mov.BatchNumber, &mov.SerialNumber, &mov.ExpirationDate, &mov.Notes,
			&mov.UserID, &mov.BranchID, &mov.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, mov)
	}
	return movements, nil
}

func (r *repository) GetProductKardex(ctx context.Context, productID int64, branchID int64) ([]InventoryMovement, error) {
	query := `
		SELECT m.id, m.inventory_id, m.movement_type, m.movement_reason, m.quantity, m.previous_balance, m.new_balance,
			m.reference_id, m.reference_type, m.batch_number, m.serial_number, m.expiration_date, m.notes,
			m.user_id, m.branch_id, m.created_at
		FROM inventory_movement m
		JOIN inventory i ON i.id = m.inventory_id
		WHERE i.product_id = $1 AND i.branch_id = $2
		ORDER BY m.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, productID, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get kardex: %w", err)
	}
	defer rows.Close()

	var movements []InventoryMovement
	for rows.Next() {
		var mov InventoryMovement
		if err := rows.Scan(
			&mov.ID, &mov.InventoryID, &mov.MovementType, &mov.MovementReason, &mov.Quantity,
			&mov.PreviousBalance, &mov.NewBalance, &mov.ReferenceID, &mov.ReferenceType,
			&mov.BatchNumber, &mov.SerialNumber, &mov.ExpirationDate, &mov.Notes,
			&mov.UserID, &mov.BranchID, &mov.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, mov)
	}
	return movements, nil
}

// Movement Reasons

func (r *repository) ListReasons(ctx context.Context) ([]MovementReason, error) {
	query := `SELECT id, code, name, description, movement_type, requires_authorization, is_active, created_at FROM movement_reason WHERE is_active = TRUE ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list reasons: %w", err)
	}
	defer rows.Close()

	var reasons []MovementReason
	for rows.Next() {
		var reason MovementReason
		if err := rows.Scan(
			&reason.ID, &reason.Code, &reason.Name, &reason.Description,
			&reason.MovementType, &reason.RequiresAuthorization, &reason.IsActive, &reason.CreatedAt,
		); err != nil {
			return nil, err
		}
		reasons = append(reasons, reason)
	}
	return reasons, nil
}

func (r *repository) GetReason(ctx context.Context, code string) (*MovementReason, error) {
	reason := &MovementReason{}
	query := `SELECT id, code, name, description, movement_type, requires_authorization, is_active, created_at FROM movement_reason WHERE code = $1`

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&reason.ID, &reason.Code, &reason.Name, &reason.Description,
		&reason.MovementType, &reason.RequiresAuthorization, &reason.IsActive, &reason.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get reason: %w", err)
	}
	return reason, nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db db.Querier) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) GetInventory(ctx context.Context, productID, branchID int64) (*Inventory, error) {
	return s.repo.GetInventory(ctx, productID, branchID)
}

func (s *service) ListInventory(ctx context.Context, enterpriseID int64, filters InventoryFilters) ([]Inventory, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.ListInventory(ctx, enterpriseID, filters)
}

func (s *service) GetLowStock(ctx context.Context, enterpriseID int64) ([]Inventory, error) {
	return s.repo.GetLowStock(ctx, enterpriseID)
}

func (s *service) GetInventoryByProduct(ctx context.Context, productID int64) ([]Inventory, error) {
	return s.repo.GetInventoryByProduct(ctx, productID)
}

func (s *service) UpdateStock(ctx context.Context, productID, branchID int64, quantity int, reason, referenceType string, referenceID *int64, userID int64, notes string) (*Inventory, error) {
	// Get current inventory
	inv, err := s.repo.GetInventory(ctx, productID, branchID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new inventory if not exists
			inv = &Inventory{
				ProductID:  productID,
				BranchID:   branchID,
				Quantity:   0,
				MinStock:   0,
			}
			if err := s.repo.CreateInventory(ctx, inv); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	previousBalance := inv.Quantity
	newBalance := previousBalance + quantity

	if newBalance < 0 {
		return nil, fmt.Errorf("insufficient stock: current %d, requested %d", previousBalance, quantity)
	}

	// Determine movement type
	movementType := MovementTypeEntry
	if quantity < 0 {
		movementType = MovementTypeExit
	}

	// Create movement record
	mov := &InventoryMovement{
		InventoryID:    inv.ID,
		MovementType:   movementType,
		MovementReason: reason,
		Quantity:       abs(quantity),
		PreviousBalance: previousBalance,
		NewBalance:     newBalance,
		ReferenceID:   referenceID,
		ReferenceType: referenceType,
		Notes:         notes,
		UserID:        userID,
		BranchID:      branchID,
	}

	if err := s.repo.CreateMovement(ctx, mov); err != nil {
		return nil, err
	}

	// Update inventory
	inv.Quantity = newBalance
	if err := s.repo.UpdateInventory(ctx, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

func (s *service) GetMovement(ctx context.Context, id int64) (*InventoryMovement, error) {
	return s.repo.GetMovement(ctx, id)
}

func (s *service) ListMovements(ctx context.Context, enterpriseID int64, filters MovementFilters) ([]InventoryMovement, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.ListMovements(ctx, enterpriseID, filters)
}

func (s *service) GetProductKardex(ctx context.Context, productID int64, branchID int64) ([]InventoryMovement, error) {
	return s.repo.GetProductKardex(ctx, productID, branchID)
}

func (s *service) ListReasons(ctx context.Context) ([]MovementReason, error) {
	return s.repo.ListReasons(ctx)
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
