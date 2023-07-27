package models

type BranchPrimaryKey struct {
	Id string `json:"id"`
}

type CreateBranch struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	CompanyID string `json:"company_id"`
}

type Branch struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CompanyID string `json:"company_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	CompanyName string  `json:"company_name"`
}

type UpdateBranch struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CompanyID string `json:"company_id"`
}

type BranchGetListRequest struct {
	Offset          int    `json:"offset"`
	Limit           int    `json:"limit"`
	SearchBYName    string `json:"search_by_name"`
	SearchByAddress string `json:"search_by_address"`
}

type BranchGetListResponse struct {
	Count   int       `json:"count"`
	Branchs []*Branch `json:"Branchs"`
}
