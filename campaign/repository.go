package campaign

import (
	"bwahttprouter/exception"
	"bwahttprouter/helper"
	"bwahttprouter/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Campaign, error)
	FindByUserId(ctx context.Context, userId int) ([]Campaign, error)
	FindById(ctx context.Context, campaignId int) (Campaign, error)
	SaveCampaign(ctx context.Context, campaign Campaign) (Campaign, error)
	UpdateCampaign(ctx context.Context, campaign Campaign) (Campaign, error)
	SaveImage(ctx context.Context, campaignImage CampaignImage) (CampaignImage, error)
	MarkAllImagesAsNonPrimary(ctx context.Context, campaignId int) (bool, error)
}
type repository struct {
	db *sql.DB
}

func NewCampaignRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll(ctx context.Context) ([]Campaign, error) {

	var campaigns []Campaign
	var campaign_image CampaignImage
	var campaign Campaign
	var userr user.User

	// defer r.db.Close()

	// helper.PanicIfError(err, " error create tx in FindAll repository")
	sql_campaign := "select * from campaigns"
	result_campaign, err := r.db.QueryContext(ctx, sql_campaign)
	helper.PanicIfError(err, " error create tx in select campaign repository")
	defer result_campaign.Close()

	sql_campaign_images := "select * from campaign_images where campaign_id = ?"
	statement_campaign_images, err := r.db.PrepareContext(ctx, sql_campaign_images)
	helper.PanicIfError(err, " error create statement campaign images repository")
	defer statement_campaign_images.Close()

	sql_user := "select * from users where id = ?"
	statement_campaign_user, err := r.db.PrepareContext(ctx, sql_user)
	helper.PanicIfError(err, " error create statement campaign  user  repository")
	defer statement_campaign_user.Close()

	err = errors.New("no campaign found")
	for result_campaign.Next() {
		var campaign_images []CampaignImage

		// var campaign_image CampaignImage

		err = result_campaign.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Perks,
			&campaign.BackerCount,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
			// &campaign.CampaignImages,
			// &campaign.User)
		)
		helper.PanicIfError(err, "error in scan campaigns")

		result_campaign_images, err := statement_campaign_images.QueryContext(ctx, campaign.ID)
		helper.PanicIfError(err, " error create tx in select campaign images repository")

		// defer helper.CommitOrRollback(tx)
		for result_campaign_images.Next() {

			err = result_campaign_images.Scan(
				&campaign_image.ID,
				&campaign_image.CampaignID,
				&campaign_image.FileName,
				&campaign_image.IsPrimary,
				&campaign_image.CreatedAt,
				&campaign_image.UpdatedAt,
			)
			helper.PanicIfError(err, " error in scan campaign images repository")
			campaign_images = append(campaign_images, campaign_image)
		}

		result_campaign_user := statement_campaign_user.QueryRowContext(ctx, campaign.UserID)

		err = result_campaign_user.Scan(
			&userr.ID,
			&userr.Name,
			&userr.Occupation,
			&userr.Email,
			&userr.PasswordHash,
			&userr.AvatarFileName,
			&userr.Role,
			&userr.Token,
			&userr.CreatedAt,
			&userr.UpdatedAt,
		)
		if err != nil {
			fmt.Println("no user for ", campaign.ID)
			userr = user.User{}
		}
		campaign.User = userr
		campaign.CampaignImages = campaign_images

		campaigns = append(campaigns, campaign)
		err = nil
	}
	exception.PanicIfNotFound(err, "error not found data")
	return campaigns, nil
}

