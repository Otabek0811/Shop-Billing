package postgres

import (
	"app/config"
	"app/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db          *pgxpool.Pool
	branch      *branchRepo
	staff_tarif *staff_tarifRepo
	staff       *staffRepo
	sale        *saleRepo
	staff_transaction *staff_transactionRepo
	b_logic     *biznesRepo
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {
	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}

	connect.MaxConns = cfg.PostgresMaxConnection

	pgxpool, err := pgxpool.ConnectConfig(context.Background(), connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (s *store) Close() {
	s.db.Close()
}

func (s *store) Branch() storage.BranchRepoI {

	if s.branch == nil {
		s.branch = NewBranchRepo(s.db)
	}

	return s.branch
}

func (s *store) Staff_Tarif() storage.Staff_TarifRepoI {

	if s.staff_tarif == nil {
		s.staff_tarif = NewStaff_TarifRepo(s.db)
	}

	return s.staff_tarif
}
func (s *store) Staff() storage.StaffRepoI {

	if s.staff == nil {
		s.staff = NewStaffRepo(s.db)
	}

	return s.staff
}

func (s *store) Sale() storage.SaleRepoI {

	if s.sale == nil {
		s.sale = NewSaleRepo(s.db)
	}

	return s.sale
}
func (s *store) StaffTransaction() storage.StaffTransactionRepoI {

	if s.staff_transaction == nil {
		s.staff_transaction = NewStaffTransactionRepo(s.db)
	}

	return s.staff_transaction
}

func (s *store) Biznes_Logic() storage.BiznesRepoI {

	if s.b_logic == nil {
		s.b_logic = NewBiznesRepo(s.db)
	}

	return s.b_logic
}

