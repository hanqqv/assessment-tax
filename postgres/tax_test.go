// go:build unit

package postgres

import (
	"testing"

	"github.com/hanqqv/assessment-tax/tax"
	"github.com/stretchr/testify/assert"
)

type StubTax struct {
	calculateTax tax.Tax
	err          error
}

func (s *StubTax) CalculateTax(userInfo tax.UserInfo) (tax.Tax, error) {
	return s.calculateTax, s.err
}

func TestCalculateTax(t *testing.T) {
	test := []struct {
		name     string
		userInfo tax.UserInfo
		wantTax  tax.Tax
	}{
		{
			name:     "Tax = 0.0 when TotalIncome = 0.0",
			userInfo: tax.UserInfo{TotalIncome: 0.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 1000.0",
			userInfo: tax.UserInfo{TotalIncome: 10000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 100000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 100000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 100000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 100000.0 & WHT = 5000.0",
			userInfo: tax.UserInfo{TotalIncome: 100000.0, WHT: 5000, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 100000.0 & WHT = 5000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 100000.0, WHT: 5000, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 140000.0",
			userInfo: tax.UserInfo{TotalIncome: 140000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 150000.0",
			userInfo: tax.UserInfo{TotalIncome: 150000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 150000.0 & WHT = 5000.0",
			userInfo: tax.UserInfo{TotalIncome: 150000.0, WHT: 5000, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 150000.0 & WHT = 5000.0 &Donation allowance = 50000.0",
			userInfo: tax.UserInfo{TotalIncome: 150000.0, WHT: 5000, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 50000.0}}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 150000.0 & WHT = 5000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 150000.0, WHT: 5000, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 150000.0 & WHT = 5000.0 & Donation allowance = 200000.0",
			userInfo: tax.UserInfo{TotalIncome: 150000.0, WHT: 5000, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 200000.0}}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 160000.0",
			userInfo: tax.UserInfo{TotalIncome: 160000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 1000.0 when TotalIncome = 220000.0",
			userInfo: tax.UserInfo{TotalIncome: 220000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 1000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 1000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 4000.0 when TotalIncome = 220000.0 & WHT = 5000.0",
			userInfo: tax.UserInfo{TotalIncome: 220000.0, WHT: 5000.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -4000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: -4000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 220000.0 & WHT = 5000.0 & Donation allowance = 50000.0",
			userInfo: tax.UserInfo{TotalIncome: 220000.0, WHT: 5000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 50000.0}}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 220000.0 & WHT = 5000.0 & Donation allowance = 200000.0",
			userInfo: tax.UserInfo{TotalIncome: 220000.0, WHT: 5000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 200000.0}}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: -5000.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 0.0 when TotalIncome = 220000.0 & Donation allowance = 80000.0",
			userInfo: tax.UserInfo{TotalIncome: 220000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 80000.0}}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 28000.0 when TotalIncome = 490000.0",
			userInfo: tax.UserInfo{TotalIncome: 490000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 28000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 28000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 29000.0 when TotalIncome = 500000.0",
			userInfo: tax.UserInfo{TotalIncome: 500000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 29000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 29000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 3000.0 when TotalIncome = 500000.0 & WHT = 32000.0",
			userInfo: tax.UserInfo{TotalIncome: 500000.0, WHT: 32000.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -3000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: -3000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 8000.0 when TotalIncome = 500000.0 & WHT = 32000.0 & Donation allowance = 50000.0",
			userInfo: tax.UserInfo{TotalIncome: 500000.0, WHT: 32000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 50000.0}}},
			wantTax: tax.Tax{Tax: -8000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: -8000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 19000.0 when TotalIncome = 500000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 500000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: 19000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 19000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 19044.45 when TotalIncome = 500444.5 & Donation allowance = 5500000.0",
			userInfo: tax.UserInfo{TotalIncome: 500444.5, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 5500000.0}}},
			wantTax: tax.Tax{Tax: 19044.45, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 19044.45},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 30000.0 when TotalIncome = 510000.0",
			userInfo: tax.UserInfo{TotalIncome: 510000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 30000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 30000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 20000.0 when TotalIncome = 510000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 510000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: 20000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 20000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 35000.0 when TotalIncome = 560000.0",
			userInfo: tax.UserInfo{TotalIncome: 560000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 35000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 35000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 5000.0 when TotalIncome = 560000.0 & WHT = 40000.0",
			userInfo: tax.UserInfo{TotalIncome: 560000.0, WHT: 40000.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -5000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: -5000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 12000.0 when TotalIncome = 560000.0 & WHT = 40000.0 & Donation allowance = 70000.0",
			userInfo: tax.UserInfo{TotalIncome: 560000.0, WHT: 40000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 70000.0}}},
			wantTax: tax.Tax{Tax: -12000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: -12000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 30000.0 when TotalIncome = 560000.0 & Donation allowance = 50000.0",
			userInfo: tax.UserInfo{TotalIncome: 560000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 50000.0}}},
			wantTax: tax.Tax{Tax: 30000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 30000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 35000.0 when TotalIncome = 660000.0 & Donation allowance = 400000.0",
			userInfo: tax.UserInfo{TotalIncome: 660000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 400000.0}}},
			wantTax: tax.Tax{Tax: 35000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 35000.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 99500.0 when TotalIncome = 990000.0",
			userInfo: tax.UserInfo{TotalIncome: 990000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 99500.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 99500.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 101000.0 when TotalIncome = 1000000.0",
			userInfo: tax.UserInfo{TotalIncome: 1000000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 101000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 101000.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 1000.0 when TotalIncome = 1000000.0 & WHT = 102000.0",
			userInfo: tax.UserInfo{TotalIncome: 1000000.0, WHT: 102000.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -1000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: -1000.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 14500.07 when TotalIncome = 1000000.0 & WHT = 102000.0 & Donation allowance = 90000.50",
			userInfo: tax.UserInfo{TotalIncome: 1000000.0, WHT: 102000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 90000.50}}},
			wantTax: tax.Tax{Tax: -14500.07, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: -14500.07},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 110000.0 when TotalIncome = 1060000.0",
			userInfo: tax.UserInfo{TotalIncome: 1060000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 110000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 110000.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 1000.0 when TotalIncome = 1060000.0 & WHT = 111000.0",
			userInfo: tax.UserInfo{TotalIncome: 1060000.0, WHT: 111000, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -1000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: -1000.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 52500.0 when TotalIncome = 1060000.0 & WHT = 50000.0 & Donation Allowance = 50000.0",
			userInfo: tax.UserInfo{TotalIncome: 1060000.0, WHT: 50000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 50000.0}}},
			wantTax: tax.Tax{Tax: 52500.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 52500.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 118000.0 when TotalIncome = 1100000.0",
			userInfo: tax.UserInfo{TotalIncome: 1100000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 118000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 118000.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 110000.0 when TotalIncome = 1160000.0 & Donation allowance = 300000.0",
			userInfo: tax.UserInfo{TotalIncome: 1160000.0, WHT: 0.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 300000.0}}},
			wantTax: tax.Tax{Tax: 110000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 110000.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 278000.0 when TotalIncome = 1900000.0",
			userInfo: tax.UserInfo{TotalIncome: 1900000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 278000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 278000.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 298000.0 when TotalIncome = 2000000.0",
			userInfo: tax.UserInfo{TotalIncome: 2000000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 298000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 298000.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 1999.9 when TotalIncome = 2000000.50 & WHT = 300000.0",
			userInfo: tax.UserInfo{TotalIncome: 2000000.50, WHT: 300000.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -1999.9, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: -1999.9},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 177999.5 when TotalIncome = 2000000.0 & WHT = 100000.5 & Donation Allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 2000000.0, WHT: 100000.5, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: 177999.5, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 177999.5},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 310000.0 when TotalIncome = 2060000.0",
			userInfo: tax.UserInfo{TotalIncome: 2060000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 310000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 310000.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 10000.5 when TotalIncome = 2060000.0 & WHT = 320000.50",
			userInfo: tax.UserInfo{TotalIncome: 2060000.0, WHT: 320000.50, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -10000.5, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: -10000.5},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 9999.89 when TotalIncome = 2160001.75 & WHT = 320000.50 & Donation allowance = 700000.0",
			userInfo: tax.UserInfo{TotalIncome: 2160001.75, WHT: 320000.50, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 700000.0}}},
			wantTax: tax.Tax{Tax: -9999.89, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: -9999.89},
			}},
		},
		{
			name:     "TaxRefund = 10000.33 when TotalIncome = 2060000.5 & WHT = 320000.50",
			userInfo: tax.UserInfo{TotalIncome: 2060000.5, WHT: 320000.50, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -10000.33, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: -10000.33},
			}},
		},
		{
			name:     "TaxRefund = 9999.95 when TotalIncome = 2060000.0 & WHT = 300000.0 & Donation allowance = 99999.75",
			userInfo: tax.UserInfo{TotalIncome: 2060000.0, WHT: 300000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 99999.75}}},
			wantTax: tax.Tax{Tax: -9999.95, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: -9999.95},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "TaxRefund = 10000.0 when TotalIncome = 2060000.0 & WHT = 300000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 2060000.0, WHT: 300000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: -10000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: -10000.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		},
		{
			name:     "Tax = 324000.0 when TotalIncome = 2100000.0",
			userInfo: tax.UserInfo{TotalIncome: 2100000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 324000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 324000.0},
			}},
		},
		{
			name:     "Tax = 1339000.0 when TotalIncome = 5000000.0",
			userInfo: tax.UserInfo{TotalIncome: 5000000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 1339000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 1339000.0},
			}},
		},
		{
			name:     "TaxRefund = 1000.0 when TotalIncome = 5000000.0 & WHT = 1340000.0",
			userInfo: tax.UserInfo{TotalIncome: 5000000.0, WHT: 1340000.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: -1000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: -1000.0},
			}},
		},
		{
			name:     "Tax = 204000.0 when TotalIncome = 5000000.0 & WHT = 1100000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 5000000.0, WHT: 1100000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: 204000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 204000.0},
			}},
		},
		{
			name:     "Tax = 1954000.0 when TotalIncome = 10000000.0 & WHT = 1100000.0 & Donation allowance = 100000.0",
			userInfo: tax.UserInfo{TotalIncome: 10000000.0, WHT: 1100000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 100000.0}}},
			wantTax: tax.Tax{Tax: 1954000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 1954000.0},
			}},
		},
		{
			name:     "Tax = 1954000.0 when TotalIncome = 10000000.0 & WHT = 1100000.0 & Donation allowance = 1000000.0",
			userInfo: tax.UserInfo{TotalIncome: 10000000.0, WHT: 1100000.0, Allowances: []tax.Allowances{{AllowanceType: "donation", Amount: 1000000.0}}},
			wantTax: tax.Tax{Tax: 1954000.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 1954000.0},
			}},
		},
		{
			name:     "Negative TotalIncome",
			userInfo: tax.UserInfo{TotalIncome: -1000.0, WHT: 0.0, Allowances: []tax.Allowances{}},
			wantTax: tax.Tax{Tax: 0.0, TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 0.0},
				{Level: "500,001-1,000,000", Tax: 0.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			}},
		}}

	p := &Postgres{}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.CalculateTax(tt.userInfo)
			assert.NoError(t, err, "expected no error but got %v", err)
			assert.Equal(t, tt.wantTax, got, "expected tax %v but got %v", tt.wantTax, got)
		})
	}
}
