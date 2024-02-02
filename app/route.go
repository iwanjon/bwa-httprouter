package app

import (
	"bwahttprouter/auth"
	"bwahttprouter/exception"
	"bwahttprouter/handler"
	"bwahttprouter/user"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewUserRouter(router *httprouter.Router, userHandler handler.UserHandler, directory http.Dir, jwtservice auth.Service, userservice user.Service, authmiddleware func(jwtservice auth.Service, userservice user.Service, h httprouter.Handle) httprouter.Handle) *httprouter.Router {
	// router := httprouter.New()
	router.ServeFiles("/files/*filepath", directory)
	router.POST("/api/v1/users", userHandler.RegisterUser)
	router.POST("/api/v1/sessions", userHandler.LoginUser)
	// router.POST("/api/v1/checkemail", middleware.EmailChecker(userHandler.CheckEmail))
	router.POST("/api/v1/email_checkers", userHandler.CheckEmail)
	router.POST("/api/v1/avatars", authmiddleware(jwtservice, userservice, userHandler.UploadAvatar))
	router.POST("/api/v1/users/fetch", authmiddleware(jwtservice, userservice, userHandler.FetchUser))
	router.PanicHandler = exception.ErrorHandler

	return router
}

func NewCampaignHandler(router *httprouter.Router, campaignHandler handler.CampaignHandler, jwtservice auth.Service, userservice user.Service, authmiddleware func(jwtservice auth.Service, userservice user.Service, h httprouter.Handle) httprouter.Handle) *httprouter.Router {

	router.GET("/api/v1/campaigns", campaignHandler.GetCampaigns)
	router.POST("/api/v1/campaigns", authmiddleware(jwtservice, userservice, campaignHandler.CreateCampaign))
	router.PUT("/api/v1/campaigns/:campaignid", authmiddleware(jwtservice, userservice, campaignHandler.UpdateCampaign))
	router.GET("/api/v1/campaigns/:campaignid", campaignHandler.GetCampaign)
	router.POST("/api/v1/campaign-images", authmiddleware(jwtservice, userservice, campaignHandler.UploadCampaignImage))

	return router
}

func NewTransactionHandler(router *httprouter.Router, transactionHandler handler.TransactionHandler, jwtservice auth.Service, userservice user.Service, authmiddleware func(jwtservice auth.Service, userservice user.Service, h httprouter.Handle) httprouter.Handle) *httprouter.Router {

	router.GET("/api/v1/campaigns/:campaignid/transactions", authmiddleware(jwtservice, userservice, transactionHandler.GetCampaignTransactions))
	router.GET("/api/v1/transactions", authmiddleware(jwtservice, userservice, transactionHandler.GetUserTransactions))
	router.POST("/api/v1/transactions", authmiddleware(jwtservice, userservice, transactionHandler.CreateTransaction))
	router.POST("/api/v1/transactions/notification", transactionHandler.GetNotif)
	return router
}
