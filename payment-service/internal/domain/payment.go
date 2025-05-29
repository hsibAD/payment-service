package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidPaymentID     = errors.New("invalid payment ID")
	ErrInvalidOrderID       = errors.New("invalid order ID")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrInvalidCurrency      = errors.New("invalid currency")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusCancelled  PaymentStatus = "CANCELLED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodMetaMask   PaymentMethod = "METAMASK"
)

type Payment struct {
	ID            string
	OrderID       string
	UserID        string
	Amount        float64
	Currency      string
	Status        PaymentStatus
	PaymentMethod PaymentMethod
	TransactionID string
	ErrorMessage  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreditCardInfo struct {
	CardNumber     string
	ExpiryMonth    string
	ExpiryYear     string
	CVV            string
	CardholderName string
}

type MetaMaskInfo struct {
	WalletAddress  string
	TransactionHash string
	ContractAddress string
	AmountWei      string
}

func NewPayment(
	orderID string,
	userID string,
	amount float64,
	currency string,
	method PaymentMethod,
) (*Payment, error) {
	if orderID == "" {
		return nil, ErrInvalidOrderID
	}

	if userID == "" {
		return nil, ErrInvalidUserID
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if currency == "" {
		return nil, ErrInvalidCurrency
	}

	if method != PaymentMethodCreditCard && method != PaymentMethodMetaMask {
		return nil, ErrInvalidPaymentMethod
	}

	return &Payment{
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		Status:        PaymentStatusPending,
		PaymentMethod: method,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (p *Payment) UpdateStatus(status PaymentStatus) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

func (p *Payment) SetTransactionID(txID string) {
	p.TransactionID = txID
	p.UpdatedAt = time.Now()
}

func (p *Payment) SetError(err string) {
	p.ErrorMessage = err
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
}

func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusProcessing
}

func (p *Payment) CanBeRetried() bool {
	return p.Status == PaymentStatusFailed || p.Status == PaymentStatusCancelled
}

func (p *Payment) MarkAsProcessing() {
	if p.Status == PaymentStatusPending {
		p.Status = PaymentStatusProcessing
		p.UpdatedAt = time.Now()
	}
}

func (p *Payment) MarkAsCompleted(transactionID string) {
	if p.IsPending() {
		p.Status = PaymentStatusCompleted
		p.TransactionID = transactionID
		p.UpdatedAt = time.Now()
	}
}

func (p *Payment) Cancel() {
	if p.IsPending() {
		p.Status = PaymentStatusCancelled
		p.UpdatedAt = time.Now()
	}
}

func (p *Payment) Refund() error {
	if p.Status != PaymentStatusCompleted {
		return errors.New("only completed payments can be refunded")
	}

	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
} 