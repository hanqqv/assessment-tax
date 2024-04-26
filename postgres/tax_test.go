// go:build unit

package postgres

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
	t.Run("CalculateTax Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err, "an error was not expected when opening a stub database connection")
		defer db.Close()

		p := &Postgres{DB: db}

		mock.ExpectQuery("SELECT amount FROM deductions_setting WHERE allowance_type = \\$1").
			WithArgs("personal").
			WillReturnError(errors.New("mock error"))

		userInfo := tax.UserInfo{
			TotalIncome: 600000.0,
			Allowances:  []tax.Allowances{},
			WHT:         0.0,
		}

		_, err = p.CalculateTax(userInfo)
		assert.Error(t, err, "CalculateTax should return an error")
	})
	t.Run("CalculateTax Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err, "an error was not expected when opening a stub database connection")
		defer db.Close()

		p := &Postgres{DB: db}

		wantTax := tax.Tax{
			Tax: 42500.0,
			TaxLevel: []tax.TaxLevel{
				{Level: "0-150,000", Tax: 0.0},
				{Level: "150,001-500,000", Tax: 35000.0},
				{Level: "500,001-1,000,000", Tax: 7500.0},
				{Level: "1,000,001-2,000,000", Tax: 0.0},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
			},
		}

		mock.ExpectQuery("SELECT amount FROM deductions_setting WHERE allowance_type = \\$1").
			WithArgs("personal").
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(50000.0))

		userInfo := tax.UserInfo{
			TotalIncome: 600000.0,
			Allowances:  []tax.Allowances{},
			WHT:         0.0,
		}

		gotTax, err := p.CalculateTax(userInfo)
		assert.NoError(t, err, "CalculateTax returned an error: %v", err)
		assert.Equal(t, wantTax, gotTax, "CalculateTax returned incorrect tax: got %v want %v", gotTax, wantTax)

		for i, v := range gotTax.TaxLevel {
			if v != wantTax.TaxLevel[i] {
				t.Errorf("CalculateTax returned incorrect tax level: got %v want %v",
					v, wantTax.TaxLevel[i])
			}
		}
	})
}
func TestGetPersonalDeduction(t *testing.T) {
	t.Run("GetPersonalDeduction Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err, "an error was not expected when opening a stub database connection")
		defer db.Close()

		p := &Postgres{DB: db}

		wantDeduction := 40000.0
		mock.ExpectQuery("SELECT amount FROM deductions_setting WHERE allowance_type = \\$1").
			WithArgs("personal").
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(wantDeduction))

		gotDeduction, err := p.getPersonalDeduction()

		assert.NoError(t, err, "GetPersonalDeduction returned an error: %v", err)
		assert.Equal(t, wantDeduction, gotDeduction, "GetPersonalDeduction returned incorrect deduction: got %v want %v", gotDeduction, wantDeduction)
	})
	t.Run("GetPersonalDeduction Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err, "an error was not expected when opening a stub database connection")
		defer db.Close()

		p := &Postgres{DB: db}

		mock.ExpectQuery("SELECT amount FROM deductions_setting WHERE allowance_type = \\$1").
			WithArgs("personal").
			WillReturnError(err)

		_, gotErr := p.getPersonalDeduction()
		assert.Error(t, gotErr, "GetPersonalDeduction did not return an error")
	})
}
func TestSettingPersonalDeduction(t *testing.T) {
	t.Run("Given deduction amount to update SettingPersonalDeduction Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err, "an error was not expected when opening a stub database connection")
		defer db.Close()

		p := &Postgres{DB: db}

		expectedDeduction := 100000.0

		mock.ExpectQuery("UPDATE deductions_setting SET amount = \\$1 WHERE allowance_type = \\$2 RETURNING amount").
			WithArgs(100000.0, "personal").
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(expectedDeduction))

		setting := tax.Setting{Amount: 100000.0}
		deduction, err := p.SettingPersonalDeduction(setting)

		assert.NoError(t, err, "SettingPersonalDeduction returned an error: %v", err)
		assert.Equal(t, expectedDeduction, deduction, "SettingPersonalDeduction returned incorrect deduction: got %v want %v", deduction, expectedDeduction)
	})
	t.Run("Given deduction amount to update SettingPersonalDeduction Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err, "an error was not expected when opening a stub database connection")
		defer db.Close()

		p := &Postgres{DB: db}

		mock.ExpectQuery("UPDATE deductions_setting SET amount = \\$1 WHERE allowance_type = \\$2 RETURNING amount").
			WithArgs(5000.0, "personal").
			WillReturnError(err)

		setting := tax.Setting{Amount: 5000.0}
		_, gotErr := p.SettingPersonalDeduction(setting)
		assert.Error(t, gotErr, "SettingPersonalDeduction did not return an error")
	})
}
