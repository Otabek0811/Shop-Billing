package postgres

import (
	"app/api/models"
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
)

type staff_transactionRepo struct {
	db *pgxpool.Pool
}

func NewStaffTransactionRepo(db *pgxpool.Pool) *staff_transactionRepo {
	return &staff_transactionRepo{
		db: db,
	}
}

func (r *staff_transactionRepo) GetByID(ctx context.Context, req *models.StaffTransactionPrimaryKey) (*models.StaffTransaction, error) {

	var (
		query string

		id               sql.NullString
		transaction_type sql.NullString
		description      sql.NullString
		price            sql.NullFloat64
		sale_id          sql.NullString
		staff_id         sql.NullString
		source_type      sql.NullString
		createdAt        sql.NullString
		updatedAt        sql.NullString
		staff_name       sql.NullString
	)

	query = `
		SELECT
			st.id,
			st.transaction_type,
			st.description,
			st.price,
			st.sale_id,
			st.staff_id,
			st.source_type,
			st.created_at,
			st.updated_at,
			s.name
		FROM staff_transaction as st
		join staff as s on s.id = st.staff_id
		WHERE st.id=$1 
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&transaction_type,
		&description,
		&price,
		&sale_id,
		&staff_id,
		&source_type,
		&createdAt,
		&updatedAt,
		&staff_name,
	)

	if err != nil {
		return nil, err
	}

	return &models.StaffTransaction{
		Id:              id.String,
		TransactionType: transaction_type.String,
		Description:     description.String,
		Price:           price.Float64,
		SaleID:          sale_id.String,
		StaffID:         staff_id.String,
		SourceType:      source_type.String,
		CreatedAt:       createdAt.String,
		UpdatedAt:       updatedAt.String,
		StaffName:       staff_id.String,
	}, nil
}

func (r *staff_transactionRepo) GetList(ctx context.Context, req *models.StaffTransactionGetListRequest) (*models.StaffTransactionGetListResponse, error) {

	var (
		resp    = &models.StaffTransactionGetListResponse{}
		query   string
		where   = " WHERE deleted_at is NULL"
		offset  = " OFFSET 0"
		limit   = " LIMIT 10"
		orderBY = " ORDER BY"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			transaction_type,
			description,
			price,
			sale_id,
			staff_id,
			source_type,
			created_at,
			updated_at
		FROM staff_transaction
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.SearchSaleID != "" {
		where += ` AND sale_id ILIKE '%' || '` + req.SearchSaleID + `' || '%'`
	}

	if req.SearchSourceType != "" {
		where += ` AND source_type ILIKE '%' || '` + req.SearchSourceType + `' || '%'`
	}
	if req.SearchSourceType != "" {
		where += ` AND staff_id = '` + req.SearchSourceType + `'`
	}

	if req.SortPriceType != "" {
		orderBY += `  price ` + req.SortPriceType
		query += where + orderBY + offset + limit
	} else {
		query += where + offset + limit
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id               sql.NullString
			transaction_type sql.NullString
			description      sql.NullString
			price            sql.NullFloat64
			sale_id          sql.NullString
			staff_id         sql.NullString
			source_type      sql.NullString
			createdAt        sql.NullString
			updatedAt        sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&transaction_type,
			&description,
			&price,
			&sale_id,
			&staff_id,
			&source_type,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.StaffTransactions = append(resp.StaffTransactions, &models.StaffTransaction{
			Id:              id.String,
			TransactionType: transaction_type.String,
			Description:     description.String,
			Price:           price.Float64,
			SaleID:          sale_id.String,
			StaffID:         staff_id.String,
			SourceType:      source_type.String,
			CreatedAt:       createdAt.String,
			UpdatedAt:       updatedAt.String,
		})
	}

	sort.Slice(resp.StaffTransactions[:], func(i, j int) bool {
		return resp.StaffTransactions[i].CreatedAt > resp.StaffTransactions[j].CreatedAt
	})

	return resp, nil
}
