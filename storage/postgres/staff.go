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
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type staffRepo struct {
	db *pgxpool.Pool
}

func NewStaffRepo(db *pgxpool.Pool) *staffRepo {
	return &staffRepo{
		db: db,
	}
}

func (r *staffRepo) Create(ctx context.Context, req *models.CreateStaff) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO staff(id, name, staff_type, balance, tarif_id, branch_id, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.StaffType,
		req.Balance,
		helper.NewNullString(req.TarifID),
		helper.NewNullString(req.BranchID),
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *staffRepo) GetByID(ctx context.Context, req *models.StaffPrimaryKey) (*models.Staff, error) {

	var (
		query string

		id         sql.NullString
		name       sql.NullString
		staff_type sql.NullString
		balance    sql.NullFloat64
		tarif_id   sql.NullString
		branch_id  sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString

		tarifObj  pgtype.JSONB
		branchObj pgtype.JSON
	)

	query = `
		SELECT
			s.id,
			s.name,
			s.staff_type,
			s.balance,
			s.tarif_id,
			s.branch_id,
			s.created_at,
			s.updated_at,
			JSON_BUILD_OBJECT (
				'id',st.id,
				'name',st.name,
				'type_tarif',st.type_tarif,
				'amount_for_cash',st.amountforcash,
				'amount_for_card',st.amountforcard,
				'created_at',st.created_at,
				'updated_at',st.updated_at	
			) AS tarif_data,
			JSON_BUILD_OBJECT (
				'id',b.id,
				'name',b.name,
				'address',b.address,
				'company_id',b.company_id,
				'created_at',b.created_at,
				'updated_at',b.updated_at	
			) AS branch_data
		FROM staff as s
		LEFT JOIN staff_tarif AS st ON st.id=s.tarif_id
		LEFT JOIN branch AS b ON b.id=s.branch_id
		
		WHERE s.id=$1 and s.deleted_at is null
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&staff_type,
		&balance,
		&tarif_id,
		&branch_id,
		&createdAt,
		&updatedAt,
		&tarifObj,
		&branchObj,
	)

	if err != nil {
		return nil, err
	}

	tarif := models.Staff_Tarif{}
	tarifObj.AssignTo(&tarif)

	branch := models.Branch{}
	branchObj.AssignTo(&branch)
	return &models.Staff{
		Id:         id.String,
		Name:       name.String,
		StaffType:  staff_type.String,
		Balance:    balance.Float64,
		TarifID:    tarif_id.String,
		BranchID:   branch_id.String,
		CreatedAt:  createdAt.String,
		UpdatedAt:  updatedAt.String,
		TarifData:  &tarif,
		BranchData: &branch,
	}, nil
}

func (r *staffRepo) GetList(ctx context.Context, req *models.StaffGetListRequest) (*models.StaffGetListResponse, error) {

	var (
		resp   = &models.StaffGetListResponse{}
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
			staff_type,
			balance,
			tarif_id,
			branch_id,
			created_at,
			updated_at
		FROM staff
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.SearchByName != "" {
		where += ` AND name ILIKE '%' || '` + req.SearchByName + `' || '%'`
	}
	if req.SearchByBranchID != "" {
		where += ` AND branch_id ILIKE '%' || '` + req.SearchByBranchID + `' || '%'`
	}
	if req.SearchBYStaffType != "" {
		where += ` AND staff_type ILIKE '%' || '` + req.SearchBYStaffType + `' || '%'`
	}
	if req.SearchByTarifID != "" {
		where += ` AND tarif_id ILIKE '%' || '` + req.SearchByTarifID + `' || '%'`
	}
	if req.SearchBalanceFrom != "" {
		where += ` AND balance  >= ||` + req.SearchBalanceFrom
	}
	if req.SearchBalanceTo != "" {
		where += ` AND balance <||` + req.SearchBalanceTo
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id         sql.NullString
			name       sql.NullString
			staff_type sql.NullString
			balance    sql.NullFloat64
			tarif_id   sql.NullString
			branch_id  sql.NullString
			createdAt  sql.NullString
			updatedAt  sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&staff_type,
			&balance,
			&tarif_id,
			&branch_id,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Staffs = append(resp.Staffs, &models.Staff{
			Id:        id.String,
			Name:      name.String,
			StaffType: staff_type.String,
			Balance:   balance.Float64,
			TarifID:   tarif_id.String,
			BranchID:  branch_id.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	sort.Slice(resp.Staffs[:], func(i, j int) bool {
		return resp.Staffs[i].CreatedAt > resp.Staffs[j].CreatedAt
	})

	return resp, nil
}

func (r *staffRepo) Update(ctx context.Context, req *models.UpdateStaff) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			staff
		SET
			name = :name,
			staff_type = :staff_type,
			balance = :balance,
			tarif_id = :tarif_id,
			branch_id = :branch_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":         req.Id,
		"name":       req.Name,
		"staff_type": req.StaffType,
		"balance":    req.Balance,
		"tarif_id":   req.TarifID,
		"branch_id":  req.BranchID,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staffRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

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
			staff
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

func (r *staffRepo) Delete(ctx context.Context, req *models.StaffPrimaryKey) error {

	_, err := r.db.Exec(ctx,  "Update staff SET deleted_at=now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
