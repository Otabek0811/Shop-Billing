CREATE TABLE company (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE   Staff(
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    company_id UUID REFERENCES company(id)
);

CREATE TABLE staff (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    AmountForCash   NUMERIC DEFAULT 0,
    AmountForCard   NUMERIC DEFAULT 0
);

CREATE TABLE staff(
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    staff_type VARCHAR NOT NULL,
    balance NUMERIC DEFAULT 0,
    tarif_id UUID REFERENCES staff(id),
    Staff_id UUID REFERENCES Staff(id)

);

CREATE TABLE sale (
    id UUID PRIMARY KEY,
    Staff_id UUID REFERENCES Staff(id),
    assistant_id UUID REFERENCES staff(id),
    cashier_id UUID NOT NULL REFERENCES staff(id) ,
    price NUMERIC NOT NULL,
    payment_type VARCHAR NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'success',
    client_name VARCHAR
);


CREATE TABLE staff_transaction(
    id UUID PRIMARY KEY,
    transaction_type VARCHAR DEFAULT 'Topup',
    description VARCHAR ,
    amount NUMERIC NOT NULL,
    sale_id UUID REFERENCES sale(id),
    staff_id UUID REFERENCES staff(id)
);






//   Check count>=10 or Total_Price >= 1_500_000
	// countQuery := `
	// 	SELECT
	// 		count(*),
	// 		sum(price)
	// 	FROM sale
	// 	WHERE status = 'success'
	// `
	// err = trx.QueryRow(ctx, countQuery).Scan(
	// 	&count,
	// 	&totalPrice,
	// )

	// if err != nil {
	// 	return "", nil
	// }






    // checkBonusAssistantQuery = `
		SELECT 
	 		Count(*),
	 		SUM(st.price)

	 	FROM staff_transaction AS st
	 	JOIN sale AS sl ON sl.id = st.sale_id
	 	where st.transaction_type = 'Topup' and sl.assistant_id = $1
	 `

	// err = trx.QueryRow(ctx, checkBonusAssistantQuery).Scan(
	// 	&countAssistant,
	// 	&totalPriceAssistant,
	// )

	// if err != nil {
	// 	return err
	// }

	// checkBonusCashierQuery = `
	// 	SELECT 
	// 		COUNT(*),
	// 		SUM(st.price)

	// 	FROM staff_transaction AS st
	// 	JOIN sale AS sl ON sl.id = st.sale_id
	// 	JOIN staff AS s ON s.id = sl.assistant_id
	// 	where st.transaction_type = 'Topup'
	// `

	// err = trx.QueryRow(ctx, checkBonusCashierQuery).Scan(
	// 	&countCashier,
	// 	&totalPriceCashier,
	// )

	// if err != nil {
	// 	return err
	// }

	// if countCashier >= 10 && totalPriceCashier >= 1500000 {
	// 	source_type="bonus"
	// 	transaction_type= "Topup"
	// 	description="bonus added"

	// 	query2 := `
	// 		INSERT INTO staff_transaction(id, sale_id, transaction_type, price, description, source_type, updated_at)
	// 		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	// 	`

	// 	_, err = trx.Exec(ctx, query2,
	// 		id2,
	// 		req.Id,
	// 		transaction_type,
	// 		req.Price,
	// 		description,
	// 		source_type,
	// 	)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	//  add bonus to cashier balance

	// 	queryGetData = `
	// 		SELECT 
	// 			id,
	// 			name,
	// 			staff_type,
	// 			balance,
	// 			tarif_id,
	// 			branch_id,
	// 			created_at
	// 		FROM staff
	// 	`

	// 	balanceQuery = `
	// 		UPDATE
	// 			staff
	// 		SET
	// 			name=$2,
	// 			staff_type=$3,
	// 			balance=$4,
	// 			tarif_id=$5,
	// 			branch_id=balance+$6,
	// 			updated_at=NOW()
	// 		WHERE id=$1
	// 	`

	// 	err = trx.QueryRow(ctx, queryGetData, req.CashierID).Scan(
	// 		&cashierData.Id,
	// 		&cashierData.Name,
	// 		&cashierData.StaffType,
	// 		&cashierData.Balance,
	// 		&cashierData.TarifID,
	// 		&cashierData.BranchID,
	// 		&cashierData.CreatedAt,
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	_, err = trx.Exec(ctx, balanceQuery,
	// 		req.CashierID,
	// 		cashierData.Name,
	// 		cashierData.StaffType,
	// 		BonusForCashier,
	// 		cashierData.TarifID,
	// 		cashierData.BranchID,
	// 		cashierData.CreatedAt,
	// 	)

	// 	if err != nil {
	// 		return err
	// 	}


		




	// }

	// Insert Data in Staff_transaction Table

	





	SELECT 
		COUNT(*),
		SUM(st.price)

		FROM staff_transaction AS st
	 	JOIN sale AS sl ON sl.id = st.sale_id
	 	JOIN staff AS s ON s.id = sl.assistant_id
	 	where st.source_type = 'bonus' and date(st.created_at)=current_date;



        	SELECT 
		COUNT(*),
		SUM(st.price)

		FROM staff_transaction AS st
	 	JOIN sale AS sl ON sl.id = st.sale_id
	 	JOIN staff AS s ON s.id = sl.assistant_id
	 	where st.source_type='bonus' and date(st.created_at)>current_date;









		SELECT 
			st.name,
			b.name,
			st.staff_type,
			SUM(price)
		FROM sale as s
		JOIN staff as st on st.id=s.assistant_id
		JOIN branch as b on b.id = s.branch_id
		WHERE date(s.created_at)<'2026-01-01' and date(s.created_at)>='2006-01-01' and s.assistant_id is not null 
		GROUP BY st.name, b.name,st.staff_type
		HAVING sum(s.price)>=(
			SELECT 
			sum(price)
			FROM sale
			group by assistant_id
			order by sum(price) desc
			limit 1
		)
		order by sum(s.price) desc



		SELECT  

		date(s.created_at),
		b.name,
		SUM(s.price)
		from sale as s 
		join branch as b on b.id = s.branch_id
		where s.status = 'success'
		group by b.name, date(s.created_at)

		order by sum(s.price)



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
				'amountforcash',st.amountforcash,
				'amountforcard',st.amountforcard,
				'created_at',st.created_at,
				'updated_at',st.updated_at	
			) as tarif_data,

			JSON_BUILD_OBJECT (
				'id',b.id,
				'name',b.name,
				'address',b.address,
				'company_id',b.company_id,
				'created_at',b.created_at,
				'updated_at',b.updated_at	
			) as branch_data

		FROM staff as s
		LEFT JOIN staff_tarif AS st ON st.id=s.tarif_id
		LEFT JOIN branch AS b ON b.id=s.branch_id
		
		WHERE s.id='9a3e1b06-68b4-446c-8487-cc308dda66e9' and s.deleted_at is null



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
				'amountforcash',st.amountforcash,
				'amountforcard',st.amountforcard,
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
		
		WHERE s.id='9a3e1b06-68b4-446c-8487-cc308dda66e9' and s.deleted_at is null;









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
			JSON_BUILD_OBJECT (
				'id',b.id,
				'name',b.name,
				'address',b.address,
				'company_id',b.company_id,
				'created_at',b.created_at,
				'updated_at',b.updated_at	
			) AS branch_data,
			JSON_BUILD_OBJECT (
				'id',sta.id,
				'name',sta.name,
				'staff_type',sta.staff_type,
				'balance',sta.balance,
				'tarif_id',sta.tarif_id,
				'branch_id',sta.branch_id,
				'created_at',sta.created_at,
				'updated_at',sta.updated_at
			) AS assistant_data,
			JSON_BUILD_OBJECT (
				'id',stc.id,
				'name',stc.name,
				'staff_type',stc.staff_type,
				'balance',stc.balance,
				'tarif_id',stc.tarif_id,
				'branch_id',stc.branch_id,
				'created_at',stc.created_at,
				'updated_at',stc.updated_at	
			) AS cashier_data
		FROM sale as s
		LEFT JOIN branch AS b ON b.id=s.branch_id
		LEFT JOIN staff AS sta ON sta.id = s.assistant_id
		LEFT JOIN staff AS stc ON stc.id = s.cashier_id
		WHERE s.id='a43470ce-e011-4972-818d-1d10af794f17' and s.deleted_at is NULL;