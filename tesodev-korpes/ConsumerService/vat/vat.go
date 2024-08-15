package vat

type VatCalculator struct {
	VatRate float64
}

// NewVatCalculator KDV hesaplama için yeni bir VatCalculator döner
func NewVatCalculator(vatRate float64) *VatCalculator {
	return &VatCalculator{VatRate: vatRate}
}

// SetVatRate KDV oranını günceller
func (vc *VatCalculator) SetVatRate(vatRate float64) {
	vc.VatRate = vatRate
}

// CalculateVat Fiyatın üzerine KDV ekleyerek toplam fiyatı döndürür
func (vc *VatCalculator) CalculateVat(price float64) float64 {
	vatAmount := price * (vc.VatRate / 100)
	return price + vatAmount
}
