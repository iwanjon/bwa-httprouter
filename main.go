package main

import (
	"bwahttprouter/app"
	"bwahttprouter/auth"
	"bwahttprouter/campaign"
	"bwahttprouter/exception"
	"bwahttprouter/handler"
	"bwahttprouter/helper"
	"bwahttprouter/middleware"
	"bwahttprouter/payment"
	"bwahttprouter/transaction"
	"bwahttprouter/user"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func main() {

	db := app.NewDB()
	validate := validator.New()
	defer func() {
		fmt.Println("ready to close")
		db.Close()
		fmt.Println("ready to close 2")
	}()
	fmt.Println(db)
	directory := http.Dir("images/")
	paymentService := payment.NewPaymentService()
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, validate)
	authUser := auth.Newjwtservice()
	userHandler := handler.NewUserHandler(userService, authUser)

	// request := context.Background()

	campaignRepository := campaign.NewCampaignRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	transactionRepository := transaction.NewRepository(db)
	transService := transaction.NewServiceTransaction(transactionRepository, campaignRepository, paymentService)
	transHandler := handler.NewTransactionHandler(transService)

	// tras, err := transactionRepository.GetByCampaignID(request, 4)
	// fmt.Println(tras, "mmmm", err)

	// input := campaign.CreateCampaignInput{
	// 	Name:             "momom",
	// 	ShortDescription: "madang short",
	// 	Description:      "madang long",
	// 	GoalAmount:       123000,
	// 	Perks:            "madsng",
	// 	User:             user.User{ID: 43},
	// }
	// getdeatail := campaign.GetCampaignDetailInput{
	// 	ID: 3,
	// }

	// campaignd := campaign.Campaign{
	// 	ID:               8,
	// 	UserID:           43,
	// 	Name:             "gogog",
	// 	ShortDescription: "nothiww ",
	// 	Description:      "just blawwwwwnk",
	// 	Perks:            "muwwwwmo",
	// 	BackerCount:      0,
	// 	GoalAmount:       100000000,
	// 	CurrentAmount:    0,
	// 	Slug:             "gugogo",
	// }

	// getimageInput := campaign.CreateCampaignImageInput{
	// 	CampaignID: 4,
	// 	IsPrimary:  true,
	// 	User:       user.User{ID: 42},
	// }
	// campaignn, err := campaignService.CreateCampaign(request, input)
	// fmt.Println(campaignn, "ddddddddddddddddd", err)

	// campaingss, err := campaignService.FindCampaigns(request, 43)
	// fmt.Println(campaingss, "ddddddddddddddddd", err)

	// detailcam, err := campaignService.GetCampaignById(request, getdeatail)
	// fmt.Println(detailcam, "mmmmm", err)

	// image, err := campaignService.SaveCampaignImage(request, getimageInput, "mmmmm")
	// fmt.Println(image, "mddddddddddddmmmm", err)

	// updated, err := campaignRepository.UpdateCampaign(request, campaignd)
	// fmt.Println(updated, "mddddddddddddmmmm", err)

	// aa, e := nn.FindAll(request)
	// eee, ee := nn.MarkAllImagesAsNonPrimary(request, 32)
	// fmt.Println(eee, ee)

	// aa, e := nn.FindById(request, 7)
	// fmt.Println(aa, "aaaaaaaa", e, "dfgfgfgf")

	// for a, i := range aa {
	// 	fmt.Println("aaa", a, "ererere", i)
	// 	fmt.Println("aaa", "ererere")
	// 	fmt.Println("aaa", "ererere")
	// 	fmt.Println("aaa", "ererere")

	// }

	// campaign_imaeg := campaign.CampaignImage{
	// 	CampaignID: 43,
	// 	FileName:   "madangsik.com",
	// 	IsPrimary:  1,
	// }

	// campaign := campaign.Campaign{
	// 	ID:               8,
	// 	UserID:           43,
	// 	Name:             "gogogwwwwwwwwwwog",
	// 	ShortDescription: "nothiwwwwwng ",
	// 	Description:      "just blawwwwwnk",
	// 	Perks:            "muwwwwmo",
	// 	BackerCount:      0,
	// 	GoalAmount:       100000000,
	// 	CurrentAmount:    0,
	// 	Slug:             "gugogogwwwwwo-mowwmo",
	// }

	// qq, ee := nn.SaveImage(request, campaign_imaeg)
	// qq, ee := campaignRepository.SaveWaeImage(request, campaign_imaeg, campaign)
	// fmt.Println("kokok", qq, "lll", ee, campaign)
	// token := auth.Newjwtservice()
	// to, errr := token.GenerateJWTToken(43)
	// fmt.Println(to, "rrrr", errr)
	// rr, errr := token.ValidateToken(to)
	// fmt.Println(rr, "rrrr", errr)

	router := httprouter.New()
	// router.ServeFiles("/files/*filepath", directory)
	// router.POST("/api/v1/user", userHandler.RegisterUser)
	// router.POST("/api/v1/login", userHandler.LoginUser)
	// // router.POST("/api/v1/checkemail", middleware.EmailChecker(userHandler.CheckEmail))
	// router.POST("/api/v1/checkemail", userHandler.CheckEmail)
	// router.POST("/api/v1/uploadavatar", middleware.AuthChecker(authUser, userService, userHandler.UploadAvatar))
	// userHandler handler.UserHandler, directory http.Dir, jwtservice auth.Service, userservice user.Service, authmiddleware func(jwtservice auth.Service, userservice user.Service, h httprouter.Handle

	userRouter := app.NewUserRouter(router, userHandler, directory, authUser, userService, middleware.AuthChecker)
	// router.POST("/api/v1/createcampaign", middleware.AuthChecker(authUser, userService, campaignHandler.CreateCampaign))
	campaignRouter := app.NewCampaignHandler(userRouter, campaignHandler, authUser, userService, middleware.AuthChecker)
	transRouter := app.NewTransactionHandler(campaignRouter, transHandler, authUser, userService, middleware.AuthChecker)

	router.PanicHandler = exception.ErrorHandler
	server := http.Server{
		Addr:    "localhost:3000",
		Handler: transRouter,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err, "error main")
}
