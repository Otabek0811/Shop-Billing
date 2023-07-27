package models

type SalePrimaryKey struct {
	Id string `json:"id"`
}

type CreateSale struct {
	AssistentID string  `json:"assistent_id"`
	CashierID   string  `json:"cashier_id"`
	BranchID    string  `json:"branch_id"`
	Price       float64 `json:"price"`
	PaymentType string  `json:"payment_type"`
	Status      string  `json:"status"`
	ClientName  string  `json:"client_name"`
}

type Sale struct {
	Id          string  `json:"id"`
	AssistentID string  `json:"assistent_id"`
	CashierID   string  `json:"cashier_id"`
	BranchID    string  `json:"branch_id"`
	Price       float64 `json:"price"`
	PaymentType string  `json:"payment_type"`
	Status      string  `json:"status"`
	ClientName  string  `json:"client_name"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   string  `json:"deleted_at"`
	BranchName string `json:"branch_name"`
	AssistantName  string  `json:"assistant_name"`
	CashiertName string  `json:"cashier_name"`
}


type UpdateSale struct {
	Id          string  `json:"id"`
	AssistentID string  `json:"assistent_id"`
	CashierID   string  `json:"cashier_id"`
	BranchID    string  `json:"branch_id"`
	Price       float64 `json:"price"`
	PaymentType string  `json:"payment_type"`
	Status      string  `json:"status"`
	ClientName  string  `json:"client_name"`
}

type SaleGetListRequest struct {
	Offset              int    `json:"offset"`
	Limit               int    `json:"limit"`
	SearchBranchID      string `json:"search_branch_id"`
	SearchPaymentType   string `json:"search_payment_type"`
	SearchClientName    string `json:"search_client_name"`
	SearchAssistantID   string `json:"search_assistant_id"`
	SearchCashierID     string `json:"search_cashier_id"`
	SearchStatus        string `json:"search_status"`
	SearchCreatedAtFrom string `json:"search_created_at_from"`
	SearchCreatedAtTo   string `json:"search_created_at_to"`
	SortPriceType    string    	`json:"sort_price_type"`
}

type SaleGetListResponse struct {
	Count int     `json:"count"`
	Sales []*Sale `json:"Sales"`
}

