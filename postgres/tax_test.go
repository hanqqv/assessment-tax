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
		wantTax  float64
	}{
		{"Tax = 0.0 when TotalIncome = 0.0", tax.UserInfo{TotalIncome: 0.0, Allowances: []tax.Allowances{}}, 0.0},
		{"Tax = 0.0 when TotalIncome = 1000.0", tax.UserInfo{TotalIncome: 10000.0, Allowances: []tax.Allowances{}}, 0.0},
		{"Tax = 0.0 when TotalIncome = 100000.0", tax.UserInfo{TotalIncome: 100000.0, Allowances: []tax.Allowances{}}, 0.0},
		{"Tax = 0.0 when TotalIncome = 140000.0", tax.UserInfo{TotalIncome: 140000.0, Allowances: []tax.Allowances{}}, 0.0},
		{"Tax = 0.0 when TotalIncome = 150000.0", tax.UserInfo{TotalIncome: 150000.0, Allowances: []tax.Allowances{}}, 0.0},
		{"Tax = 0.0 when TotalIncome = 160000.0", tax.UserInfo{TotalIncome: 160000.0, Allowances: []tax.Allowances{}}, 0.0},
		{"Tax = 1000.0 when TotalIncome = 220000.0", tax.UserInfo{TotalIncome: 220000.0, Allowances: []tax.Allowances{}}, 1000.0},
		{"Tax = 28000.0 when TotalIncome = 490000.0", tax.UserInfo{TotalIncome: 490000.0, Allowances: []tax.Allowances{}}, 28000.0},
		{"Tax = 29000.0 when TotalIncome = 500000.0", tax.UserInfo{TotalIncome: 500000.0, Allowances: []tax.Allowances{}}, 29000.0},
		{"Tax = 30000.0 when TotalIncome = 510000.0", tax.UserInfo{TotalIncome: 510000.0, Allowances: []tax.Allowances{}}, 30000.0},
		{"Tax = 35000.0 when TotalIncome = 560000.0", tax.UserInfo{TotalIncome: 560000.0, Allowances: []tax.Allowances{}}, 35000.0},
		{"Tax = 99500.0 when TotalIncome = 990000.0", tax.UserInfo{TotalIncome: 990000.0, Allowances: []tax.Allowances{}}, 99500.0},
		{"Tax = 101000.0 when TotalIncome = 1000000.0", tax.UserInfo{TotalIncome: 1000000.0, Allowances: []tax.Allowances{}}, 101000.0},
		{"Tax = 110000.0 when TotalIncome = 1060000.0", tax.UserInfo{TotalIncome: 1060000.0, Allowances: []tax.Allowances{}}, 110000.0},
		{"Tax = 118000.0 when TotalIncome = 1100000.0", tax.UserInfo{TotalIncome: 1100000.0, Allowances: []tax.Allowances{}}, 118000.0},
		{"Tax = 278000.0 when TotalIncome = 1900000.0", tax.UserInfo{TotalIncome: 1900000.0, Allowances: []tax.Allowances{}}, 278000.0},
		{"Tax = 298000.0 when TotalIncome = 2000000.0", tax.UserInfo{TotalIncome: 2000000.0, Allowances: []tax.Allowances{}}, 298000.0},
		{"Tax = 310000.0 when TotalIncome = 2060000.0", tax.UserInfo{TotalIncome: 2060000.0, Allowances: []tax.Allowances{}}, 310000.0},
		{"Tax = 324000.0 when TotalIncome = 2100000.0", tax.UserInfo{TotalIncome: 2100000.0, Allowances: []tax.Allowances{}}, 324000.0},
		{"Tax = 1339000.0 when TotalIncome = 5000000.0", tax.UserInfo{TotalIncome: 5000000.0, Allowances: []tax.Allowances{}}, 1339000.0},
		{"Negative TotalIncome", tax.UserInfo{TotalIncome: -1000.0, Allowances: []tax.Allowances{}}, 0.0},
	}

	p := &Postgres{}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.CalculateTax(tt.userInfo)
			assert.NoError(t, err, "expected no error but got %v", err)
			assert.Equal(t, tt.wantTax, got.Tax, "expected tax %v but got %v", tt.wantTax, got.Tax)
		})
	}

}
