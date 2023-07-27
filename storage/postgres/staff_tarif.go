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

type staff_tarifRepo struct {
	db *pgxpool.Pool
}

func NewStaff_TarifRepo(db *pgxpool.Pool) *staff_tarifRepo {
	return &staff_tarifRepo{
		db: db,
	}
}

func (r *staff_tarifRepo) Create(ctx context.Context, req *models.CreateStaff_Tarif) (string, error) {
	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO staff_tarif(id, name, type_tarif, amountforcash, amountforcard, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.TypeTarif,
		req.AmountForCash,
		req.AmountForCard,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *staff_tarifRepo) GetByID(ctx context.Context, req *models.Staff_TarifPrimaryKey) (*models.Staff_Tarif, error) {

	var (
		query string

		id              sql.NullString
		name            sql.NullString
		type_tarif      sql.NullString
		amount_for_cash sql.NullFloat64
		amount_for_card sql.NullFloat64
		createdAt       sql.NullString
		updatedAt       sql.NullString
	)

	query = `
		SELECT
			id,
			name,
			type_tarif,
			amountforcash,
			amountforcard,
			created_at,
			updated_at
		FROM staff_tarif
		WHERE id=$1 and deleted_at is NULL
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&type_tarif,
		&amount_for_cash,
		&amount_for_card,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Staff_Tarif{
		Id:            id.String,
		Name:          name.String,
		TypeTarif:     type_tarif.String,
		AmountForCash: amount_for_cash.Float64,
		AmountForCard: amount_for_card.Float64,
		CreatedAt:     createdAt.String,
		UpdatedAt:     updatedAt.String,
	}, nil
}

func (r *staff_tarifRepo) GetList(ctx context.Context, req *models.Staff_TarifGetListRequest) (*models.Staff_TarifGetListResponse, error) {

	var (
		resp   = &models.Staff_TarifGetListResponse{}
		query  string
		where  = " WHERE deleted_at is NULL "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			type_tarif,
			amountforcash,
			amountforcard,
			created_at,
			updated_at
		FROM staff_tarif
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id              sql.NullString
			name            sql.NullString
			type_tarif      sql.NullString
			amount_for_cash sql.NullFloat64
			amount_for_card sql.NullFloat64
			createdAt       sql.NullString
			updatedAt       sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&type_tarif,
			&amount_for_cash,
			&amount_for_card,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Staff_Tarifs = append(resp.Staff_Tarifs, &models.Staff_Tarif{
			Id:            id.String,
			Name:          name.String,
			TypeTarif:     type_tarif.String,
			AmountForCash: amount_for_cash.Float64,
			AmountForCard: amount_for_card.Float64,
			CreatedAt:     createdAt.String,
			UpdatedAt:     updatedAt.String,
		})
	}

	sort.Slice(resp.Staff_Tarifs[:], func(i, j int) bool {
		return resp.Staff_Tarifs[i].CreatedAt > resp.Staff_Tarifs[j].CreatedAt
	})
	return resp, nil
}

func (r *staff_tarifRepo) Update(ctx context.Context, req *models.UpdateStaff_Tarif) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			staff_tarif
		SET
			name = :name,
			type_tarif = :type_tarif,
			amountforcash = :amountforcash,
			amountforcard = :amountforcard,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":              req.Id,
		"name":            req.Name,
		"type_tarif":      req.TypeTarif,
		"amount_for_cash": req.AmountForCash,
		"amount_for_card": req.AmountForCard,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *staff_tarifRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

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
			staff_tarif
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

func (r *staff_tarifRepo) Delete(ctx context.Context, req *models.Staff_TarifPrimaryKey) error {

	_, err := r.db.Exec(ctx,  "Update staff_tarif SET deleted_at=now() WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
