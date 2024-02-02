package transaction

import (
	"bwahttprouter/campaign"
	"bwahttprouter/exception"
	"bwahttprouter/helper"
	"bwahttprouter/user"
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	GetByCampaignID(ctx context.Context, campaignID int) ([]Transaction, error)
	GetByUserId(ctx context.Context, userId int) ([]Transaction, error)
	GetByTransactionId(ctx context.Context, Id int) (Transaction, error)
	SaveTransaction(ctx context.Context, trans Transaction) (Transaction, error)
	Update(ctx context.Context, trans Transaction) (Transaction, error)
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Update(ctx context.Context, trans Transaction) (Transaction, error) {
	tx, err := r.db.Begin()
	helper.PanicIfError(err, " error cretae tx in transactio n repo")
	defer helper.CommitOrRollback(tx)

	sql_string := "update transactions set  campaign_id =?, user_id      =?, amount    =?, status      =?, code      =?, payment_url  =? where id =? "
	stmt, err := tx.PrepareContext(ctx, sql_string)
	helper.PanicIfError(err, " error in create stmt update transaction repo")
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, trans.CampaignID, trans.UserID, trans.Amount, trans.Status, trans.Code, trans.PaymentURL, trans.ID)
	helper.PanicIfError(err, "error in exceute update transaction update repo")

	return trans, nil

}

func (r *repository) SaveTransaction(ctx context.Context, trans Transaction) (Transaction, error) {
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error ctreate tx in save transaction repo")
	defer helper.CommitOrRollback(tx)
	sql_insert := "insert into transactions(campaign_id, user_id , amount , status , code , payment_url) values( ?,?,?,?,?,?);"
	stmt, err := tx.PrepareContext(ctx, sql_insert)
	helper.PanicIfError(err, "error in create smt in save transaction repo")
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, trans.CampaignID, trans.UserID, trans.Amount, trans.Status, trans.Code, trans.PaymentURL)
	helper.PanicIfError(err, "error in execute insert in transaction repo")
	// err := r.db.Create(&trans).Error
	// // err := r.db.Omit("PaymentURL").Create(&trans).Error
	// if err != nil {
	// 	return trans, err
	// }
	idint64, err := result.LastInsertId()
	helper.PanicIfError(err, "error in getting last insert id transaction repo")
	trans.ID = int(idint64)
	return trans, nil
}

func (r *repository) GetByTransactionId(ctx context.Context, Id int) (Transaction, error) {
	var transaction Transaction
	select_stmt := "select * from transactions where id = ?"
	stmt, err := r.db.PrepareContext(ctx, select_stmt)
	helper.PanicIfError(err, "error in crreate stetement in get by id trnsaction repo")
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, Id)

	err = result.Scan(
		&transaction.ID,
		&transaction.CampaignID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Code,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.PaymentURL,
	)
	exception.PanicIfNotFound(err, "error in scan transaction by id repo")

	// err := r.db.Where("id = ?", Id).Order("id desc").Find(&transaction).Error

	// // err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	// if err != nil {
	// 	fmt.Println(err, "error repo get bu transid")
	// 	return transaction, err
	// }
	return transaction, nil
}

