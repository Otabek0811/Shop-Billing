package postgres

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	uuid "github.com/google/uuid"
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

func (r *staff_transactionRepo) Create(ctx context.Context, req *models.CreateStaffTransaction) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO staff_transaction(id, transaction_type, description, price,sale_id, source_type, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.TransactionType,
		req.Description,
		req.Price,
		helper.NewNullString(req.SaleID),
		req.SourceType,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *staff_transactionRepo) GetByID(ctx context.Context, req *models.StaffTransactionPrimaryKey) (*models.StaffTransaction, error) {

	var (
		query string

		id               sql.NullString
		transaction_type sql.NullString
		description      sql.NullString
		price            sql.NullString
		sale_id          sql.NullString
		source_type      sql.NullString
		createdAt        sql.NullString
		updatedAt        sql.NullString
	)

	query = `
		SELECT
			id,
			transaction_type,
			description,
			price,
			sale_id,
			source_type,
			created_at,
			updated_at
		FROM staff_transaction
		WHERE id=$1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&transaction_type,
		&description,
		&price,
		&sale_id,
		&source_type,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.StaffTransaction{
		Id:              id.String,
		TransactionType: transaction_type.String,
		Description:     description.String,
		Price:           price.String,
		SaleID:          sale_id.String,
		SourceType:      source_type.String,
		CreatedAt:       createdAt.String,
		UpdatedAt:       updatedAt.String,
	}, nil
}

func (r *staff_transactionRepo) GetList(ctx context.Context, req *models.StaffTransactionGetListRequest) (*models.StaffTransactionGetListResponse, error) {

	var (
		resp   = &models.StaffTransactionGetListResponse{}
		query  string
		where  = " WHERE deleted_at is NULL"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			transaction_type,
			description,
			price,
			sale_id,
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

	if req.SearchStaffID != "" {
		where += ` AND staff_id ILIKE '%' || '` + req.SearchStaffID + `' || '%'`
	}
	if req.SearchSourceType != "" {
		where += ` AND sorce_type ILIKE '%' || '` + req.SearchSourceType + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id               sql.NullString
			transaction_type sql.NullString
			description      sql.NullString
			price            sql.NullString
			sale_id          sql.NullString
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
			Price:           price.String,
			SaleID:          sale_id.String,
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

func (r *staff_transactionRepo) Update(ctx context.Context, req *models.UpdateStaffTransaction) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			staff_transaction
		SET
			'transaction_type' = :transaction_type,
			'description' = :description,
			'price' = :price,
			'sale_id' = :sale_id,
			'source_type' = :source_type,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":               req.Id,
		"transaction_type": req.TransactionType,
		"description":      req.Description,
		"price":            req.Price,
		"sale_id":          req.SaleID,
		"source_type":      req.SourceType,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staff_transactionRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

	var (
		query string
		set   string
	)

	if len(req.Fields) <= 0 {
		return 0, errors.New("no fields")
	}

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s, ", key, key)
	}

	query = `
		UPDATE
			staff_transaction
		SET ` + set + ` updated_at = now()
		WHERE id = :id
	`

	req.Fields["id"] = req.ID

	query, args := helper.ReplaceQueryParams(query, req.Fields)
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staff_transactionRepo) Delete(ctx context.Context, req *models.StaffTransactionPrimaryKey) error {

	_, err := r.db.Exec(ctx,  "Update staff_transaction SET deleted_at=now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
