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
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund,omitempty"`
}