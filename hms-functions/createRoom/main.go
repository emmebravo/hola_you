package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type RequestBody struct {
	Room string `json:"room"`
}

func createManagementToken() string {
	appAccessKey := os.Getenv("APP_ACCESS_KEY")
	appSecret := os.Getenv("APP_SECRET")

	mySigningKey := []byte(appSecret)
	expiresAt := uint32(24*3600)
	now := uint32(time.Now().UTC().Unix())
	expirationTime := now + expiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                "access_key": appAccessKey,
		"type": "management",
		"version": 2,
		"jti": uuid.New().String(),
		"iat": now,
		"exp": expirationTime,
		"nbf": now,
        })

	// sign in & get encoded token (str) after using secret
	signedToken, _ := token.SignedString(mySigningKey)
	return signedToken
}

func handleInternalServerError(errorMessage string) (*events.APIGatewayProxyResponse, error){
	err := errors.New(errorMessage)

	return &events.APIGatewayProxyResponse{
                StatusCode: http.StatusInternalServerError,
		Headers: map[string]string{"Content-Type": "application/json"},
                Body:       err.Error(),
        }, err
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error){
	var rb RequestBody
	managementToken := createManagementToken()

	b := []byte(request.Body)
	err1 := json.Unmarshal(b, &rb)
	if err1!= nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
		}, errors.New("Provide room name in the request body")
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"name":		strings.ToLower(rb.Room),
		"active":	true,
	})
	payload := bytes.NewBuffer(postBody)

	roomURL := os.Getenv("ROOM_URL")
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, roomURL, payload)

	if err != nil {
		return handleInternalServerError(err.Error())
	}

	// Auth Header
	req.Header.Add("Authorization", "Bearer " + managementToken)
	req.Header.Add("Content-Type", "application/json")

	// send HTTP req
	res, err := client.Do(req)
	if err != nil {
		return handleInternalServerError(err.Error())
	}
	defer res.Body.Close()

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return handleInternalServerError(err.Error())
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:		res.StatusCode,
		Headers:		map[string]string{"Content-Type": "application/json"},
		Body:			string(resp),
		IsBase64Encoded:	false,
	}, nil
}



func main() {
	// starts serverless fcn for API
	lambda.Start(handler)
}
