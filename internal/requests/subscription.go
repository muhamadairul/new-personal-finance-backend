package requests

// QrisPaymentRequest holds QRIS checkout payload
type QrisPaymentRequest struct {
	Plan string `json:"plan" validate:"required,oneof=monthly yearly"`
}

// VaPaymentRequest holds Virtual Account checkout payload
type VaPaymentRequest struct {
	Plan     string `json:"plan" validate:"required,oneof=monthly yearly"`
	BankCode string `json:"bank_code" validate:"required,oneof=BCA BNI BRI MANDIRI PERMATA"`
}

// EwalletPaymentRequest holds E-Wallet checkout payload
type EwalletPaymentRequest struct {
	Plan        string `json:"plan" validate:"required,oneof=monthly yearly"`
	ChannelCode string `json:"channel_code" validate:"required,oneof=OVO DANA SHOPEEPAY LINKAJA"`
}
