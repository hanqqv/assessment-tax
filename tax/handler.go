package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	CalculateTax(userInfo UserInfo) (Tax, error)
	SettingPersonalDeduction(setting Setting) (float64, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

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

func (h *Handler) CalculateTaxHandler(c echo.Context) error {
	var userInfo UserInfo
	if err := c.Bind(&userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}
	if err := h.validationUserInfo(userInfo); err.Message != "" {
		return c.JSON(http.StatusBadRequest, err)
	}

	for _, allowance := range userInfo.Allowances {
		if !h.isValidAllowanceType(allowance.AllowanceType) {
			return c.JSON(http.StatusBadRequest, Err{Message: "invalid allowance type"})
		}
	}

	tax, err := h.store.CalculateTax(userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "failed to calculate tax"})
	}

	if tax.Tax < 0.0 {
		refund(&tax)
	}

	return c.JSON(http.StatusOK, tax)

}

func refund(tax *Tax) {
	tax.TaxRefund = tax.Tax * -1
	tax.Tax = 0.0
}

func (h *Handler) SettingPersonalDeductionHandler(c echo.Context) error {
	var setting Setting
	if err := c.Bind(&setting); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}

	if setting.Amount == 0.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "amount is required"})
	}

	if setting.Amount < 10000.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "personal deduction amount must be greater than or equal to 10,000.0"})
	}
	if setting.Amount > 100000.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "personal deduction amount must be less than or equal to 100,000.0"})
	}

	personalDeduction, err := h.store.SettingPersonalDeduction(setting)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "failed to set personal deduction"})
	}

	response := PersonalDeductionResponse{
		PersonalDeduction: personalDeduction,
	}

	return c.JSON(http.StatusOK, response)
}
