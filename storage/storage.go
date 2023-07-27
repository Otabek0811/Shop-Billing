package storage

import (
	"app/api/models"
	"context"
)

type StorageI interface {
	Close()
	Branch() BranchRepoI
	Staff_Tarif() Staff_TarifRepoI
	Staff() StaffRepoI
	Sale()   SaleRepoI
	StaffTransaction() StaffTransactionRepoI
	Biznes_Logic()    BiznesRepoI
}

type BranchRepoI interface {
	Create(context.Context, *models.CreateBranch) (string, error)
	GetByID(context.Context, *models.BranchPrimaryKey) (*models.Branch, error)
	GetList(context.Context, *models.BranchGetListRequest) (*models.BranchGetListResponse, error)
	Update(context.Context, *models.UpdateBranch) (int64, error)
	Patch(context.Context, *models.PatchRequest) (int64, error)
	Delete(context.Context, *models.BranchPrimaryKey) error
}

type Staff_TarifRepoI interface {
	Create(context.Context, *models.CreateStaff_Tarif) (string, error)
	GetByID(context.Context, *models.Staff_TarifPrimaryKey) (*models.Staff_Tarif, error)
	GetList(context.Context, *models.Staff_TarifGetListRequest) (*models.Staff_TarifGetListResponse, error)
	Update(context.Context, *models.UpdateStaff_Tarif) (int64, error)
	Patch(context.Context, *models.PatchRequest) (int64, error)
	Delete(context.Context, *models.Staff_TarifPrimaryKey) error
}

type StaffRepoI interface {
	Create(context.Context, *models.CreateStaff) (string, error)
	GetByID(context.Context, *models.StaffPrimaryKey) (*models.Staff, error)
	GetList(context.Context, *models.StaffGetListRequest) (*models.StaffGetListResponse, error)
	Update(context.Context, *models.UpdateStaff) (int64, error)
	Patch(context.Context, *models.PatchRequest) (int64, error)
	Delete(context.Context, *models.StaffPrimaryKey) error

}



type SaleRepoI interface {
	Create(context.Context, *models.CreateSale) (string, error)
	GetByID(context.Context, *models.SalePrimaryKey) (*models.Sale, error)
	GetList(context.Context, *models.SaleGetListRequest) (*models.SaleGetListResponse, error)
	Update(context.Context, *models.UpdateSale) (int64, error)
	Patch(context.Context, *models.PatchRequest) (int64, error)
	Delete(context.Context, *models.SalePrimaryKey) error
}

type StaffTransactionRepoI interface{
	Create(context.Context, *models.CreateStaffTransaction) (string, error)
	GetByID(context.Context, *models.StaffTransactionPrimaryKey) (*models.StaffTransaction, error)
	GetList(context.Context, *models.StaffTransactionGetListRequest) (*models.StaffTransactionGetListResponse, error)
	Update(context.Context, *models.UpdateStaffTransaction) (int64, error)
	Patch(context.Context, *models.PatchRequest) (int64, error)
	Delete(context.Context, *models.StaffTransactionPrimaryKey) error
}


type BiznesRepoI interface{
	GetTopStaff(ctx context.Context, req *models.TopStaffRequest) (*models.TopStaffResponse, error)
	GetTopBranchByDate(ctx context.Context) (*models.TopBranchReponse, error)
	Do_Staff_Transaction(ctx context.Context, req *models.Sale) error 
}