func (r *repository) FindByUserId(ctx context.Context, userId int) ([]Campaign, error) {
	var campaigns []Campaign
	var campaign Campaign
	var campaign_image CampaignImage
	var userr user.User

	// defer r.db.Close()
	// tx, err := r.db.Begin()
	// helper.PanicIfError(err, "error in create transaction find campaign by user id")

	// sql_campaign := "select * from campaigns where user_id = ?"
	// result_campaign, err := tx.QueryContext(ctx, sql_campaign, userId)
	// helper.PanicIfError(err, " error in select campaign by user id repository")

	sql_campaign := "select * from campaigns where user_id = ?"
	statement_campaign, err := r.db.PrepareContext(ctx, sql_campaign)
	helper.PanicIfError(err, " error in select campaign by user id repository")

	defer statement_campaign.Close()

	sql_user := "select * from users where id = ?"
	statement_campaign_user, err := r.db.PrepareContext(ctx, sql_user)
	helper.PanicIfError(err, " error create statement campaign  user  repository")
	defer statement_campaign_user.Close()

	sql_campaign_images := "select * from campaign_images where campaign_id = ?"
	statement_campaign_images, err := r.db.PrepareContext(ctx, sql_campaign_images)
	helper.PanicIfError(err, " error create statement campaign images repository")
	defer statement_campaign_images.Close()

	// defer helper.CommitOrRollback(tx)

	result_campaign, err := statement_campaign.QueryContext(ctx, userId)
	helper.PanicIfError(err, " error in select campaign repository")

	err = errors.New("no campaign found")
	for result_campaign.Next() {
		var campaign_images []CampaignImage

		err = result_campaign.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Perks,
			&campaign.BackerCount,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.Slug,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		helper.PanicIfError(err, "error in scan repository find campaign by user id")

		result_campaign_images, err := statement_campaign_images.QueryContext(ctx, campaign.ID)
		helper.PanicIfError(err, " error create tx in select campaign images repository")

		// defer helper.CommitOrRollback(tx)
		for result_campaign_images.Next() {

			err = result_campaign_images.Scan(
				&campaign_image.ID,
				&campaign_image.CampaignID,
				&campaign_image.FileName,
				&campaign_image.IsPrimary,
				&campaign_image.CreatedAt,
				&campaign_image.UpdatedAt,
			)
			helper.PanicIfError(err, " error in scan campaign images repository")
			campaign_images = append(campaign_images, campaign_image)
		}
		// result_campaign_images.Close()
		result_campaign_user := statement_campaign_user.QueryRowContext(ctx, campaign.UserID)

		err = result_campaign_user.Scan(
			&userr.ID,
			&userr.Name,
			&userr.Occupation,
			&userr.Email,
			&userr.PasswordHash,
			&userr.AvatarFileName,
			&userr.Role,
			&userr.Token,
			&userr.CreatedAt,
			&userr.UpdatedAt,
		)
		if err != nil {
			fmt.Println("no user for ", campaign.ID)
			userr = user.User{}
		}
		campaign.User = userr
		campaign.CampaignImages = campaign_images

		campaigns = append(campaigns, campaign)
		err = nil
	}
	exception.PanicIfNotFound(err, "error not found data")
	return campaigns, nil
}

