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

	if netAmount > 150000.0 {
		taxAmount += 0.0
		defaultTaxLevels[0].Tax = 0.0
		netAmount -= 150000.0
	} else {
		taxAmount += 0.0
		defaultTaxLevels[0].Tax = 0.0
		netAmount = 0.0

	}

	if netAmount > 350000.0 {
		taxAmount += 0.10 * 350000.0
		defaultTaxLevels[1].Tax = 0.10 * 350000.0
		netAmount -= 350000.0
	} else {
		taxAmount += 0.10 * netAmount
		defaultTaxLevels[1].Tax = 0.10 * netAmount
		netAmount = 0.0
	}

	if netAmount > 500000.0 {
		taxAmount += 0.15 * 500000.0
		defaultTaxLevels[2].Tax = 0.15 * 500000.0
		netAmount -= 500000.0
	} else {
		taxAmount += 0.15 * netAmount
		defaultTaxLevels[2].Tax = 0.15 * netAmount
		netAmount = 0.0
	}

	if netAmount > 1000000.0 {
		taxAmount += 0.20 * 1000000.0
		defaultTaxLevels[3].Tax = 0.20 * 1000000.0
		netAmount -= 1000000.0
	} else {
		taxAmount += 0.20 * netAmount
		defaultTaxLevels[3].Tax = 0.20 * netAmount
		netAmount = 0.0

	}

	if netAmount > 0.0 {
		taxAmount += 0.35 * netAmount
		defaultTaxLevels[4].Tax = 0.35 * netAmount
		netAmount = 0.0
	}

	for i := range defaultTaxLevels {
		defaultTaxLevels[i].Tax = math.Round(defaultTaxLevels[i].Tax*100) / 100
	}

	taxAmount -= userInfo.WHT
	taxAmount = math.Round(taxAmount*100) / 100

	return tax.Tax{Tax: taxAmount, TaxLevel: append([]tax.TaxLevel(nil), defaultTaxLevels...)}, nil
}
