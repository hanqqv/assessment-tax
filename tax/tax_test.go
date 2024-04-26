// go:build unit

package tax

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type StubTax struct {
	calculateTax Tax
	err          error
}

func (s *StubTax) CalculateTax(userInfo UserInfo) (Tax, error) {
	return s.calculateTax, s.err
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
