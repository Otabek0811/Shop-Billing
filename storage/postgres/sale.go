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

type saleRepo struct {
	db *pgxpool.Pool
}

func NewSaleRepo(db *pgxpool.Pool) *saleRepo {
	return &saleRepo{
		db: db,
	}
}

func (r *saleRepo) Create(ctx context.Context, req *models.CreateSale) (string, error) {

	trx, err := r.db.Begin(ctx)
	if err != nil {
		return "", nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO sale(id, branch_id, assistant_id, cashier_id, price, payment_type, status, client_name, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`

	_, err = trx.Exec(ctx, query,
		id,
		helper.NewNullString(req.BranchID),
		helper.NewNullString(req.AssistentID),
		helper.NewNullString(req.CashierID),
		req.Price,
		req.PaymentType,
		req.Status,
		req.ClientName,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *saleRepo) GetByID(ctx context.Context, req *models.SalePrimaryKey) (*models.Sale, error) {

	var (
		query string

		id            sql.NullString
		branch_id     sql.NullString
		assistant_id  sql.NullString
		cashier_id    sql.NullString
		price         sql.NullFloat64
		payment_type  sql.NullString
		status        sql.NullString
		client_name   sql.NullString
		createdAt     sql.NullString
		updatedAt     sql.NullString
		branchName    sql.NullString
		assistantName sql.NullString
		cashierName   sql.NullString
	)

	query = `
		SELECT
			s.id,
			s.branch_id,
			s.assistant_id,
			s.cashier_id,
			s.price,
			s.payment_type,
			s.status,
			s.client_name,
			s.created_at,
			s.updated_at,
			b.name,
			sta.name,
			stc.name
		FROM sale as s
		LEFT JOIN branch AS b ON b.id=s.branch_id
		LEFT JOIN staff AS sta ON sta.id = s.assistant_id
		LEFT JOIN staff AS stc ON stc.id = s.cashier_id
		WHERE s.id=$1 and s.deleted_at is NULL
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&branch_id,
		&assistant_id,
		&cashier_id,
		&price,
		&payment_type,
		&status,
		&client_name,
		&createdAt,
		&updatedAt,
		&branchName,
		&assistantName,
		&cashierName,
	)

	if err != nil {
		return nil, err
	}

	return &models.Sale{
		Id:            id.String,
		BranchID:      branch_id.String,
		AssistentID:   assistant_id.String,
		CashierID:     cashier_id.String,
		Price:         price.Float64,
		PaymentType:   payment_type.String,
		Status:        status.String,
		ClientName:    client_name.String,
		CreatedAt:     createdAt.String,
		UpdatedAt:     updatedAt.String,
		BranchName:    branchName.String,
		CashiertName:  cashierName.String,
		AssistantName: assistantName.String,
	}, nil
}

func (r *saleRepo) GetList(ctx context.Context, req *models.SaleGetListRequest) (*models.SaleGetListResponse, error) {

	var (
		resp    = &models.SaleGetListResponse{}
		query   string
		where   = " WHERE deleted_at is NULL"
		offset  = " OFFSET 0"
		limit   = " LIMIT 10"
		orderBY = " "
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			branch_id,
			assistant_id,
			cashier_id,
			price,
			payment_type,
			status,
			client_name,
			created_at,
			updated_at
		FROM sale
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.SearchBranchID != "" {
		where += ` AND branch_id ILIKE '%' || '` + req.SearchBranchID + `' || '%'`
	}
	if req.SearchClientName != "" {
		where += ` AND client_name ILIKE '%' || '` + req.SearchClientName + `' || '%'`
	}
	if req.SearchPaymentType != "" {
		where += ` AND payment_type ILIKE '%' || '` + req.SearchPaymentType + `' || '%'`
	}
	if req.SearchAssistantID != "" {
		where += ` AND assistant_id ILIKE '%' || '` + req.SearchAssistantID + `' || '%'`
	}
	if req.SearchCashierID != "" {
		where += ` AND cashier_id ILIKE '%' || '` + req.SearchCashierID + `' || '%'`
	}
	if req.SearchStatus != "" {
		where += ` AND status ILIKE '%' || '` + req.SearchStatus + `' || '%'`
	}
	if req.SearchCreatedAtFrom != "" {
		where += ` AND date(created_at) >= || '` + req.SearchCreatedAtFrom + `'`
	}
	if req.SearchCreatedAtTo != "" {
		where += ` AND date(created_at) <= || '` + req.SearchCreatedAtTo + `'`
	}
	if req.SortPriceType != "" {
		orderBY += ` price '` + req.SortPriceType + `'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id           sql.NullString
			branch_id    sql.NullString
			assistant_id sql.NullString
			cashier_id   sql.NullString
			price        sql.NullFloat64
			payment_type sql.NullString
			status       sql.NullString
			client_name  sql.NullString
			createdAt    sql.NullString
			updatedAt    sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&branch_id,
			&assistant_id,
			&cashier_id,
			&price,
			&payment_type,
			&status,
			&client_name,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Sales = append(resp.Sales, &models.Sale{
			Id:          id.String,
			BranchID:    branch_id.String,
			AssistentID: assistant_id.String,
			CashierID:   cashier_id.String,
			Price:       price.Float64,
			PaymentType: payment_type.String,
			Status:      status.String,
			ClientName:  client_name.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}

	sort.Slice(resp.Sales[:], func(i, j int) bool {
		return resp.Sales[i].CreatedAt > resp.Sales[j].CreatedAt
	})
	return resp, nil
}

func (r *saleRepo) Update(ctx context.Context, req *models.UpdateSale) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			sale
		SET
			branch_id = :branch_id,
			assistant_id = :assistant_id,
			cashier_id = :cashier_id,
			price   = :price,
			payment_type = : payment_type,
			status = :status,
			client_name = :client_name,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"branch_id":    req.BranchID,
		"assistant_id": req.AssistentID,
		"cashier_id":   req.CashierID,
		"price":        req.Price,
		"payment_type": req.PaymentType,
		"status":       req.Status,
		"client_name":  req.ClientName,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *saleRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

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
			sale
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

func (r *saleRepo) Delete(ctx context.Context, req *models.SalePrimaryKey) error {

	_, err := r.db.Exec(ctx,  "Update sale SET deleted_at=now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
