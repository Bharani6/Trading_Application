package profile

type KYCSubmitRequest struct {
	Mobile       string           `json:"mobile"`
	FatherName   string           `json:"father_name"`
	MotherName   string           `json:"mother_name"`
	Country      string           `json:"country"`
	State        string           `json:"state"`
	City         string           `json:"city"`
	Address      string           `json:"address"`
	Pincode      string           `json:"pincode"`
	IncomeRange  string           `json:"income_range"`
	BankAccounts []BankAccountDTO `json:"bank_accounts" binding:"dive"`
	Nominees     []NomineeDTO     `json:"nominees" binding:"omitempty,max=3,dive"`
}

type BankAccountDTO struct {
	AccountType   string `json:"account_type" binding:"required"`
	IFSC          string `json:"ifsc" binding:"required,len=11"`
	BankName      string `json:"bank_name" binding:"required"`
	Branch        string `json:"branch"`
	AccountNumber string `json:"account_number" binding:"required,numeric"`
}

type NomineeDTO struct {
	Name                 string  `json:"name" binding:"required"`
	DOB                  string  `json:"dob" binding:"required"`
	PAN                  string  `json:"pan" binding:"required"`
	Relationship         string  `json:"relationship" binding:"required"`
	Percentage           float64 `json:"percentage" binding:"required,gt=0"`
	GuardianName         string  `json:"guardian_name"`
	GuardianRelationship string  `json:"guardian_relationship"`
	GuardianPAN          string  `json:"guardian_pan"`
	GuardianDOB          string  `json:"guardian_dob"`
}
