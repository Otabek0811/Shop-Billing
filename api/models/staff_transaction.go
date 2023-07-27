package models

type StaffTransactionPrimaryKey struct {
	Id string `json:"id"`
}

type CreateStaffTransaction struct {
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
	Price          string `json:"price"`
	SaleID          string `json:"sale_id"`
	SourceType     string  `json:"source_type"`
}

type StaffTransaction struct {
	Id string `json:"id"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
	Price          string `json:"price"`
	SaleID          string `json:"sale_id"`
	SourceType     string  `json:"source_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
}

type UpdateStaffTransaction struct {
	Id string `json:"id"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
	Price          string `json:"price"`
	SaleID          string `json:"sale_id"`
	SourceType     string  `json:"source_type"`
}

type StaffTransactionGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	SearchSaleID string `json:"search_sale_id"`
	SearchStaffID string `json:"search_staff_id"`
	SearchSourceType string `json:"search_source_type"`
}

type StaffTransactionGetListResponse struct {
	Count             int                 `json:"count"`
	StaffTransactions []*StaffTransaction `json:"staff_transactions"`
}
