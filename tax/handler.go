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
	SettingMaxKReceipt(setting Setting) (float64, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
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

	if err := h.validationPersonalDeductionSetting(setting); err.Message != "" {
		return c.JSON(http.StatusBadRequest, err)
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

func (h *Handler) SettingMaxKReceiptHandler(c echo.Context) error {
	var setting Setting
	if err := c.Bind(&setting); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}
	if setting.Amount == 0.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "amount is required"})
	}
	if setting.Amount > 100000.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "max k-receipt amount must be less than or equal to 100,000.0"})
	}
	if setting.Amount < 0.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "max k-receipt amount must be greater than 0.0"})
	}

	maxKReceipt, err := h.store.SettingMaxKReceipt(setting)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "failed to set max k-receipt"})
	}
	response := KReceiptResponse{
		KReceipt: maxKReceipt,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) CalculateTaxCSVHandler(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid file : key must be taxFile"})
	}

	taxResponseCSV, errors := h.processTaxFile(&MultipartFileHeader{file})
	if errors.Message != "" {
		return c.JSON(http.StatusBadRequest, errors)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"taxes": taxResponseCSV})
}
