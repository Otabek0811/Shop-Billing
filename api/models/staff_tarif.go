package models

type Staff_TarifPrimaryKey struct {
	Id string `json:"id"`
}

type CreateStaff_Tarif struct {
	Name          string `json:"name"`
	TypeTarif     string `json:"type_tarif"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
}

type Staff_Tarif struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TypeTarif     string `json:"type_tarif"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
}

type UpdateStaff_Tarif struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TypeTarif     string `json:"type_tarif"`
	AmountForCash float64 `json:"amount_for_cash"`
	AmountForCard float64 `json:"amount_for_card"`
}

type Staff_TarifGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type Staff_TarifGetListResponse struct {
	Count        int            `json:"count"`
	Staff_Tarifs []*Staff_Tarif `json:"Staff_Tarifs"`
}
