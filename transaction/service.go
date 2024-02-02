package transaction

import (
	"bwahttprouter/campaign"
	"bwahttprouter/exception"
	"bwahttprouter/helper"
	"bwahttprouter/payment"
	"context"
	"errors"
	"strconv"
)

type service struct {
	repo               Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionByCampaignID(ctx context.Context, input GetCampaignTransactionsInput) ([]Transaction, error)
	GetTransactionByUserID(ctx context.Context, UserID int) ([]Transaction, error)
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error)
	ProcessPayment(ctx context.Context, trans TransactionNotificationInput) error
}

func NewServiceTransaction(r Repository, campaignRepository campaign.Repository, payservice payment.Service) *service {
	return &service{r, campaignRepository, payservice}
}

func (s *service) ProcessPayment(ctx context.Context, input TransactionNotificationInput) error {
	trans_id, err := strconv.Atoi(input.OrderID)
	helper.PanicIfError(err, " erro in conv order id service ")
	// fmt.Println(trans_id)
	transaction, err := s.repo.GetByTransactionId(ctx, trans_id)
	helper.PanicIfError(err, " error in get transaction by id service")
	// if err != nil {
	// 	return err
	// }

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedtransaction, err := s.repo.Update(ctx, transaction)
	helper.PanicIfError(err, " error in update transaction service")
	// if err != nil {
	// 	return err
	// }

	campaign, err := s.campaignRepository.FindById(ctx, updatedtransaction.CampaignID)
	helper.PanicIfError(err, " error in get campaign by id")
	// if err != nil {
	// 	return err
	// }

	if transaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedtransaction.Amount

		_, err := s.campaignRepository.UpdateCampaign(ctx, campaign)
		helper.PanicIfError(err, " error in update campaign by id service transaction")
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func (s *service) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	var trans Transaction

	trans.Amount = input.Amount
	trans.User = input.User
	trans.CampaignID = input.CampaignID
	trans.Status = "pending"

	trans.Code = "muamama"
	trans.UserID = input.User.ID

	// fmt.Println(trans, "madang trans")
	newtrans, err := s.repo.SaveTransaction(ctx, trans)
	helper.PanicIfError(err, " error in create transaction service")
	// if err != nil {
	// 	return newtrans, err
	// }

	paymentTransaction := payment.Transaction{
		ID:     newtrans.ID,
		Amount: newtrans.Amount,
	}

	url, err := s.paymentService.GetPaymentUrl(paymentTransaction, input.User)
	// if err != nil {
	// 	fmt.Println("error url")
	// 	return newtrans, err
	// }
	helper.PanicIfError(err, " error in create payment url service")
	newtrans.PaymentURL = url
	newtranss, err := s.repo.Update(ctx, newtrans)
	helper.PanicIfError(err, " erro in update transaction url serviced")
	// if err != nil {
	// 	return newtranss, err
	// }
	return newtranss, nil
}

func (s *service) GetTransactionByUserID(ctx context.Context, UserID int) ([]Transaction, error) {
	var transactions []Transaction
	// fmt.Println("gettransbyuserid", UserID)
	transactions, err := s.repo.GetByUserId(ctx, UserID)
	helper.PanicIfError(err, " error in get transaction by user id service")
	// if err != nil {
	// 	return transactions, err
	// }
	return transactions, nil
}

func (s *service) GetTransactionByCampaignID(ctx context.Context, input GetCampaignTransactionsInput) ([]Transaction, error) {
	var transactions []Transaction
	campaign, err := s.campaignRepository.FindById(ctx, input.ID)
	helper.PanicIfError(err, " erro in get acampaign by id service")
	// if err != nil {
	// 	return transactions, err
	// }

	if campaign.UserID != input.User.ID {
		// fmt.Println("error not owner check")
		// fmt.Println("error not owner check", campaign.UserID, input.User.ID)
		// return transactions, errors.New("not an owner of campaign")
		exception.PanicIfNotOwner(errors.New("error not owner of campaign"), " error in campaign authorized ")
	}

	transactions, err = s.repo.GetByCampaignID(ctx, input.ID)
	helper.PanicIfError(err, " erro in get find transacstion by campaign id service")
	// if err != nil {
	// 	return transactions, err
	// }
	return transactions, nil
}
