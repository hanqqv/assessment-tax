package postgres

import (
	"github.com/hanqqv/assessment-tax/tax"
)

func (p *Postgres) CalculateTax(userInfo tax.UserInfo) (tax.Tax, error) {
	return calculate(userInfo)
}

func (p *Postgres) SettingPersonalDeduction(setting tax.Setting) (float64, error) {
	row := p.DB.QueryRow("UPDATE deductions_setting SET amount = $1 WHERE allowance_type = $2 RETURNING amount", setting.Amount, "personal")
	var personalDeduction float64
	err := row.Scan(&personalDeduction)
	if err != nil {
		return 0, err
	}
	return personalDeduction, nil
}
