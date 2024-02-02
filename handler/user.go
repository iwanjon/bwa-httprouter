package handler

import (
	"bwahttprouter/auth"
	"bwahttprouter/helper"
	"bwahttprouter/middleware"
	"bwahttprouter/user"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type UserHandler interface {
	RegisterUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	LoginUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	CheckEmail(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	UploadAvatar(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	FetchUser(w http.ResponseWriter, r *http.Request, param httprouter.Params)
}

type userHandler struct {
	userService user.Service
	authUser    auth.Service
}

func NewUserHandler(u user.Service, a auth.Service) *userHandler {
	return &userHandler{u, a}
}

func (h *userHandler) RegisterUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var input user.RegisterUser
	helper.ReadFromRequestBody(request, &input)
	userRegistered, err := h.userService.RegisterUser(request.Context(), input)
	helper.PanicIfError(err, " erroor register user handler")
	jwtToken, err := h.authUser.GenerateJWTToken(userRegistered.ID)
	helper.PanicIfError(err, "error generate token")

	val := sql.NullString{
		String: jwtToken,
		Valid:  true,
	}
	userRegistered.Token = val
	formatedUser := user.Formatuser(userRegistered, val.String)
	response := helper.APIResponse("success register user", http.StatusOK, "success", formatedUser)
	helper.WriteToResponseBody(writer, response)

}

func (h *userHandler) LoginUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	var input user.LoginInput
	helper.ReadFromRequestBody(request, &input)
	userRegistered, err := h.userService.LoginUser(request.Context(), input)
	helper.PanicIfError(err, "error in login handler")

	jwtToken, err := h.authUser.GenerateJWTToken(userRegistered.ID)
	helper.PanicIfError(err, "error generate token")

	val := sql.NullString{
		String: jwtToken,
		Valid:  true,
	}
	userRegistered.Token = val
	formatedUser := user.Formatuser(userRegistered, val.String)
	repsonse := helper.APIResponse("success login", http.StatusAccepted, "success", formatedUser)
	helper.WriteToResponseBody(writer, repsonse)
	// writer.Write(repsonse)

}

func (h *userHandler) CheckEmail(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var input user.CheckEmailInput
	helper.ReadFromRequestBody(r, &input)
	available, err := h.userService.CheckEmailAvailable(r.Context(), input)
	helper.PanicIfError(err, "error handler check email")
	data := map[string]bool{"is_available": available}

	respnse := helper.APIResponse("email already registered", http.StatusOK, "success", data)
	helper.WriteToResponseBody(w, respnse)

}

func (h *userHandler) UploadAvatar(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	fmt.Println(r.Context(), "cocococooc", r.Context().Value(middleware.Contectkey))
	struct_context_interface := r.Context().Value(middleware.Contectkey)
	struct_context, ok := struct_context_interface.(middleware.StructContextKey)
	if !ok {
		helper.PanicIfError(errors.New("error in struct context"), "error in conver struct context in upload avatar handler")
	}

	// userid := r.Context().Value(middleware.Contectkey)
	// fmt.Println(userid, "user id")
	// user_id, ok := userid.(int)

	// if !ok {
	// 	helper.PanicIfError(errors.New("error in convert user id"), "error in conver user id in upload avatar handler")
	// }

	user_id := struct_context.User_id
	// user_id := int(user_id_f)
	// request.ParseMultipartForm(32 << 20)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("images/%d-%s", user_id, fileHeader.Filename)
	// fileDestination, err := os.Create("./images/" + strconv.Itoa(user_id) + "-" + fileHeader.Filename)
	fileDestination, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(fileDestination, file)
	if err != nil {
		panic(err)
	}
	fmt.Println(fileDestination, fileDestination.Name(), fileHeader.Filename, "mudang wae")
	// user, err := h.userService.SaveAvatar(r.Context(), user_id, fileDestination.Name())
	_, err = h.userService.SaveAvatar(r.Context(), user_id, fileDestination.Name())
	if err != nil {
		panic(err)
	}
	data := make(map[string]bool)
	data["is_uploaded"] = true
	fmt.Println(fileDestination, fileDestination.Name(), fileHeader.Filename, "mudang wae")
	respnse := helper.APIResponse("upload image success", http.StatusOK, "success", data)
	helper.WriteToResponseBody(w, respnse)
}

func (h *userHandler) FetchUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	struct_context_interface := r.Context().Value(middleware.Contectkey)
	struct_context, ok := struct_context_interface.(middleware.StructContextKey)
	if !ok {
		helper.PanicIfError(errors.New("error in struct context"), "error in conver struct context in upload avatar handler")
	}

	// userid := r.Context().Value(middleware.Contectkey)
	// fmt.Println(userid, "user id")
	// user_id, ok := userid.(int)

	// if !ok {
	// 	helper.PanicIfError(errors.New("error in convert user id"), "error in conver user id in upload avatar handler")
	// }

	currentUser := struct_context.CurrentUser

	formatter := user.Formatuser(currentUser, "")

	response := helper.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)
	helper.WriteToResponseBody(w, response)
}

// func (h *userHandler) FetchUser(c *gin.Context) {

// 	currentUser := c.MustGet("curresntuser").(user.User)

// 	formatter := user.Formatuser(currentUser, "")

// 	response := helper.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)

// 	c.JSON(http.StatusOK, response)

// }
