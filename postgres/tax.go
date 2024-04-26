package postgres

import (
	"math"

	"github.com/hanqqv/assessment-tax/tax"
)

func (p *Postgres) CalculateTax(userInfo tax.UserInfo) (tax.Tax, error) {
	personalDeduction := 60000.0
	var taxAmount float64
	netAmount := userInfo.TotalIncome - personalDeduction

	if netAmount <= 150000.0 {
		taxAmount = 0.0
		taxAmount = math.Round(taxAmount*100) / 100
		return tax.Tax{Tax: taxAmount}, nil
	}

	if netAmount <= 500000.0 {
		taxAmount = 0.10 * (netAmount - 150000.0)
		taxAmount = math.Round(taxAmount*100) / 100
		return tax.Tax{Tax: taxAmount}, nil
	}

	if netAmount <= 1000000.0 {
		taxAmount = 35000.0 + (0.15 * (netAmount - 500000.0))
		taxAmount = math.Round(taxAmount*100) / 100
		return tax.Tax{Tax: taxAmount}, nil
	}

	if netAmount <= 2000000.0 {
		taxAmount = 35000.0 + 75000.0 + (0.20 * (netAmount - 1000000.0))
		taxAmount = math.Round(taxAmount*100) / 100
		return tax.Tax{Tax: taxAmount}, nil
	}
	taxAmount = 35000.0 + 75000.0 + 200000.0 + (0.35 * (netAmount - 2000000.0))
	taxAmount = math.Round(taxAmount*100) / 100
	return tax.Tax{Tax: taxAmount}, nil
}
