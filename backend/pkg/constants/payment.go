package constants

// PaymentMethod represents a payment option
type PaymentMethod struct {
	Value       string
	Label       string
	Emoji       string
	SettingsKey string
}

// Available payment methods
var PaymentMethods = []PaymentMethod{
	{Value: "telebirr", Label: "TeleBirr", Emoji: "📱", SettingsKey: "merchant_telebirr"},
	{Value: "safaricom", Label: "M-PESA", Emoji: "💚", SettingsKey: "merchant_safaricom"},
	{Value: "cbe", Label: "CBE Birr", Emoji: "🏦", SettingsKey: "merchant_cbe"},
	{Value: "bank", Label: "Bank Transfer", Emoji: "🏛️", SettingsKey: "merchant_bank_account"},
}

// ValidPaymentMethods returns a map of valid payment method values
func ValidPaymentMethods() map[string]bool {
	valid := make(map[string]bool)
	for _, pm := range PaymentMethods {
		valid[pm.Value] = true
	}
	return valid
}

// IsValidPaymentMethod checks if a payment method is valid
func IsValidPaymentMethod(method string) bool {
	valid := ValidPaymentMethods()
	return valid[method]
}
