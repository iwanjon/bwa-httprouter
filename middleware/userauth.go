package middleware

import (
	"bwahttprouter/auth"
	"bwahttprouter/helper"
	"bwahttprouter/user"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

// type authmidlleware struct {
// 	jwtservice auth.Service
// 	handler    httprouter.Handle
// }

// func NewAuthMiddleware(j auth.Service, h httprouter.Handle) *authmidlleware {
// 	return &authmidlleware{j, h}
// }
type ContextKey string
type StructContextKey struct {
	User_id     int
	CurrentUser user.User
}

var Contectkey ContextKey = "user_id"

func AuthChecker(jwtservice auth.Service, userservice user.Service, h httprouter.Handle) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		autho := request.Header.Get("Authorization")
		fmt.Println(autho, "rtrtr")
		if !strings.Contains(autho, "Bearer") {
			helper.PanicIfError(errors.New("error in authorization"), "error in bearer contains checker")
		}
		arrayToken := strings.Split(autho, " ")
		fmt.Println(arrayToken, "array token")
		if len(arrayToken) != 2 {
			helper.PanicIfError(errors.New("error in validate authorization"), "eroor authorization len split")
		}
		jwtToken := arrayToken[1]

		tok, err := jwtservice.ValidateToken(jwtToken)
		helper.PanicIfError(err, "error invalidate tok")

		claim, ok := tok.Claims.(jwt.MapClaims)
		fmt.Println(claim, "mafang", ok)
		if !ok || !tok.Valid {
			helper.PanicIfError(errors.New("error in validation token"), "error in if validation")
		}
		user_id := claim["user_id"]
		fmt.Println(user_id, "all id")
		useridfloat, ok := user_id.(float64)
		userid := int(useridfloat)
		if !ok {
			helper.PanicIfError(errors.New("error in conver to int auth checker"), "error conver user id to int auth checker")
		}
		currentUser, err := userservice.GetUserById(request.Context(), userid)
		context_value := StructContextKey{
			User_id:     userid,
			CurrentUser: currentUser,
		}
		helper.PanicIfError(err, "error current user")

		ctx := context.WithValue(request.Context(), Contectkey, context_value)
		fmt.Println(ctx, "ctx auth id")

		fmt.Println(request.WithContext(ctx).Context().Value(Contectkey), "momomo")

		// claim, ok := token.Claims.(jwt.MapCaims)
		// fmt.Println("proses clam", claim)
		// if !ok || !token.Valid {
		// 	response := helper.APIResponse("Unathorized", http.Statusnauthorized, "unauthorized", nil)
		// 	c.AborWithStatusJSON(http.StatusUnauthorized, response)
		// 	eturn
		// }
		// userId := int(claim["user_id"].(float64))

		// user, err := usrService.GetUserById(userId)
		// if err != nil {
		// 	response := helper.APIResponse("Unathorized", http.StatuUnauthorized, "unauthorized", nil)
		// 	c.AborWithStatusJSON(http.StatusUnauthorized, response)
		// 	eturn
		// }

		h(writer, request.WithContext(ctx), params)

	}

}