func (r *repository) FindById(ctx context.Context, campaignId int) (Campaign, error) {
	// var campaigns []Campaign
	var campaign Campaign
	var campaign_image CampaignImage
	var userr user.User

	// defer r.db.Close()

	sql_campaign := "select * from campaigns where id = ?"
	statement_campaign, err := r.db.PrepareContext(ctx, sql_campaign)
	helper.PanicIfError(err, " error in select campaign by user id repository find campaign by campaign id")
	defer statement_campaign.Close()

	sql_user := "select * from users where id = ?"
	statement_campaign_user, err := r.db.PrepareContext(ctx, sql_user)
	helper.PanicIfError(err, " error create statement campaign  user  repository find campaign by campaign id ")
	defer statement_campaign_user.Close()

	sql_campaign_images := "select * from campaign_images where campaign_id = ?"
	statement_campaign_images, err := r.db.PrepareContext(ctx, sql_campaign_images)
	helper.PanicIfError(err, " error create statement campaign images repository find campaign by campaign id")
	defer statement_campaign_images.Close()

	fmt.Println(campaignId, "rrrrrrrrrrr")
	result_campaign := statement_campaign.QueryRowContext(ctx, campaignId)
	fmt.Println(campaignId, "reeer", result_campaign, "avb,")
	// helper.PanicIfError(err, " error in select campaign repository find campaign by campaign id")

	// err = errors.New("no campaign found")
	// for result_campaign.Next() {
	var campaign_images []CampaignImage

	err = result_campaign.Scan(
		&campaign.ID,
		&campaign.UserID,
		&campaign.Name,
		&campaign.ShortDescription,
		&campaign.Description,
		&campaign.Perks,
		&campaign.BackerCount,
		&campaign.GoalAmount,
		&campaign.CurrentAmount,
		&campaign.Slug,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)
	helper.PanicIfError(err, "error in scan repository find campaign by user id find campaign by campaign id")

	result_campaign_images, err := statement_campaign_images.QueryContext(ctx, campaign.ID)
	helper.PanicIfError(err, " error create tx in select campaign images repository find campaign by campaign id")

	for result_campaign_images.Next() {

		err = result_campaign_images.Scan(
			&campaign_image.ID,
			&campaign_image.CampaignID,
			&campaign_image.FileName,
			&campaign_image.IsPrimary,
			&campaign_image.CreatedAt,
			&campaign_image.UpdatedAt,
		)
		helper.PanicIfError(err, " error in scan campaign images repository find campaign by campaign id")
		campaign_images = append(campaign_images, campaign_image)
	}

	result_campaign_user := statement_campaign_user.QueryRowContext(ctx, campaign.UserID)
	fmt.Println(campaign.ID, "campgin id ")
	err = result_campaign_user.Scan(
		&userr.ID,
		&userr.Name,
		&userr.Occupation,
		&userr.Email,
		&userr.PasswordHash,
		&userr.AvatarFileName,
		&userr.Role,
		&userr.Token,
		&userr.CreatedAt,
		&userr.UpdatedAt,
	)
	if err != nil {
		fmt.Println("no user for ", campaign.ID)
		userr = user.User{}
	}
	campaign.User = userr
	campaign.CampaignImages = campaign_images

	// campaigns = append(campaigns, campaign)
	// err = nil
	// }
	// exception.PanicIfNotFound(err, "error not found data")

	return campaign, err
}

func (r *repository) MarkAllImagesAsNonPrimary(ctx context.Context, campaignId int) (bool, error) {
	// var campaignImages []CampaignImage
	var campaign_image CampaignImage
	// defer r.db.Close()

	sql_get_images := " select * from campaign_images where campaign_id = ?"
	statement_get_campaign_images, err := r.db.PrepareContext(ctx, sql_get_images)
	helper.PanicIfError(err, "  error in prepare statement  get campaign images")
	defer func() {
		statement_get_campaign_images.Close()
		fmt.Println("madang dap")
	}()

	result_campaign_images := statement_get_campaign_images.QueryRowContext(ctx, campaignId)

	err = result_campaign_images.Scan(
		&campaign_image.ID,
		&campaign_image.CampaignID,
		&campaign_image.FileName,
		&campaign_image.IsPrimary,
		&campaign_image.CreatedAt,
		&campaign_image.UpdatedAt,
	)
	exception.PanicIfNotFound(err, "error not found data")
	// helper.PanicIfError(err, "error in execute query get campagin imagesby campaign id")

	tx, err := r.db.Begin()
	helper.PanicIfError(err, " error create transaction markallimagesasnonprimary")
	defer helper.CommitOrRollback(tx)

	sql_update_image_isaprimary := "update campaign_images set is_primary = 0 where campaign_id = ?"

	statement_update_images, err := tx.PrepareContext(ctx, sql_update_image_isaprimary)
	helper.PanicIfError(err, " error create update isprymary")
	defer statement_update_images.Close()

	result_update_images, err := statement_update_images.ExecContext(ctx, campaignId)
	helper.PanicIfError(err, "error in execute update images")
	fmt.Println(result_update_images, "xxxxx")

	// sql_get_images := " select * from campaign_images where campaign_id = ? and is_primary = 1"
	// statement_get_campaign_images, err := r.db.PrepareContext(ctx, sql_get_images)
	// helper.PanicIfError(err, "  error in prepare statement  get campaign images")
	// defer statement_get_campaign_images.Close()

	// result_campaign_images, err := statement_get_campaign_images.QueryContext(ctx, campaignId)
	// helper.PanicIfError(err, "error in execute query get campagin imagesby campaign id")

	return true, nil
}

