package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type RequestBody struct {
	UserID	string	`json:"userId"`
	RoomID	string	`json:"roomId"`
	Role	string	`json:"role"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var rb RequestBody

	b := []byte(request.Body)
	err1 := json.Unmarshal(b, &rb)
	if err1 != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
		}, errors.New("Provide userId, roomId, and role in the req.body")
	}

	appAccessKey := os.Getenv("APP_ACCESS_KEY")
	appSecret := os.Getenv("APP_SECRET")

	mySigningKey := []byte(appSecret)
	expiresAt := uint32(24 * 3600)
	now := uint32(time.Now().UTC().Unix())
	expirationTime := now + expiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"access_key":	appAccessKey,
		"type":		"app",
		"version":	2,
		"roomId":	rb.RoomID,
		"userId":	rb.UserID,
		"role":		rb.Role,
		"jti":		uuid.New().String(),
		"iat":		now,
		"exp":		expirationTime,
		"nbf":		now,
	})

	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 	http.StatusInternalServerError,
			Headers:	map[string]string{"Content-Type": "application/json"},
			Body:		"Internal Server Error",
		}, err
	}

	// return auth token to client
	return &events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            signedToken,
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
