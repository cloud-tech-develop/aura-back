package purchases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPurchaseOrderNotFound = errors.New("orden de compra no encontrada")
	ErrPurchaseNotFound      = errors.New("compra no encontrada")
	ErrInvalidItems          = errors.New("items inválidos")
	ErrOrderAlreadyCompleted = errors.New("la orden ya fue completada")
	ErrInvalidPaymentAmount  = errors.New("monto de pago inválido")
)

type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) generateOrderNumber(branchID int64) string {
	return fmt.Sprintf("PO-%d-%d", branchID, time.Now().UnixNano())
}

func (s *service) generatePurchaseNumber(branchID int64) string {
	return fmt.Sprintf("PUR-%d-%d", branchID, time.Now().UnixNano())
}

// ─── HU-PUR-001: Create Purchase Order ───────────────────────────────────────

func (s *service) CreatePurchaseOrder(ctx context.Context, userID int64, req CreatePurchaseOrderRequest) (*PurchaseOrder, error) {
	if len(req.Items) == 0 {
		return nil, ErrInvalidItems
	}

	po := &PurchaseOrder{
		OrderNumber:  s.generateOrderNumber(req.BranchID),
		SupplierID:   req.SupplierID,
		BranchID:     req.BranchID,
		UserID:       userID,
		OrderDate:    time.Now(),
		ExpectedDate: req.ExpectedDate,
		Notes:        req.Notes,
	}

	var subtotal, discountTotal, taxTotal, total float64

	poID, err := s.repo.CreatePurchaseOrder(ctx, po)
	if err != nil {
		return nil, fmt.Errorf("creando orden: %w", err)
	}
	po.ID = poID

	for _, itemReq := range req.Items {
		lineTotal := itemReq.Quantity*itemReq.UnitCost - itemReq.DiscountAmount
		taxAmount := lineTotal * itemReq.TaxRate / 100

		item := &PurchaseOrderItem{
			PurchaseOrderID: poID,
			ProductID:       itemReq.ProductID,
			Quantity:        itemReq.Quantity,
			UnitCost:        itemReq.UnitCost,
			DiscountAmount:  itemReq.DiscountAmount,
			TaxRate:         itemReq.TaxRate,
			LineTotal:       lineTotal,
		}
		if err := s.repo.CreatePurchaseOrderItem(ctx, item); err != nil {
			return nil, fmt.Errorf("creando item: %w", err)
		}

		subtotal += lineTotal
		discountTotal += itemReq.DiscountAmount
		taxTotal += taxAmount
		total += lineTotal + taxAmount
	}

	po.Subtotal = subtotal
	po.DiscountTotal = discountTotal
	po.TaxTotal = taxTotal
	po.Total = total
	if err := s.repo.UpdatePurchaseOrder(ctx, po); err != nil {
		return nil, fmt.Errorf("actualizando totales: %w", err)
	}

	return po, nil
}

// ─── HU-PUR-002: Receive Goods ────────────────────────────────────────────────

func (s *service) ReceiveGoods(ctx context.Context, userID int64, req ReceiveGoodsRequest) (*Purchase, error) {
	// Verify order exists
	order, err := s.repo.GetPurchaseOrderByID(ctx, req.PurchaseOrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPurchaseOrderNotFound
		}
		return nil, fmt.Errorf("obteniendo orden: %w", err)
	}

	if order.Status == StatusCancelled || order.Status == StatusReceived {
		return nil, ErrOrderAlreadyCompleted
	}

	purchase := &Purchase{
		PurchaseNumber:  s.generatePurchaseNumber(order.BranchID),
		PurchaseOrderID: &order.ID,
		SupplierID:      order.SupplierID,
		BranchID:        order.BranchID,
		UserID:          userID,
		Notes:           req.Notes,
	}

	var subtotal, discountTotal, taxTotal, total float64

	purchaseID, err := s.repo.CreatePurchase(ctx, purchase)
	if err != nil {
		return nil, fmt.Errorf("creando compra: %w", err)
	}
	purchase.ID = purchaseID

	for _, itemReq := range req.Items {
		lineTotal := itemReq.Quantity*itemReq.UnitCost - itemReq.DiscountAmount
		taxAmount := lineTotal * itemReq.TaxRate / 100

		item := &PurchaseItem{
			PurchaseID:     purchaseID,
			ProductID:      itemReq.ProductID,
			Quantity:       itemReq.Quantity,
			UnitCost:       itemReq.UnitCost,
			DiscountAmount: itemReq.DiscountAmount,
			TaxRate:        itemReq.TaxRate,
			LineTotal:      lineTotal,
		}
		if err := s.repo.CreatePurchaseItem(ctx, item); err != nil {
			return nil, fmt.Errorf("creando item de compra: %w", err)
		}

		subtotal += lineTotal
		discountTotal += itemReq.DiscountAmount
		taxTotal += taxAmount
		total += lineTotal + taxAmount
	}

	purchase.Subtotal = subtotal
	purchase.DiscountTotal = discountTotal
	purchase.TaxTotal = taxTotal
	purchase.Total = total
	purchase.PendingAmount = total
	if err := s.repo.UpdatePurchase(ctx, purchase); err != nil {
		return nil, fmt.Errorf("actualizando totales de compra: %w", err)
	}

	// Update order status
	order.Status = StatusPartial
	if err := s.repo.UpdatePurchaseOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("actualizando estado de orden: %w", err)
	}

	return purchase, nil
}

