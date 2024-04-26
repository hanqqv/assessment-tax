package postgres

import (
	"math"

	"github.com/hanqqv/assessment-tax/tax"
)

func (p *Postgres) CalculateTax(userInfo tax.UserInfo) (tax.Tax, error) {
	personalDeduction := 60000.0
	var taxAmount float64
	netAmount := userInfo.TotalIncome - personalDeduction
	defaultTaxLevels := []tax.TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	for _, allowance := range userInfo.Allowances {
		if allowance.AllowanceType == "donation" && allowance.Amount > 100000.0 {
			allowance.Amount = 100000.0
		}
		netAmount -= allowance.Amount
	}

	if netAmount <= 150000.0 {
		taxAmount = 0.0 - userInfo.WHT
		taxAmount = math.Round(taxAmount*100) / 100
		defaultTaxLevels[0].Tax = taxAmount
		return tax.Tax{Tax: taxAmount, TaxLevel: append([]tax.TaxLevel(nil), defaultTaxLevels...)}, nil
	}

	if netAmount <= 500000.0 {
		taxAmount = 0.10*(netAmount-150000.0) - userInfo.WHT
		taxAmount = math.Round(taxAmount*100) / 100
		defaultTaxLevels[1].Tax = taxAmount
		return tax.Tax{Tax: taxAmount, TaxLevel: append([]tax.TaxLevel(nil), defaultTaxLevels...)}, nil
	}

	if netAmount <= 1000000.0 {
		taxAmount = 35000.0 + (0.15 * (netAmount - 500000.0)) - userInfo.WHT
		taxAmount = math.Round(taxAmount*100) / 100
		defaultTaxLevels[2].Tax = taxAmount
		return tax.Tax{Tax: taxAmount, TaxLevel: append([]tax.TaxLevel(nil), defaultTaxLevels...)}, nil
	}

	if netAmount <= 2000000.0 {
		taxAmount = 35000.0 + 75000.0 + (0.20 * (netAmount - 1000000.0)) - userInfo.WHT
		taxAmount = math.Round(taxAmount*100) / 100
		defaultTaxLevels[3].Tax = taxAmount
		return tax.Tax{Tax: taxAmount, TaxLevel: append([]tax.TaxLevel(nil), defaultTaxLevels...)}, nil
	}
	taxAmount = 35000.0 + 75000.0 + 200000.0 + (0.35 * (netAmount - 2000000.0)) - userInfo.WHT
	taxAmount = math.Round(taxAmount*100) / 100
	defaultTaxLevels[4].Tax = taxAmount
	return tax.Tax{Tax: taxAmount, TaxLevel: append([]tax.TaxLevel(nil), defaultTaxLevels...)}, nil
}
