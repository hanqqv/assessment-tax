package tax

func (h *Handler) validationUserInfo(userInfo UserInfo) Err {
	if userInfo.TotalIncome == 0.0 {
		return Err{Message: "total income is required"}
	}
	if userInfo.TotalIncome < 0.0 {
		return Err{Message: "total income must be greater than 0.0"}
	}
	if userInfo.WHT < 0.0 {
		return Err{Message: "wht must be greater than or equal to 0.0"}
	}
	if userInfo.WHT > userInfo.TotalIncome {
		return Err{Message: "wht must be less than or equal to total income"}
	}
	for _, allowance := range userInfo.Allowances {
		if allowance.AllowanceType == "" {
			return Err{Message: "missing allowanceType key"}
		}
		if allowance.Amount < 0.0 {
			return Err{Message: "allowance amount must be greater than or equal to 0.0"}
		}
		if allowance.AllowanceType == "personal" {
			return Err{Message: "user can not fill personal allowance"}
		}
	}

	return Err{}
}

func (h *Handler) isValidAllowanceType(allowanceType string) bool {
	validAllowanceTypes := map[string]bool{
		"donation":  true,
		"k-receipt": true,
		"personal":  true,
	}

	_, ok := validAllowanceTypes[allowanceType]
	return ok
}

func (h *Handler) validationPersonalDeductionSetting(setting Setting) Err {
	if setting.Amount == 0.0 {
		return Err{Message: "amount is required"}
	}
	if setting.Amount < 10000.0 {
		return Err{Message: "personal deduction amount must be greater than or equal to 10,000.0"}
	}
	if setting.Amount > 100000.0 {
		return Err{Message: "personal deduction amount must be less than or equal to 100,000.0"}
	}

	return Err{}
}
