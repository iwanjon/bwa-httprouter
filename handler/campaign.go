package handler

import (
	"bwahttprouter/campaign"
	"bwahttprouter/exception"
	"bwahttprouter/helper"
	"bwahttprouter/middleware"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type campaignHandler struct {
	service campaign.Service
}
type CampaignHandler interface {
	CreateCampaign(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateCampaign(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetCampaigns(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetCampaign(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UploadCampaignImage(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

func NewCampaignHandler(s campaign.Service) *campaignHandler {
	return &campaignHandler{s}
}

func (h *campaignHandler) CreateCampaign(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	fmt.Println(request.Context(), "cocococooc", request.Context().Value(middleware.Contectkey))
	struct_context_interface := request.Context().Value(middleware.Contectkey)
	struct_context, ok := struct_context_interface.(middleware.StructContextKey)
	if !ok {
		helper.PanicIfError(errors.New("error in struct context"), "error in conver struct context in create campaign handler")
	}

	// userid := r.Context().Value(middleware.Contectkey)
	// fmt.Println(userid, "user id")
	// user_id, ok := userid.(int)

	// if !ok {
	// 	helper.PanicIfError(errors.New("error in convert user id"), "error in conver user id in upload avatar handler")
	// }

	currentUser := struct_context.CurrentUser

	var input campaign.CreateCampaignInput

	helper.ReadFromRequestBody(request, &input)
	input.User = currentUser
	newcampaign, err := h.service.CreateCampaign(request.Context(), input)
	helper.PanicIfError(err, " error in create campaign handler")
	campaignFormater := campaign.FormatCampaign(newcampaign)
	response := helper.APIResponse("success", http.StatusOK, "success create campaign", campaignFormater)
	helper.WriteToResponseBody(writer, response)
}

func (h *campaignHandler) UpdateCampaign(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	struct_context_interface := request.Context().Value(middleware.Contectkey)
	struct_context, ok := struct_context_interface.(middleware.StructContextKey)
	if !ok {
		helper.PanicIfError(errors.New("error in struct context"), "error in conver struct context in create campaign handler")
	}

	// userid := r.Context().Value(middleware.Contectkey)
	// fmt.Println(userid, "user id")
	// user_id, ok := userid.(int)

	// if !ok {
	// 	helper.PanicIfError(errors.New("error in convert user id"), "error in conver user id in upload avatar handler")
	// }

	currentUser := struct_context.CurrentUser

	var inputid campaign.GetCampaignDetailInput

	id := params.ByName("campaignid")
	if id == "" {
		exception.PanicIfNotFound(errors.New("not found error"), "error in get dynamic url")
	}
	idint, err := strconv.Atoi(params.ByName("campaignid"))
	exception.PanicIfNotFound(err, "error in get dynamic url int")
	inputid.ID = idint

	var input campaign.CreateCampaignInput
	input.User = currentUser
	helper.ReadFromRequestBody(request, &input)

	updatedCampaign, err := h.service.UpdateCampaign(request.Context(), inputid, input)
	helper.PanicIfError(err, "error in update campaign")
	campaignFormater := campaign.FormatCampaign(updatedCampaign)
	response := helper.APIResponse("success", http.StatusOK, "success create campaign", campaignFormater)
	helper.WriteToResponseBody(writer, response)

}
func (h *campaignHandler) GetCampaigns(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userid := request.URL.Query().Get("user_id")
	user_id := 0
	if userid != "" {
		userIdInt, err := strconv.Atoi(userid)
		helper.PanicIfError(err, "error in get user id Get Campaign handler")
		user_id = userIdInt
	}

	campaigns, err := h.service.FindCampaigns(request.Context(), user_id)
	helper.PanicIfError(err, "error in getin campaigns")
	campaignFormater := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("success", http.StatusOK, "success create campaign", campaignFormater)
	helper.WriteToResponseBody(writer, response)
}
func (h *campaignHandler) GetCampaign(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	campaignId := params.ByName("campaignid")
	campaignIdInt, err := strconv.Atoi(campaignId)
	exception.PanicIfNotFound(err, "no campaign id found in get campaign by id handler")
	campdinIdStruct := campaign.GetCampaignDetailInput{
		ID: campaignIdInt,
	}
	campaignById, err := h.service.GetCampaignById(request.Context(), campdinIdStruct)
	exception.PanicIfNotFound(err, "error in finding campaign by is handler")
	campaignformater := campaign.FormatCampaignDetail(campaignById)
	response := helper.APIResponse("succes", http.StatusOK, "success", campaignformater)

	helper.WriteToResponseBody(writer, response)

}
func (h *campaignHandler) UploadCampaignImage(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	struct_context_interface := request.Context().Value(middleware.Contectkey)
	struct_context, ok := struct_context_interface.(middleware.StructContextKey)
	if !ok {
		helper.PanicIfError(errors.New("error in struct context"), "error in conver struct context in upload avatar handler")
	}

	user_User := struct_context.CurrentUser

	var input campaign.CreateCampaignImageInput
	file, fileheader, err := request.FormFile("file")
	helper.PanicIfError(err, "error in upload campaign image handler")

	isPrimary := request.PostFormValue("is_primary")
	campaignId := request.PostFormValue("campaign_id")
	campaignIdInt, err := strconv.Atoi(campaignId)
	helper.PanicIfError(err, "error in convert campaign id upload campaign handler")

	isPrimaryBool, err := strconv.ParseBool(isPrimary)
	helper.PanicIfError(err, "error in convert primary boolean upload campaign handler")

	input.User = user_User
	input.CampaignID = campaignIdInt
	input.IsPrimary = isPrimaryBool

	path := fmt.Sprintf("images/camapign-%d-%s", input.CampaignID, fileheader.Filename)
	fileDestination, err := os.Create(path)
	helper.PanicIfError(err, "erro in create path for upload campaign image")

	_, err = io.Copy(fileDestination, file)
	helper.PanicIfError(err, " error in copy file to destination upload campaign Handler")
	_, err = h.service.SaveCampaignImage(request.Context(), input, fileDestination.Name())
	helper.PanicIfError(err, "error in save campaign image uploadcampaign handler")
	data := make(map[string]bool)
	data["is_uploaded"] = true

	respnse := helper.APIResponse("upload image success", http.StatusOK, "success", data)
	helper.WriteToResponseBody(writer, respnse)
}
