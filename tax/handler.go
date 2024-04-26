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
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

func (h *Handler) isValidAllowanceType(allowanceType string) bool {
	validAllowanceTypes := map[string]bool{
		"donation":  true,
		"k-receipt": true,
	}

	_, ok := validAllowanceTypes[allowanceType]
	return ok
}

func (h *Handler) CalculateTaxHandler(c echo.Context) error {
	var userInfo UserInfo
	if err := c.Bind(&userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}
	if userInfo.TotalIncome == 0.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "total income is required"})
	}

	if userInfo.TotalIncome < 0.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "total income must be greater than 0.0"})
	}
	if userInfo.WHT < 0.0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "wht must be greater than or equal to 0.0"})
	}
	if userInfo.WHT > userInfo.TotalIncome {
		return c.JSON(http.StatusBadRequest, Err{Message: "wht must be less than or equal to total income"})
	}

	for _, allowance := range userInfo.Allowances {
		if !h.isValidAllowanceType(allowance.AllowanceType) {
			return c.JSON(http.StatusBadRequest, Err{Message: "invalid allowance type"})
		}
		if allowance.Amount < 0.0 {
			return c.JSON(http.StatusBadRequest, Err{Message: "allowance amount must be greater than or equal to 0.0"})
		}
	}

	tax, err := h.store.CalculateTax(userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "failed to calculate tax"})
	}
	return c.JSON(http.StatusOK, tax)

}