// ─── HU-PUR-003: Record Purchase Payment ──────────────────────────────────────

func (s *service) RecordPayment(ctx context.Context, userID int64, req RecordPaymentRequest) (*PurchasePayment, error) {
	purchase, err := s.repo.GetPurchaseByID(ctx, req.PurchaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPurchaseNotFound
		}
		return nil, fmt.Errorf("obteniendo compra: %w", err)
	}

	if purchase.Status == StatusCancelled {
		return nil, fmt.Errorf("compra cancelada no acepta pagos")
	}

	if req.Amount > purchase.PendingAmount {
		return nil, ErrInvalidPaymentAmount
	}

	payment := &PurchasePayment{
		PurchaseID:      req.PurchaseID,
		PaymentMethod:   req.PaymentMethod,
		Amount:          req.Amount,
		ReferenceNumber: &req.ReferenceNumber,
		Notes:           req.Notes,
		UserID:          userID,
	}

	if err := s.repo.CreatePurchasePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("registrando pago: %w", err)
	}

	// Update purchase paid amount
	purchase.PaidAmount += req.Amount
	purchase.PendingAmount = purchase.Total - purchase.PaidAmount
	if purchase.PendingAmount <= 0 {
		purchase.Status = StatusCompleted
	} else {
		purchase.Status = StatusPartial
	}

	if err := s.repo.UpdatePurchase(ctx, purchase); err != nil {
		return nil, fmt.Errorf("actualizando compra: %w", err)
	}

	return payment, nil
}

// ─── HU-PUR-004: Cancel Purchase ──────────────────────────────────────────────

func (s *service) CancelPurchase(ctx context.Context, purchaseID int64, reason string) error {
	purchase, err := s.repo.GetPurchaseByID(ctx, purchaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPurchaseNotFound
		}
		return fmt.Errorf("obteniendo compra: %w", err)
	}

	if purchase.Status == StatusCancelled {
		return fmt.Errorf("compra ya cancelada")
	}

	if purchase.PaidAmount > 0 {
		return fmt.Errorf("no se puede cancelar una compra con pagos realizados")
	}

	purchase.Status = StatusCancelled
	purchase.Notes = reason
	if err := s.repo.UpdatePurchase(ctx, purchase); err != nil {
		return fmt.Errorf("cancelando compra: %w", err)
	}

	return nil
}

// ─── HU-PUR-005: View Purchase History ─────────────────────────────────────────

func (s *service) GetPurchaseHistory(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Purchase, int64, error) {
	return s.repo.ListPurchases(ctx, supplierID, status, startDate, endDate, page, limit)
}

// ─── HU-PUR-006: Supplier Account Summary ──────────────────────────────────────

func (s *service) GetSupplierSummary(ctx context.Context, supplierID int64) (*SupplierSummary, error) {
	return s.repo.GetSupplierSummary(ctx, supplierID)
}

// ─── Additional ──────────────────────────────────────────────────────────────

func (s *service) GetPurchaseOrderByID(ctx context.Context, id int64) (*PurchaseOrder, error) {
	return s.repo.GetPurchaseOrderByID(ctx, id)
}

func (s *service) GetPurchaseByID(ctx context.Context, id int64) (*Purchase, error) {
	return s.repo.GetPurchaseByID(ctx, id)
}

func (s *service) ListPurchaseOrders(ctx context.Context, supplierID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]PurchaseOrder, int64, error) {
	return s.repo.ListPurchaseOrders(ctx, supplierID, status, startDate, endDate, page, limit)
}
