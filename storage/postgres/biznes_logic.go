package postgres

import (
	"app/api/models"
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type biznesRepo struct {
	db *pgxpool.Pool
}

func NewBiznesRepo(db *pgxpool.Pool) *biznesRepo {
	return &biznesRepo{
		db: db,
	}
}

func (r *biznesRepo) GetTopStaff(ctx context.Context, req *models.TopStaffRequest) (*models.TopStaffResponse, error) {
	var (
		resp           = &models.TopStaffResponse{}
		query          string
		queryAssistant string
		queryCashier   string

		whereA  = " WHERE  s.status = 'success' and s.assistant_id is not null "
		havingA string

		whereC  = " WHERE  s.status = 'success'  "
		havingC string
	)

	//  GET TOP ASSISTANT

	queryAssistant = `
		SELECT 
			st.name,
			b.name,
			st.staff_type,
			SUM(s.price)
		FROM sale as s
		JOIN staff as st on st.id=s.assistant_id
		JOIN branch as b on b.id = s.branch_id
	`

	if req.FromDate != "" {
		whereA += fmt.Sprintf("and date(s.created_at)>='%s'", req.FromDate)
		whereC += fmt.Sprintf(" and date(s.created_at)>='%s'", req.FromDate)

	}

	if req.ToDate != "" {
		whereA += fmt.Sprintf(" and date(s.created_at)<'%s'", req.ToDate)
		whereC += fmt.Sprintf(" and date(s.created_at)<'%s'", req.ToDate)
	}

	havingA = `
		GROUP BY st.name, b.name,st.staff_type
		HAVING sum(s.price)>=(
			SELECT 
			sum(price)
			FROM sale
			where status = 'success'
			group by assistant_id
			order by sum(price) desc
			limit 1
		)
		order by sum(s.price) desc
	`

	query = queryAssistant + whereA + havingA

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			name      sql.NullString
			branch    sql.NullString
			staffType sql.NullString
			sum       sql.NullFloat64
		)

		err = rows.Scan(
			&name,
			&branch,
			&staffType,
			&sum,
		)
		if err != nil {
			return nil, err
		}

		resp.TopStaffs = append(resp.TopStaffs, &models.TopStaff{
			Name:      name.String,
			Branch:    branch.String,
			Total_Sum: sum.Float64,
			StaffType: staffType.String,
		})

	}

	//   GET TOP CASHIER

	queryCashier = `
		SELECT 
			st.name,
			b.name,
			st.staff_type,
			SUM(s.price)
		FROM sale as s
		JOIN staff as st on st.id=s.cashier_id
		JOIN branch as b on b.id = s.branch_id
	`

	havingC = `
		GROUP BY st.name, b.name,st.staff_type
		HAVING sum(s.price)>=(
			SELECT 
			sum(price)
			FROM sale
			where status = 'success'
			group by cashier_id
			order by sum(price) desc
			limit 1
		)
		order by sum(s.price) desc
	`

	query = queryCashier + whereC + havingC

	rows, err = r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			name      sql.NullString
			branch    sql.NullString
			staffType sql.NullString
			sum       sql.NullFloat64
		)

		err = rows.Scan(
			&name,
			&branch,
			&staffType,
			&sum,
		)
		if err != nil {
			return nil, err
		}

		resp.TopStaffs = append(resp.TopStaffs, &models.TopStaff{
			Name:      name.String,
			Branch:    branch.String,
			Total_Sum: sum.Float64,
			StaffType: staffType.String,
		})

	}

	return resp, nil

}