func (r *repository) SaveCampaign(ctx context.Context, campaign Campaign) (Campaign, error) {

	// defer r.db.Close()

	tx, err := r.db.Begin()
	helper.PanicIfError(err, " error create tx in save campaign")
	defer helper.CommitOrRollback(tx)

	sql_query := "insert into campaigns( user_id, name, short_description, description, perks, backer_count, goal_amount, current_amount,slug) values(?,?,?,?,?,?,?,?,?)"
	// var campaign Campaign
	statement_inser_campaign, err := tx.PrepareContext(ctx, sql_query)
	helper.PanicIfError(err, " error in create stetement save campaign")
	defer statement_inser_campaign.Close()

	result, err := statement_inser_campaign.ExecContext(ctx, campaign.UserID, campaign.Name, campaign.ShortDescription, campaign.Description, campaign.Perks, campaign.BackerCount, campaign.GoalAmount, campaign.CurrentAmount, campaign.Slug)
	helper.PanicIfError(err, "error in exceute save campaign")
	idint64, err := result.LastInsertId()
	helper.PanicIfError(err, "error in get the last insert id")
	campaign.ID = int(idint64)
	return campaign, nil
}

func (r *repository) UpdateCampaign(ctx context.Context, campaign Campaign) (Campaign, error) {
	// var campaign Campaign
	// defer r.db.Close()
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error in create tx in update campaign reposiotory")
	defer helper.CommitOrRollback(tx)
	sql_statement := "update campaigns set user_id =? , name = ? , short_description= ? , description= ? , perks= ? , backer_count= ? , goal_amount= ? , current_amount= ? ,slug= ? where id = ?"
	staement_update, err := tx.PrepareContext(ctx, sql_statement)
	helper.PanicIfError(err, "error create stetement update campaign")
	defer staement_update.Close()
	_, err = staement_update.ExecContext(ctx, campaign.UserID, campaign.Name, campaign.ShortDescription, campaign.Description, campaign.Perks, campaign.BackerCount, campaign.GoalAmount, campaign.CurrentAmount, campaign.Slug, campaign.ID)
	helper.PanicIfError(err, "error in exceute update campaign query")

	return campaign, nil
}
func (r *repository) SaveImage(ctx context.Context, campaignImage CampaignImage) (CampaignImage, error) {
	// var campaignImage CampaignImage
	// defer r.db.Close()
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error in create tx in save images")
	defer helper.CommitOrRollback(tx)

	sql_statement := "insert into campaign_images(campaign_id, file_name, is_primary) values(?,?,?)"
	stetement, err := tx.PrepareContext(ctx, sql_statement)
	helper.PanicIfError(err, "error in create statement save images")
	defer stetement.Close()

	result, err := stetement.ExecContext(ctx, campaignImage.CampaignID, campaignImage.FileName, campaignImage.IsPrimary)
	helper.PanicIfError(err, " error in execute save image campaign")
	id, err := result.LastInsertId()
	helper.PanicIfError(err, "error in getting last id campain images repository")
	campaignImage.ID = int(id)

	return campaignImage, nil
}

