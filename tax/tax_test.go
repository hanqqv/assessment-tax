// go:build unit

package tax

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type StubTax struct {
	calculateTax             Tax
	settingPersonalDeduction float64
	err                      error
}

type MockFileOpener struct {
	reader *strings.Reader
}

func (m *MockFileOpener) Open() (io.ReadCloser, error) {
	return io.NopCloser(m.reader), nil
}

func (s *StubTax) CalculateTax(userInfo UserInfo) (Tax, error) {
	return s.calculateTax, s.err
}

func (s *StubTax) SettingPersonalDeduction(setting Setting) (float64, error) {
	return s.settingPersonalDeduction, s.err
}

func TestCalculateTax(t *testing.T) {
	t.Run("given user unable to calculate tax should return status 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "failed to calculate tax" }`

		stubTax := StubTax{err: errors.New("failed to calculate tax")}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code, "expected status code %d but got %d", http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given user able to calculate tax should return status 200 and tax", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := Tax{Tax: 29000.0}

		stubTax := StubTax{calculateTax: want}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code %d but got %d", http.StatusOK, rec.Code)

		var got Tax
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, want, got, "expected tax %v but got %v", want, got)
	})
	t.Run("given invalid request body should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": .0, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "invalid request body" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given invalid allowance type should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 0.0, "allowances": [{"allowanceType": "invalid", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "invalid allowance type" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given missing total income should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "total income is required" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given negative total income should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": -1000.0, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "total income must be greater than 0.0" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given negative WHT should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": -1000.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "wht must be greater than or equal to 0.0" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given WHT greater than total income should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 6000000.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "wht must be less than or equal to total income" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given missing allowance type should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 0.0, "allowances": [{"amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "missing allowanceType key" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given negative allowance amount should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 0.0, "allowances": [{"allowanceType": "donation", "amount": -1000.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "allowance amount must be greater than or equal to 0.0" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given tax refund is negative should return status 200 and tax result", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 500000.0, "wht": 30000.0, "allowances": [{"allowanceType": "donation", "amount": 0.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := Tax{Tax: 0.0, TaxRefund: 1000.0}

		stubTax := StubTax{calculateTax: Tax{Tax: -1000.0}}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code %d but got %d", http.StatusOK, rec.Code)

		var got Tax
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, want, got, "expected tax %v but got %v", want, got)
	})
	t.Run("given user fill personal allowance should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", io.NopCloser(strings.NewReader(`{"totalIncome": 5000000.0, "wht": 0.0, "allowances": [{"allowanceType": "personal", "amount": 100000.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tax/calculations")

		want := `{ "message": "user can not fill personal allowance" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
}
func TestSettingPersonalDeduction(t *testing.T) {
	t.Run("given user unable to set personal deduction should return status 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": 70000.0}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/deductions/personal")

		want := `{ "message": "failed to set personal deduction" }`

		stubTax := StubTax{err: errors.New("failed to set personal deduction")}
		p := New(&stubTax)

		err := p.SettingPersonalDeductionHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code, "expected status code %d but got %d", http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given admin able to set personal deduction should return status 200 and personal deduction", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": 70000.0}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/deductions/personal")

		want := PersonalDeductionResponse{PersonalDeduction: 70000.0}

		stubTax := StubTax{settingPersonalDeduction: 70000.0}
		p := New(&stubTax)

		err := p.SettingPersonalDeductionHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code %d but got %d", http.StatusOK, rec.Code)

		var got PersonalDeductionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, want, got, "expected personal deduction %v but got %v", want, got)
	})
	t.Run("given missing amount should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/deductions/personal")

		want := `{ "message": "amount is required" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.SettingPersonalDeductionHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given amount less than 10,000 should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": 5000.0}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/deductions/personal")

		want := `{ "message": "personal deduction amount must be greater than or equal to 10,000.0" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.SettingPersonalDeductionHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given amount greater than 100,000 should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": 150000.0}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/deductions/personal")

		want := `{ "message": "personal deduction amount must be less than or equal to 100,000.0" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.SettingPersonalDeductionHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
	t.Run("given invalid request body should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": "invalid"}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/admin/deductions/personal")

		want := `{ "message": "invalid request body" }`

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.SettingPersonalDeductionHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, want, rec.Body.String(), "expected response body %s but got %s", want, rec.Body.String())
	})
}
func TestProcessCSVLine(t *testing.T) {
	tests := []struct {
		name    string
		line    []string
		want    TaxResponseCSV
		wantErr bool
	}{
		{
			name: "valid line",
			line: []string{"100000", "0", "5000"},
			want: TaxResponseCSV{
				TotalIncome: 100000.0,
				Tax:         0.0,
				TaxRefund:   0.0,
			},
			wantErr: false,
		},
		{
			name:    "invalid line",
			line:    []string{"invalid", "50000", "5000"},
			wantErr: true,
		},
	}

	stubTax := StubTax{}
	p := New(&stubTax)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.processCSVLine(tt.line)
			if (err.Message != "") != tt.wantErr {
				t.Errorf("processCSVLine() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got, "expected tax %v but got %v", tt.want, got)
		})
	}
}
func TestProcessTaxFile(t *testing.T) {
	t.Run("given valid file format should return taxResponseCSV", func(t *testing.T) {
		stubTax := StubTax{}
		p := New(&stubTax)
		mockFile := &MockFileOpener{reader: strings.NewReader("totalIncome,wht,donation\n10000,2000,3000\n15000,2500,3500\n")}
		want := []TaxResponseCSV{
			{TotalIncome: 10000, Tax: 0, TaxRefund: 0},
			{TotalIncome: 15000, Tax: 0, TaxRefund: 0},
		}

		got, err := p.processTaxFile(mockFile)
		assert.Equal(t, err.Message, "", "expected no error but got %v", err.Message)
		assert.Equal(t, want, got, "expected tax %v but got %v", want, got)
	})
	t.Run("given invalid file format should return error message", func(t *testing.T) {
		stubTax := StubTax{}
		p := New(&stubTax)
		mockFile := &MockFileOpener{reader: strings.NewReader("totalIncome,wht,donation\n10000,2000\n15000,2500,3500\n")}
		want := Err{Message: "error reading file: invalid format"}

		got, err := p.processTaxFile(mockFile)
		assert.Equal(t, want, err, "expected error %v but got %v", want, err)
		assert.Nil(t, got, "expected nil but got %v", got)
	})
}

func TestCalculateTaxCSVHandler(t *testing.T) {
	t.Run("given valid file format should return status 200 and tax result", func(t *testing.T) {
		e := echo.New()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte("totalIncome,wht,donation\n1000,200,300\n2000,400,600"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxCSVHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code %d but got %d", http.StatusOK, rec.Code)
		assert.Equal(t, `{"taxes":[{"totalIncome":1000,"tax":0},{"totalIncome":2000,"tax":0}]}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
	t.Run("given invalid file format should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte("totalIncome,wht,donation\n1000,200,300\n2000,400"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxCSVHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"message":"error reading file: invalid format"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
	t.Run("given invalid file key should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("invalidKey", "taxes.csv")
		part.Write([]byte("totalIncome,wht,donation\n1000,200,300\n2000,400,600"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxCSVHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"message":"invalid file : key must be taxFile"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
	t.Run("given totalIncome not numeric should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte("totalIncome,wht,donation\n1000,200,300\n2000,400,600\ninvalid,200,300"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxCSVHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"message":"totalIncome must be a numeric value"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
	t.Run("given wht not numeric should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte("totalIncome,wht,donation\n1000,200,300\n2000,400,600\n1000,invalid,300"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxCSVHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"message":"wht must be a numeric value"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
	t.Run("given totalIncome negative should return status 400 and error message", func(t *testing.T) {
		e := echo.New()
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("taxFile", "taxes.csv")
		part.Write([]byte("totalIncome,wht,donation\n-1000,200,300\n2000,400,600"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubTax := StubTax{}
		p := New(&stubTax)

		err := p.CalculateTaxCSVHandler(c)

		assert.NoError(t, err, "expected no error but got %v", err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "expected status code %d but got %d", http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"message":"total income must be greater than 0.0"}`, strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
func TestParseUserInfoFromCSVLine(t *testing.T) {
	tests := []struct {
		name    string
		line    []string
		want    UserInfo
		wantErr bool
	}{
		{
			name:    "valid input",
			line:    []string{"500000", "5000", "5000"},
			want:    UserInfo{TotalIncome: 500000.0, WHT: 5000.0, Allowances: []Allowances{{AllowanceType: "donation", Amount: 5000.0}}},
			wantErr: false,
		},
		{
			name:    "invalid input",
			line:    []string{"abc", "5000", "5000"},
			want:    UserInfo{},
			wantErr: true,
		},
		{
			name:    "missing values",
			line:    []string{"5000", "5000"},
			want:    UserInfo{},
			wantErr: true,
		},
		{
			name:    "empty values",
			line:    []string{"", "5000", "5000"},
			want:    UserInfo{},
			wantErr: true,
		},
		{
			name:    "empty values",
			line:    []string{"5000", "", "5000"},
			want:    UserInfo{},
			wantErr: true,
		},
		{
			name:    "wht not numeric",
			line:    []string{"5000", "abc", "5000"},
			want:    UserInfo{},
			wantErr: true,
		},
		{
			name:    "donation not numeric",
			line:    []string{"5000", "5000", "abc"},
			want:    UserInfo{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUserInfoFromCSVLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUserInfoFromCSVLine() got error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got, "expected tax %v but got %v", tt.want, got)
		})
	}
}