func (r *biznesRepo) GetTopBranchByDate(ctx context.Context) (*models.TopBranchReponse, error) {

	var (
		resp = &models.TopBranchReponse{}

		query string
	)
	query = `
		SELECT  
			date(s.created_at),
			b.name,
			SUM(s.price)
			FROM sale AS s 
			join branch AS b ON b.id = s.branch_id
		WHERE s.status = 'success'
		GROUP BY b.name, date(s.created_at)

		ORDER BY sum(s.price)
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			name      sql.NullString
			date      sql.NullString
			total_sum sql.NullFloat64
		)

		err = rows.Scan(
			&date,
			&name,
			&total_sum,
		)

		resp.TopBranchs = append(resp.TopBranchs, &models.TopBranch{
			Date:      date.String,
			Name:      name.String,
			Total_Sum: total_sum.Float64,
		})

	}

	return resp, nil
}

func (r *biznesRepo) Do_Staff_Transaction(ctx context.Context, req *models.Sale) error {

	trx, err := r.db.Begin(ctx)
	if err != nil {
		return nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	// Staff Transaction  ------->
	var (
		BonusForCashier   = 50000
		BonusForAssistant = 50000

		id2                      = uuid.New().String()
		queryGetData             string
		balanceQuery             string
		checkBonusAssistantQuery string
		// checkBonusCashierQuery   string
		cashierData          models.Staff
		cashier_tarif        string
		cashier_tarif_cash   float64
		cashier_tarif_card   float64
		assistantData        models.Staff
		assistant_tarif      string
		assistant_tarif_cash float64
		assistant_tarif_card float64

		transaction_type string
		description      string
		source_type      string

		id          sql.NullString
		name        sql.NullString
		staffType   sql.NullString
		balance     sql.NullFloat64
		tarifID     sql.NullString
		branchID    sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
		tarif       sql.NullString
		tarif_cash  sql.NullFloat64
		tarif_card  sql.NullFloat64
		count       sql.NullInt64
		total_price sql.NullFloat64

		countCashier        int
		totalPriceCashier   float64
		countAssistant      int
		totalPriceAssistant float64
	)

	queryGetData = `
		SELECT 
			s.id,
			s.name,
			s.staff_type,
			s.balance,
			s.tarif_id,
			s.branch_id,
			s.created_at,
			s.updated_at,
			st.type_tarif,
			st.amountforcash,
			st.amountforcard
		FROM staff as s 
		JOIN staff_tarif AS st ON st.id = s.tarif_id
		WHERE s.id = $1
	`

	balanceQuery = `
		UPDATE
			staff
		SET
			name=$2,
			staff_type=$3,
			balance=$4,
			tarif_id=$5,
			branch_id=$6,
			updated_at=NOW()
		WHERE id=$1
	`

	err = trx.QueryRow(ctx, queryGetData, req.CashierID).Scan(
		&id,
		&name,
		&staffType,
		&balance,
		&tarifID,
		&branchID,
		&createdAt,
		&updatedAt,
		&tarif,
		&tarif_cash,
		&tarif_card,
	)
	if err != nil {
		return err
	}

	cashierData.Id = id.String
	cashierData.Name = name.String
	cashierData.StaffType = staffType.String
	cashierData.Balance = balance.Float64
	cashierData.TarifID = tarifID.String
	cashierData.BranchID = branchID.String
	cashierData.CreatedAt = createdAt.String
	cashierData.UpdatedAt = updatedAt.String
	cashier_tarif = tarif.String
	cashier_tarif_cash = tarif_cash.Float64
	cashier_tarif_card = tarif_card.Float64


	if req.Status == "success" {
		description = "Sales Finished Succesfully"
		transaction_type = "Topup"
		source_type = "sales"
		if req.PaymentType == "Cash" {
			if cashier_tarif == "fixed" {
				cashierData.Balance += cashier_tarif_cash
			} else {
				cashierData.Balance += (req.Price * cashier_tarif_cash / 100)
			}
		} else {
			if cashier_tarif == "fixed" {
				cashierData.Balance += cashier_tarif_card
			} else {
				cashierData.Balance += (req.Price * cashier_tarif_card / 100)
			}
		}
	} else {
		description = "Sales Cancelled"
		transaction_type = "Withdraw"
		source_type = "sales"
		if req.PaymentType == "Cash" {
			if cashier_tarif == "fixed" {
				cashierData.Balance -= cashier_tarif_cash
			} else {
				cashierData.Balance -= (req.Price * cashier_tarif_cash / 100)
			}
		} else {
			if cashier_tarif == "fixed" {
				cashierData.Balance -= cashier_tarif_card
			} else {
				cashierData.Balance -= (req.Price * cashier_tarif_card / 100)
			}
		}
	}

	_, err = trx.Exec(ctx, balanceQuery,
		req.CashierID,
		cashierData.Name,
		cashierData.StaffType,
		cashierData.Balance,
		cashierData.TarifID,
		cashierData.BranchID,
	)

	if err != nil {
		return err
	}

	if req.AssistentID != "" {
		err = trx.QueryRow(ctx, queryGetData, req.AssistentID).Scan(
			&id,
			&name,
			&staffType,
			&balance,
			&tarifID,
			&branchID,
			&createdAt,
			&updatedAt,
			&tarif,
			&tarif_cash,
			&tarif_card,
		)
		if err != nil {
			return nil
		}

		assistantData.Id = id.String
		assistantData.Name = name.String
		assistantData.StaffType = staffType.String
		assistantData.Balance = balance.Float64
		assistantData.TarifID = tarifID.String
		assistantData.BranchID = branchID.String
		assistantData.CreatedAt = createdAt.String
		assistantData.UpdatedAt = updatedAt.String
		assistant_tarif = tarif.String
		assistant_tarif_cash = tarif_cash.Float64
		assistant_tarif_card = tarif_card.Float64

		switch req.Status {
		case "success":
			description = "Sales Finished Succesfully"
			transaction_type = "Topup"
			source_type = "sales"
			if req.PaymentType == "Cash" {
				if assistant_tarif == "fixed" {
					assistantData.Balance += assistant_tarif_cash
				} else {
					assistantData.Balance += (req.Price * assistant_tarif_cash / 100)
				}
			} else {
				if assistant_tarif == "fixed" {
					assistantData.Balance += assistant_tarif_card
				} else {
					assistantData.Balance += (req.Price * assistant_tarif_card / 100)
				}
			}
		case "cancel":
			description = "Sales Cancelled"
			transaction_type = "Withdraw"
			source_type = "sales"
			if req.PaymentType == "Cash" {
				if assistant_tarif == "fixed" {
					assistantData.Balance -= assistant_tarif_cash
				} else {
					assistantData.Balance -= (req.Price * assistant_tarif_cash / 100)
				}
			} else {
				if assistant_tarif == "fixed" {
					assistantData.Balance -= assistant_tarif_card
				} else {
					assistantData.Balance -= (req.Price * assistant_tarif_card / 100)
				}
			}
		}

		// Balance updated

		_, err = trx.Exec(ctx, balanceQuery,
			req.AssistentID,
			assistantData.Name,
			assistantData.StaffType,
			assistantData.Balance,
			assistantData.TarifID,
			assistantData.BranchID,
		)

		if err != nil {
			return err
		}

	}

	query2 := `
		INSERT INTO staff_transaction(id, sale_id, transaction_type, price, description, source_type, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err = trx.Exec(ctx, query2,
		id2,
		req.Id,
		transaction_type,
		req.Price,
		description,
		source_type,
	)

	if err != nil {
		return err
	}

	////////////*************CHECK    ASSISTANT    BONUS*****************/////////////////
	if req.AssistentID != "" {
		checkBonusAssistantQuery = `
				SELECT 
					Count(*),
					SUM(st.price)

				FROM staff_transaction AS st
				JOIN sale AS sl ON sl.id = st.sale_id
				WHERE st.transaction_type = 'Topup' and st.source_type = 'sales' and sl.assistant_id = $1
			`

		err = trx.QueryRow(ctx, checkBonusAssistantQuery, req.AssistentID).Scan(
			&count,
			&total_price,
		)

		if err != nil {
			return err
		}

		countAssistant = int(count.Int64)
		totalPriceAssistant = total_price.Float64

		
		checkBonusQuery := `
		SELECT 
			COUNT(*)
		FROM staff_transaction AS st
		JOIN sale AS sl ON sl.id = st.sale_id
		JOIN staff AS s ON s.id = sl.assistant_id
		WHERE st.source_type = 'bonus' and date(st.created_at)=current_date
	`
		var c int

		err = r.db.QueryRow(ctx, checkBonusQuery).Scan(
			&c,
		)

		if countAssistant >= 10 && totalPriceAssistant >= 1500000 && c < 1 {

			id4 := uuid.New().String()
			query2 := `
					INSERT INTO staff_transaction(id, sale_id, transaction_type, price, description, source_type, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, NOW())
				`
			transaction_type = "Topup"
			description = "bonus added"
			source_type = "bonus"

			_, err = trx.Exec(ctx, query2,
				id4,
				req.Id,
				transaction_type,
				req.Price,
				description,
				source_type,
			)

			if err != nil {
				return err
			}

			assistantData.Balance += float64(BonusForAssistant)

			_, err = trx.Exec(ctx, balanceQuery,
				req.AssistentID,
				assistantData.Name,
				assistantData.StaffType,
				assistantData.Balance,
				assistantData.TarifID,
				assistantData.BranchID,
			)

			if err != nil {
				return err
			}

		}
	}

	////////////*************CHECK    CASHIER    BONUS*****************/////////////////

	checkBonusCashierQuery := `
			SELECT 
				Count(*),
				SUM(st.price)

			FROM staff_transaction AS st
			JOIN sale AS sl ON sl.id = st.sale_id
			WHERE st.transaction_type = 'Topup' and st.source_type = 'sales' and sl.cashier_id = $1
		`

	err = trx.QueryRow(ctx, checkBonusCashierQuery, req.CashierID).Scan(
		&count,
		&total_price,
	)

	if err != nil {
		return err
	}

	countCashier = int(count.Int64)
	totalPriceCashier = total_price.Float64

	
	checkBonusQuery := `
		SELECT 
			COUNT(*)
		FROM staff_transaction AS st
		JOIN sale AS sl ON sl.id = st.sale_id
		JOIN staff AS s ON s.id = sl.cashier_id
		WHERE st.source_type = 'bonus' and date(st.created_at)=current_date
	`
	var c int

	err = r.db.QueryRow(ctx, checkBonusQuery).Scan(
		&c,
	)

	if countCashier >= 10 && totalPriceCashier >= 1500000 && c < 1 {

		id3 := uuid.New().String()
		query2 := `
				INSERT INTO staff_transaction(id, sale_id, transaction_type, price, description, source_type, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW())
			`
		transaction_type = "Topup"
		description = "bonus added"
		source_type = "bonus"

		_, err = trx.Exec(ctx, query2,
			id3,
			req.Id,
			transaction_type,
			req.Price,
			description,
			source_type,
		)

		if err != nil {
			return err
		}

		cashierData.Balance += float64(BonusForCashier)

		_, err = trx.Exec(ctx, balanceQuery,
			req.CashierID,
			cashierData.Name,
			cashierData.StaffType,
			cashierData.Balance,
			cashierData.TarifID,
			cashierData.BranchID,
		)

		if err != nil {
			return err
		}

	}

	return nil
}
