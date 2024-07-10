package main

import (
	"encoding/json"
	"fmt"
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
	"gameApp/service/userservice"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	Jwt_SignKey                = "jwt_secret"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
)

func main() {
	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/health-check", healthCheckHandler)
	http.HandleFunc("/users/login", userLoginHandler)
	http.HandleFunc("/users/profile", userPeofileHandler)

	log.Println("server start on port 8088...")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func userRegisterHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		//res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, `{"error": "invalid method"}`)
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		//fmt.Fprintf(writer, `{"error": "reading request body failed"}`)
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	var uReq userservice.RegisterRequest
	err = json.Unmarshal(data, &uReq)

	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	authSvc := authservice.New(Jwt_SignKey, AccessTokenExpireDuration, RefreshTokenExpireDuration,
		AccessTokenSubject, RefreshTokenSubject)

	mysqlRepo := mysql.New()
	userSvc := userservice.NewService(authSvc, mysqlRepo)

	_, err = userSvc.Register(uReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	writer.Write([]byte(`{"message": "user created"}`))
}
func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, `{"message": "everything is good!"}`)
}

func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprintf(writer, `{"error": "invalid method"}`)
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		//fmt.Fprintf(writer, `{"error": "reading request body failed"}`)
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	var lReq userservice.LoginRequest
	err = json.Unmarshal(data, &lReq)

	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}
	authSvc := authservice.New(Jwt_SignKey, AccessTokenExpireDuration, RefreshTokenExpireDuration,
		AccessTokenSubject, RefreshTokenSubject)

	mysqlRepo := mysql.New()
	userSvc := userservice.NewService(authSvc, mysqlRepo)

	resp, err := userSvc.Login(lReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}
	data, err = json.Marshal(resp)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	writer.Write(data)
}

func userPeofileHandler(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprintf(writer, `{"error": "invalid method"}`)
	}

	authSvc := authservice.New(Jwt_SignKey, AccessTokenExpireDuration, RefreshTokenExpireDuration,
		AccessTokenSubject, RefreshTokenSubject)

	authToken := req.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(authToken)
	if err != nil {
		fmt.Fprintf(writer, `{"error": "token is not valid"}`)
	}
	mysqlRepo := mysql.New()
	userSvc := userservice.NewService(authSvc, mysqlRepo)

	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})

	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	writer.Write(data)
}
