package exception

import (
	"bwahttprouter/helper"
	"database/sql"
	"fmt"
	"net/http"

	// "programmerzamannow/belajar-golang-restful-api/helper"
	// "programmerzamannow/belajar-golang-restful-api/model/web"

	"github.com/go-playground/validator/v10"
	// "gopkg.in/go-playground/validator.v9"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {

	if notFoundError(writer, request, err) {
		return
	}

	if validationErrors(writer, request, err) {
		return
	}

	if notOwnerError(writer, request, err) {
		return
	}

	// if emailNotFound(writer, request, err) {
	// 	return
	// }

	internalServerError(writer, request, err)
}

func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {

	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		response := helper.APIResponse("error in validation", http.StatusBadRequest, "error", helper.FormatValidationError(exception))

		// webResponse := web.WebResponse{
		// 	Code:   http.StatusBadRequest,
		// 	Status: "BAD REQUEST",
		// 	Data:   exception.Error(),
		// }

		helper.WriteToResponseBody(writer, response)
		return true
	} else {
		return false
	}
}

func notOwnerError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotOwnerError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		webResponse := helper.APIResponse("Not The Owner", http.StatusNotAcceptable, "error", exception.Error)
		helper.WriteToResponseBody(writer, webResponse)
		return true
	}
	return false
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok && request.URL.Path == "/api/v1/email_checkers" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		webResponse := helper.APIResponse("email available", http.StatusOK, "success", true)
		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else if ok && request.URL.Path != "/api/v1/email_checkers" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := helper.APIResponse("not found error", http.StatusNotFound, "error", exception.Error)
		helper.WriteToResponseBody(writer, webResponse)
		return true

	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := helper.APIResponse("internal server error", http.StatusInternalServerError, "error", err)

	helper.WriteToResponseBody(writer, webResponse)
}

func emailNotFound(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	fmt.Println(request.Host, "host")
	fmt.Println(request.URL.Host, "url host")
	fmt.Println(request.URL.Path, "url path")
	if err == sql.ErrNoRows && request.URL.Path == "/api/v1/checkemail" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := helper.APIResponse("not found error", http.StatusNotFound, "error", "email not registered")
		helper.WriteToResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}

}
