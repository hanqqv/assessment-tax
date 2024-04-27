package tax

type UserInfo struct {
	TotalIncome float64      `json:"totalIncome"`
	WHT         float64      `json:"wht"`
	Allowances  []Allowances `json:"allowances"`
}

type Allowances struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type Tax struct {
	Tax       float64    `json:"tax"`
	TaxRefund float64    `json:"taxRefund,omitempty"`
	TaxLevel  []TaxLevel `json:"taxLevel"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type Setting struct {
	Amount float64 `json:"amount"`
}

type PersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

type TaxResponseCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}
