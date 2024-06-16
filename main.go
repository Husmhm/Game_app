package main

import (
	"encoding/json"
	"fmt"
	"gameApp/repository/mysql"
	"gameApp/service/userservice"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/users/register", userRegisterHandler)
	http.HandleFunc("/health-check", healthCheckHandler)
	http.HandleFunc("/users/login", userLoginHandler)

	log.Println("server start on port 8080...")
	http.ListenAndServe(":8080", nil)
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

	mysqlRepo := mysql.New()
	userSvc := userservice.NewService(mysqlRepo)

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

	mysqlRepo := mysql.New()
	userSvc := userservice.NewService(mysqlRepo)

	_, err = userSvc.Login(lReq)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
		return
	}

	writer.Write([]byte(`{"message": "user credentials is ok"}`))
}
