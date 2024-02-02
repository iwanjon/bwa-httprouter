package auth

import (
	"bwahttprouter/helper"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
// var hmacSampleSecret []byte

type Service interface {
	GenerateJWTToken(userId int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtservice struct {
}

func Newjwtservice() *jwtservice {
	return &jwtservice{}
}

var SECRET_KEY []byte = []byte("makan_malam")

func (s *jwtservice) GenerateJWTToken(userid int) (string, error) {
	// SECRET_KEY = []byte("makan_malam")
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	claim := jwt.MapClaims{
		"user_id":  userid,
		"tokenjwt": "JWT_TOKEN",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(SECRET_KEY)

	helper.PanicIfError(err, " error in create token jwt")

	fmt.Println(tokenString, err)

	return tokenString, nil
}

func (s *jwtservice) ValidateToken(tokeninput string) (*jwt.Token, error) {
	// sample token string taken from the New example
	// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokeninput, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error token input")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		fmt.Println(SECRET_KEY, "this is secret key result of validate")
		return SECRET_KEY, nil
	})
	helper.PanicIfError(err, "error jwt parser")

	fmt.Println(token, "this is token result of validate")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		fmt.Println(claims["user_id"], claims["tokenjwt"])

		return token, nil

	} else {

		// return token, errors.New("error claim converting")
		helper.PanicIfError(errors.New("error claim converting"), "error jwt parser")

	}
	return token, nil
}
