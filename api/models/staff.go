package models

type StaffPrimaryKey struct {
	Id string `json:"id"`
}

type CreateStaff struct {
	Name      string  `json:"name"`
	StaffType string  `json:"staff_type"`
	Balance   float64 `json:"balance"`
	TarifID   string  `json:"tarif_id"`
	BranchID  string  `json:"branch_id"`
}

type Staff struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	StaffType  string      `json:"staff_type"`
	Balance    float64     `json:"balance"`
	TarifID    string      `json:"tarif_id"`
	BranchID   string      `json:"branch_id"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
	DeletedAt  string      `json:"deleted_at"`
	TarifData  *Staff_Tarif `json:"tarif_data"`
	BranchData *Branch      `json:"branch_data"`
}

type UpdateStaff struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	StaffType string  `json:"staff_type"`
	TarifID   string  `json:"tarif_id"`
	BranchID  string  `json:"branch_id"`
	Balance   float64 `json:"balance"`
}

type StaffGetListRequest struct {
	Offset            int    `json:"offset"`
	Limit             int    `json:"limit"`
	SearchByName      string `json:"search_by_name"`
	SearchBYStaffType string `json:"search_by_staff_type"`
	SearchByTarifID   string `json:"search_by_tarif_id"`
	SearchByBranchID  string `json:"search_by_branch_id"`
	SearchBalanceFrom string `json:"search_balance_from"`
	SearchBalanceTo   string `json:"search_balance_to"`
}

type StaffGetListResponse struct {
	Count  int      `json:"count"`
	Staffs []*Staff `json:"Staffs"`
}
