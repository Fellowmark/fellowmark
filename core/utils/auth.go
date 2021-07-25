package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nus-utils/nus-peer-review/loggers"
)

type ClaimsData struct {
	Role string      `json:"role"`
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

func DecodeBody(body io.ReadCloser, out interface{}) error {
	var unmarshalErr *json.UnmarshalFieldError
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&out)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return errors.New("Bad Request. Wrong Type provided " + unmarshalErr.Field.Name)
		} else {
			return errors.New("Bad Request. " + err.Error())
		}
	}
	return nil
}

func HandleResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func HandleResponseWithObject(w http.ResponseWriter, object interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(object)
}

func GenerateJWT(role string, object interface{}) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))

	claims := ClaimsData{
		role,
		object,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "npr-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		loggers.ErrorLogger.Println("Something Went Wrong: %s" + err.Error())
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*ClaimsData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsData{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(*ClaimsData); ok && token.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}
