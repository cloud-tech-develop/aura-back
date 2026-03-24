package transfers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTransferNotFound = errors.New("traslado no encontrado")
	ErrInvalidItems     = errors.New("items inválidos")
	ErrSameBranch       = errors.New("las sucursales de origen y destino deben ser diferentes")
	ErrAlreadyApproved  = errors.New("traslado ya aprobado")
	ErrAlreadyShipped   = errors.New("traslado ya enviado")
	ErrAlreadyReceived  = errors.New("traslado ya recibido")
	ErrAlreadyCancelled = errors.New("traslado ya cancelado")
	ErrNotApproved      = errors.New("traslado debe estar aprobado para enviar")
	ErrNotShipped       = errors.New("traslado debe estar enviado para recibir")
	ErrCannotCancel     = errors.New("no se puede cancelar un traslado enviado o recibido")
)

type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) generateTransferNumber() string {
	return fmt.Sprintf("TRF-%d", time.Now().UnixNano())
}

// ─── HU-TRANS-001: Create Transfer Request ────────────────────────────────────

func (s *service) CreateTransfer(ctx context.Context, userID int64, req CreateTransferRequest) (*Transfer, error) {
	if len(req.Items) == 0 {
		return nil, ErrInvalidItems
	}

	if req.OriginBranchID == req.DestinationBranchID {
		return nil, ErrSameBranch
	}

	transfer := &Transfer{
		TransferNumber:      s.generateTransferNumber(),
		OriginBranchID:      req.OriginBranchID,
		DestinationBranchID: req.DestinationBranchID,
		UserID:              userID,
		Notes:               req.Notes,
	}

	transferID, err := s.repo.CreateTransfer(ctx, transfer)
	if err != nil {
		return nil, fmt.Errorf("creando traslado: %w", err)
	}
	transfer.ID = transferID

	for _, itemReq := range req.Items {
		item := &TransferItem{
			TransferID:        transferID,
			ProductID:         itemReq.ProductID,
			RequestedQuantity: itemReq.RequestedQuantity,
		}
		if err := s.repo.CreateTransferItem(ctx, item); err != nil {
			return nil, fmt.Errorf("creando item de traslado: %w", err)
		}
	}

	return transfer, nil
}

// ─── HU-TRANS-002: Approve Transfer ──────────────────────────────────────────

func (s *service) ApproveTransfer(ctx context.Context, transferID int64, userID int64) (*Transfer, error) {
	transfer, err := s.repo.GetTransferByID(ctx, transferID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransferNotFound
		}
		return nil, fmt.Errorf("obteniendo traslado: %w", err)
	}

	if transfer.Status != StatusPending {
		return nil, fmt.Errorf("traslado no está pendiente")
	}

	transfer.Status = StatusApproved
	if err := s.repo.UpdateTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("aprobando traslado: %w", err)
	}

	return transfer, nil
}

// ─── HU-TRANS-003: Ship Transfer ──────────────────────────────────────────────

func (s *service) ShipTransfer(ctx context.Context, transferID int64, userID int64, req ShipTransferRequest) (*Transfer, error) {
	transfer, err := s.repo.GetTransferByID(ctx, transferID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransferNotFound
		}
		return nil, fmt.Errorf("obteniendo traslado: %w", err)
	}

	if transfer.Status != StatusApproved && transfer.Status != StatusPending {
		return nil, ErrNotApproved
	}

	if len(req.Items) == 0 {
		return nil, ErrInvalidItems
	}

	// Update items with shipped quantities
	existingItems, err := s.repo.GetTransferItems(ctx, transferID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo items: %w", err)
	}

	// Create a map for quick lookup
	itemsMap := make(map[int64]*TransferItem)
	for i := range existingItems {
		itemsMap[existingItems[i].ProductID] = &existingItems[i]
	}

	for _, shipItem := range req.Items {
		if existingItem, exists := itemsMap[shipItem.ProductID]; exists {
			existingItem.ShippedQuantity = shipItem.ShippedQuantity
			if err := s.repo.UpdateTransferItem(ctx, existingItem); err != nil {
				return nil, fmt.Errorf("actualizando item: %w", err)
			}
		}
	}

	now := time.Now()
	transfer.Status = StatusShipped
	transfer.ShippedDate = &now
	transfer.ShippedBy = &userID
	transfer.Notes = req.Notes

	if err := s.repo.UpdateTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("enviando traslado: %w", err)
	}

	return transfer, nil
}

// ─── HU-TRANS-004: Receive Transfer ──────────────────────────────────────────

func (s *service) ReceiveTransfer(ctx context.Context, transferID int64, userID int64, req ReceiveTransferRequest) (*Transfer, error) {
	transfer, err := s.repo.GetTransferByID(ctx, transferID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransferNotFound
		}
		return nil, fmt.Errorf("obteniendo traslado: %w", err)
	}

	if transfer.Status != StatusShipped {
		return nil, ErrNotShipped
	}

	if len(req.Items) == 0 {
		return nil, ErrInvalidItems
	}

	// Update items with received quantities
	existingItems, err := s.repo.GetTransferItems(ctx, transferID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo items: %w", err)
	}

	itemsMap := make(map[int64]*TransferItem)
	for i := range existingItems {
		itemsMap[existingItems[i].ProductID] = &existingItems[i]
	}

	allFullyReceived := true
	for _, recvItem := range req.Items {
		if existingItem, exists := itemsMap[recvItem.ProductID]; exists {
			existingItem.ReceivedQuantity = recvItem.ReceivedQuantity
			if err := s.repo.UpdateTransferItem(ctx, existingItem); err != nil {
				return nil, fmt.Errorf("actualizando item: %w", err)
			}
			if recvItem.ReceivedQuantity < existingItem.ShippedQuantity {
				allFullyReceived = false
			}
		}
	}

	now := time.Now()
	transfer.ReceivedDate = &now
	transfer.ReceivedBy = &userID
	transfer.Notes = req.Notes

	if allFullyReceived {
		transfer.Status = StatusReceived
	} else {
		transfer.Status = StatusPartial
	}

	if err := s.repo.UpdateTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("recibiendo traslado: %w", err)
	}

	return transfer, nil
}

// ─── HU-TRANS-005: Cancel Transfer ──────────────────────────────────────────

func (s *service) CancelTransfer(ctx context.Context, transferID int64, userID int64, reason string) (*Transfer, error) {
	transfer, err := s.repo.GetTransferByID(ctx, transferID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransferNotFound
		}
		return nil, fmt.Errorf("obteniendo traslado: %w", err)
	}

	if transfer.Status == StatusCancelled {
		return nil, ErrAlreadyCancelled
	}

	if transfer.Status == StatusReceived {
		return nil, ErrCannotCancel
	}

	now := time.Now()
	transfer.Status = StatusCancelled
	transfer.CancelledBy = &userID
	transfer.CancelledAt = &now
	transfer.CancellationReason = reason

	if err := s.repo.UpdateTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("cancelando traslado: %w", err)
	}

	return transfer, nil
}

// ─── HU-TRANS-006: View Transfer History ──────────────────────────────────────

func (s *service) ListTransfers(ctx context.Context, originBranchID, destBranchID *int64, status string, startDate, endDate *time.Time, page, limit int) ([]Transfer, int64, error) {
	return s.repo.ListTransfers(ctx, originBranchID, destBranchID, status, startDate, endDate, page, limit)
}

// ─── Additional ──────────────────────────────────────────────────────────────

func (s *service) GetTransferByID(ctx context.Context, id int64) (*Transfer, error) {
	return s.repo.GetTransferByID(ctx, id)
}
