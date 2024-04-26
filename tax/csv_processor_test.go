// go:build unit

package tax

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockFileOpener struct {
	reader *strings.Reader
}

func (m *MockFileOpener) Open() (io.ReadCloser, error) {
	return io.NopCloser(m.reader), nil
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