func (r *repository) SaveWaeImage(ctx context.Context, campaignImage CampaignImage, campaign Campaign) (CampaignImage, error) {
	// var campaignImage CampaignImage
	// defer r.db.Close()
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error in create tx in save images")
	defer helper.CommitOrRollback(tx)

	sql_statement := "insert into campaign_images(campaign_id, file_name, is_primary) values(?,?,?)"
	stetement, err := tx.PrepareContext(ctx, sql_statement)
	helper.PanicIfError(err, "error in create statement save images")
	defer stetement.Close()

	result, err := stetement.ExecContext(ctx, campaignImage.CampaignID, campaignImage.FileName, campaignImage.IsPrimary)
	helper.PanicIfError(err, " error in execute save image campaign")
	id, err := result.LastInsertId()
	helper.PanicIfError(err, "error in getting last id campain images repository")
	campaignImage.ID = int(id)

	sql_query := "insert into campaigns( user_id, name, short_description, description, perks, backer_count, goal_amount, current_amount,slug) values(?,?,?,?,?,?,?,?,?)"
	// var campaign Campaign
	statement_inser_campaign, err := tx.PrepareContext(ctx, sql_query)
	helper.PanicIfError(err, " error in create stetement save campaign")
	defer statement_inser_campaign.Close()

	_, err = statement_inser_campaign.ExecContext(ctx, campaign.UserID, campaign.Name, campaign.ShortDescription, campaign.Description, campaign.Perks, campaign.BackerCount, campaign.GoalAmount, campaign.CurrentAmount, campaign.Slug)
	helper.PanicIfError(err, "error in exceute save campaign")

	return campaignImage, nil
}

// func (r *repository) MarkAllImagesAsNonPrimary(campaignid int) (bool, error) {
// 	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignid).Update("is_primary", false).Error
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (r *repository) Update(campaign Campaign) (Campaign, error) {
// 	err := r.db.Save(&campaign).Error

// 	if err != nil {
// 		return campaign, err
// 	}
// 	return campaign, nil

// }

// func (r *repository) Save(campaign Campaign) (Campaign, error) {
// 	err := r.db.Create(&campaign).Error
// 	if err != nil {
// 		return campaign, err
// 	}
// 	return campaign, nil
// }

// func (r *repository) SaveImage(campaignImages CampaignImage) (CampaignImage, error) {

// 	err := r.db.Create(&campaignImages).Error
// 	if err != nil {
// 		return campaignImages, err
// 	}
// 	return campaignImages, nil
// }

// func (r *repository) FindAll() ([]Campaign, error) {
// 	var campaigns []Campaign
// 	// err := r.db.Preload("CampaignImages").Find(&campaigns).Error
// 	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

// 	if err != nil {
// 		return campaigns, err
// 	}
// 	return campaigns, nil
// }

// func (r *repository) FindByUserId(userId int) ([]Campaign, error) {
// 	var campaigns []Campaign
// 	err := r.db.Where("user_id = ?", userId).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
// 	if err != nil {
// 		return campaigns, err
// 	}
// 	return campaigns, nil

// }

// func (r *repository) FindById(id int) (Campaign, error) {
// 	var campaign Campaign

// 	err := r.db.Where("id = ?", id).Preload("CampaignImages").Preload("User").Find(&campaign).Error

// 	if err != nil {
// 		return campaign, err
// 	}
// 	return campaign, nil
// }

// type Repository interface {
// 	FindAll() ([]Campaign, error)
// 	FindByUserId(userId int) ([]Campaign, error)
// 	FindById(id int) (Campaign, error)
// 	Save(campaign Campaign) (Campaign, error)
// 	Update(campaign Campaign) (Campaign, error)
// 	SaveImage(campaignImages CampaignImage) (CampaignImage, error)
// 	MarkAllImagesAsNonPrimary(campaignid int) (bool, error)
// }

// type repository struct {
// 	db *gorm.DB
// }

// func NewRepository(db *gorm.DB) *repository {
// 	return &repository{db}
// }
