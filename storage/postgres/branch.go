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

type branchRepo struct {
	db *pgxpool.Pool
}

func NewBranchRepo(db *pgxpool.Pool) *branchRepo {
	return &branchRepo{
		db: db,
	}
}

func (r *branchRepo) Create(ctx context.Context, req *models.CreateBranch) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO branch(id, name, address, company_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Address,
		helper.NewNullString(req.CompanyID),
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *branchRepo) GetByID(ctx context.Context, req *models.BranchPrimaryKey) (*models.Branch, error) {

	var (
		query string

		id          sql.NullString
		name        sql.NullString
		address     sql.NullString
		phoneNumber sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	query = `
		SELECT
			id,
			name,
			address,
			company_id,
			created_at,
			updated_at
		FROM branch
		WHERE id=$1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&address,
		&phoneNumber,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Branch{
		Id:        id.String,
		Name:      name.String,
		Address:   address.String,
		CompanyID: phoneNumber.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}, nil
}

func (r *branchRepo) GetList(ctx context.Context, req *models.BranchGetListRequest) (*models.BranchGetListResponse, error) {

	var (
		resp   = &models.BranchGetListResponse{}
		query  string
		where  = " WHERE deleted_at is NULL"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			address,
			company_id,
			created_at,
			updated_at
		FROM branch
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.SearchBYName != "" {
		where += ` AND name ILIKE '%' || '` + req.SearchBYName + `' || '%'`
	}


	fmt.Println(req.SearchBYName)
	if req.SearchByAddress != "" {
		where += ` AND address ILIKE '%' || '` + req.SearchByAddress + `' || '%'`
		fmt.Println("djfffffffffffffffffffffn")
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			name        sql.NullString
			address     sql.NullString
			phoneNumber sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&address,
			&phoneNumber,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Branchs = append(resp.Branchs, &models.Branch{
			Id:        id.String,
			Name:      name.String,
			Address:   address.String,
			CompanyID: phoneNumber.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	sort.Slice(resp.Branchs[:], func(i, j int) bool {
		return resp.Branchs[i].CreatedAt > resp.Branchs[j].CreatedAt
	})

	return resp, nil
}

func (r *branchRepo) Update(ctx context.Context, req *models.UpdateBranch) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			branch
		SET
			name = :name,
			address = :address,
			company_id = :company_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":         req.Id,
		"name":       req.Name,
		"address":    req.Address,
		"company_id": req.CompanyID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *branchRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

	var (
		query string
		set   string
	)

	fmt.Println(req)

	fmt.Println()
	fmt.Println(req.Fields)
	
	

	if len(req.Fields) <= 0 {
		return 0, errors.New("no fields")
	}

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s, ", key, key)
	}

	query = `
		UPDATE
			branch
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

func (r *branchRepo) Delete(ctx context.Context, req *models.BranchPrimaryKey) error {

	_, err := r.db.Exec(ctx, "Update branch SET deleted_at=now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