func (r *repository) GetByUserId(ctx context.Context, userId int) ([]Transaction, error) {
	var transactions []Transaction
	var transaction Transaction
	var transUser user.User
	var transCampaign campaign.Campaign
	var transCampaignImage campaign.CampaignImage

	sql_select_trans := `select
	transactions.id, transactions.campaign_id,  transactions.user_id,transactions.amount, transactions.status, transactions.code, transactions.payment_url
	from transactions
	where transactions.user_id =?`
	statement_1, err := r.db.PrepareContext(ctx, sql_select_trans)
	helper.PanicIfError(err, " err in create statement in get trans by user id repo")
	defer statement_1.Close()

	sql_2 := `select users.id, users.name, users.occupation, users.email, users.avatar_file_name, users.role 
		from users 
		where users.id = ?`
	stetement_2, err := r.db.PrepareContext(ctx, sql_2)
	helper.PanicIfError(err, "erro in statement 2")
	defer stetement_2.Close()

	sql_3 := `select campaigns.id, campaigns.user_id, campaigns.name, campaigns.short_description, campaigns.description , campaigns.perks, campaigns.backer_count,campaigns.goal_amount, campaigns.current_amount,campaigns.slug
		from campaigns 
		where campaigns.id = ?`
	stetement_3, err := r.db.PrepareContext(ctx, sql_3)
	helper.PanicIfError(err, "erro in statement 3")
	defer stetement_3.Close()

	sql_4 := "select file_name from campaign_images where campaign_id = ? and is_primary =1"
	stetement_4, err := r.db.PrepareContext(ctx, sql_4)
	helper.PanicIfError(err, " erro in create stement 4 repository transaction")
	defer stetement_4.Close()

	result_1, err := statement_1.QueryContext(ctx, userId)
	helper.PanicIfError(err, " err in execute create statement in get trans by user id repo")

	for result_1.Next() {
		// fmt.Println(result_1, "thi result")
		result_1.Scan(
			&transaction.ID,
			&transaction.CampaignID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Code,
			&transaction.PaymentURL,
		)
		fmt.Println("treee", transaction)

		result_2 := stetement_2.QueryRowContext(ctx, transaction.UserID)
		err = result_2.Scan(
			&transUser.ID,
			&transUser.Name,
			&transUser.Occupation,
			&transUser.Email,
			&transUser.AvatarFileName,
			&transUser.Role,
		)
		if err != nil {
			transUser = user.User{}
		}
		// helper.PanicIfError(err, "error in scan result_2")

		// fmt.Println("treee", user, transaction.CampaignID)
		result_3 := stetement_3.QueryRowContext(ctx, transaction.CampaignID)
		err = result_3.Scan(
			&transCampaign.ID,
			&transCampaign.UserID,
			&transCampaign.Name,
			&transCampaign.ShortDescription,
			&transCampaign.Description,
			&transCampaign.Perks,
			&transCampaign.BackerCount,
			&transCampaign.GoalAmount,
			&transCampaign.CurrentAmount,
			&transCampaign.Slug,
		)

		if err != nil {
			transCampaign = campaign.Campaign{}
		}

		result_4, err := stetement_4.QueryContext(ctx, transCampaign.ID)
		helper.PanicIfError(err, fmt.Sprintf(" errrr in result 4 %d transaction repository", transaction.CampaignID))
		var transCampaignImages []campaign.CampaignImage
		for result_4.Next() {
			result_4.Scan(
				&transCampaignImage.FileName,
			)
			transCampaignImages = append(transCampaignImages, transCampaignImage)
		}
		transCampaign.CampaignImages = transCampaignImages
		// helper.PanicIfError(err, "error inscan result_3")

		fmt.Println("rrrrrrr")
		fmt.Println("rrrrrrr")
		transaction.User = transUser
		transaction.Campaign = transCampaign
		fmt.Println("treee", transaction)
		transactions = append(transactions, transaction)

	}
	// err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userId).Order("id desc").Find(&transactions).Error
	// // err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	// if err != nil {
	// 	fmt.Println(err, "error repo get bu userid")
	// 	return transactions, err
	// }
	return transactions, nil
}

func (r *repository) GetByCampaignID(ctx context.Context, campaignID int) ([]Transaction, error) {
	var transactions []Transaction
	var transaction Transaction
	var transCampaign campaign.Campaign
	var transUser user.User

	sql_select_trans := `select
	transactions.id, transactions.campaign_id,  transactions.user_id,transactions.amount, transactions.status, transactions.code, transactions.payment_url
	from transactions
	where transactions.campaign_id =?`
	statement_1, err := r.db.PrepareContext(ctx, sql_select_trans)
	helper.PanicIfError(err, " err in create statement in get by campaign id repoo")
	defer statement_1.Close()

	sql_2 := `select users.id, users.name, users.occupation, users.email, users.avatar_file_name, users.role 
		from users 
		where users.id = ?`
	stetement_2, err := r.db.PrepareContext(ctx, sql_2)
	helper.PanicIfError(err, "erro in statement 2 get by campaign id repo")
	defer stetement_2.Close()

	sql_3 := `select campaigns.id, campaigns.user_id, campaigns.name, campaigns.short_description, campaigns.description , campaigns.perks, campaigns.backer_count,campaigns.goal_amount, campaigns.current_amount,campaigns.slug
		from campaigns 
		where campaigns.id = ?`
	stetement_3, err := r.db.PrepareContext(ctx, sql_3)
	helper.PanicIfError(err, "erro in statement 3 get by campaign id repo")
	defer stetement_3.Close()

	result_1, err := statement_1.QueryContext(ctx, campaignID)
	helper.PanicIfError(err, " err in execute create statement inget by campaign id repoo")

	for result_1.Next() {

		result_1.Scan(
			&transaction.ID,
			&transaction.CampaignID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Code,
			&transaction.PaymentURL,
		)

		result_2 := stetement_2.QueryRowContext(ctx, transaction.UserID)
		err = result_2.Scan(
			&transUser.ID,
			&transUser.Name,
			&transUser.Occupation,
			&transUser.Email,
			&transUser.AvatarFileName,
			&transUser.Role,
		)
		if err != nil {
			transUser = user.User{}
		}

		result_3 := stetement_3.QueryRowContext(ctx, transaction.CampaignID)
		err = result_3.Scan(
			&transCampaign.ID,
			&transCampaign.UserID,
			&transCampaign.Name,
			&transCampaign.ShortDescription,
			&transCampaign.Description,
			&transCampaign.Perks,
			&transCampaign.BackerCount,
			&transCampaign.GoalAmount,
			&transCampaign.CurrentAmount,
			&transCampaign.Slug,
		)
		if err != nil {
			transCampaign = campaign.Campaign{}
		}

		fmt.Println("rrrrrrr")
		fmt.Println("rrrrrrr")
		transaction.User = transUser
		transaction.Campaign = transCampaign
		fmt.Println("treee", transaction)
		transactions = append(transactions, transaction)

	}

	return transactions, nil

}
