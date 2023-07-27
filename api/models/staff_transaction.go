package models

type StaffTransactionPrimaryKey struct {
	Id string `json:"id"`
}


type StaffTransaction struct {
	Id              string  `json:"id"`
	TransactionType string  `json:"transaction_type"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	SourceType      string  `json:"source_type"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       string  `json:"deleted_at"`
	StaffName       string  `json:"staff_name"`
}



type StaffTransactionGetListRequest struct {
	Offset           int    `json:"offset"`
	Limit            int    `json:"limit"`
	SearchSaleID     string `json:"search_sale_id"`
	SearchStaffID    string `json:"search_staff_id"`
	SearchSourceType string `json:"search_source_type"`
	SortPriceType    string `json:"sort_price_type"`
}

type StaffTransactionGetListResponse struct {
	Count             int                 `json:"count"`
	StaffTransactions []*StaffTransaction `json:"staff_transactions"`
}
