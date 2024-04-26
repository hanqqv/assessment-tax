package tax

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
)

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
		return nil, Err{Message: "error reading file"}
	}

	var taxResponseCSV []TaxResponseCSV
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, Err{Message: "error reading file: invalid format"}
		}

		taxResponse, error := h.processCSVLine(line)
		if error.Message != "" {
			return nil, Err{Message: error.Message}
		}

		taxResponseCSV = append(taxResponseCSV, taxResponse)
	}

	return taxResponseCSV, Err{}
}

func (h *Handler) processCSVLine(line []string) (TaxResponseCSV, Err) {
	userInfo, err := parseUserInfoFromCSVLine(line)
	if err != nil {
		return TaxResponseCSV{}, Err{Message: err.Error()}
	}

	if err := h.validationUserInfo(userInfo); err.Message != "" {
		return TaxResponseCSV{}, Err{Message: err.Message}
	}

	tax, err := h.store.CalculateTax(userInfo)
	if err != nil {
		return TaxResponseCSV{}, Err{Message: err.Error()}
	}

	if tax.Tax < 0.0 {
		refund(&tax)
	}

	return TaxResponseCSV{
		TotalIncome: userInfo.TotalIncome,
		Tax:         tax.Tax,
		TaxRefund:   tax.TaxRefund,
	}, Err{}
}

func parseUserInfoFromCSVLine(line []string) (UserInfo, error) {
	if len(line) != 3 {
		return UserInfo{}, fmt.Errorf("invalid file format")
	}

	totalIncomeStr, whtStr, donationStr := strings.TrimSpace(line[0]), strings.TrimSpace(line[1]), strings.TrimSpace(line[2])

	if totalIncomeStr == "" || whtStr == "" || donationStr == "" {
		return UserInfo{}, fmt.Errorf("totalIncome, wht or donation value can not be empty")
	}

	totalIncome, err := strconv.ParseFloat(totalIncomeStr, 64)
	if err != nil {
		return UserInfo{}, fmt.Errorf("totalIncome must be a numeric value")
	}

	wht, err := strconv.ParseFloat(whtStr, 64)
	if err != nil {
		return UserInfo{}, fmt.Errorf("wht must be a numeric value")
	}

	donation, err := strconv.ParseFloat(donationStr, 64)
	if err != nil {
		return UserInfo{}, fmt.Errorf("donation must be a numeric value")
	}

	return UserInfo{
		TotalIncome: totalIncome,
		WHT:         wht,
		Allowances:  []Allowances{{AllowanceType: "donation", Amount: donation}},
	}, nil
}
