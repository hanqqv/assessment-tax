package tax

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) CalculateTaxCSVHandler(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid file"})
	}

	taxResponseCSV, errors := h.processTaxFile(&MultipartFileHeader{file})
	if errors.Message != "" {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"taxes": taxResponseCSV})
}

type MultipartFileHeader struct {
	*multipart.FileHeader
}

func (m *MultipartFileHeader) Open() (io.ReadCloser, error) {
	return m.FileHeader.Open()
}

type FileOpener interface {
	Open() (io.ReadCloser, error)
}

func (h *Handler) processTaxFile(file FileOpener) ([]TaxResponseCSV, Err) {
	src, err := file.Open()
	if err != nil {
		return nil, Err{Message: "failed to open file"}
	}
	defer src.Close()

	reader := csv.NewReader(src)

	_, err = reader.Read()
	if err != nil {
		return nil, Err{Message: "failed to read file"}
	}

	var taxResponseCSV []TaxResponseCSV
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, Err{Message: "failed to read file"}
		}

		userInfo, err := parseUserInfoFromCSVLine(line)
		if err != nil {
			return nil, Err{Message: err.Error()}
		}

		if err := h.validationUserInfo(userInfo); err.Message != "" {
			return nil, err
		}

		tax, err := h.store.CalculateTax(userInfo)
		if err != nil {
			return nil, Err{Message: err.Error()}
		}

		if tax.Tax < 0.0 {
			refund(&tax)
		}

		taxResponseCSV = append(taxResponseCSV, TaxResponseCSV{
			TotalIncome: userInfo.TotalIncome,
			Tax:         tax.Tax,
		})
	}
	return taxResponseCSV, Err{}
}

func parseUserInfoFromCSVLine(line []string) (UserInfo, error) {
	if len(line) != 3 {
		return UserInfo{}, errors.New("invalid file format")
	}

	totalIncomeStr, whtStr, donationStr := strings.TrimSpace(line[0]), strings.TrimSpace(line[1]), strings.TrimSpace(line[2])

	if totalIncomeStr == "" || whtStr == "" || donationStr == "" {
		return UserInfo{}, errors.New("invalid file format")
	}

	totalIncome, err := strconv.ParseFloat(totalIncomeStr, 64)
	if err != nil {
		return UserInfo{}, errors.New("totalIncome must be numeric value")
	}

	wht, err := strconv.ParseFloat(whtStr, 64)
	if err != nil {
		return UserInfo{}, errors.New("wht must be numeric value")
	}

	donation, err := strconv.ParseFloat(donationStr, 64)
	if err != nil {
		return UserInfo{}, errors.New("donation must be numeric value")
	}

	return UserInfo{
		TotalIncome: totalIncome,
		WHT:         wht,
		Allowances:  []Allowances{{AllowanceType: "donation", Amount: donation}},
	}, nil
}
