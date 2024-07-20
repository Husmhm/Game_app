package main

import (
	"gameApp/config"
	"gameApp/delivery/httpsever"
	"gameApp/repository/mysql"
	"gameApp/service/authservice"
	"gameApp/service/userservice"
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
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth: authservice.Config{
			SignKey:               Jwt_SignKey,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessTokenSubject:    AccessTokenSubject,
			RefreshTokenSubject:   RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			UserName: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}
	authService, userService := setupServices(cfg)

	server := httpsever.New(cfg, authService, userService)
	server.Serve()

	//http.HandleFunc("/users/register", userRegisterHandler)
	//http.HandleFunc("/health-check", healthCheckHandler)
	//http.HandleFunc("/users/login", userLoginHandler)
	//http.HandleFunc("/users/profile", userPeofileHandler)

}

//func healthCheckHandler(writer http.ResponseWriter, req *http.Request) {
//	fmt.Fprintf(writer, `{"message": "everything is good!"}`)
//}
//
//func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
//	if req.Method != http.MethodPost {
//		fmt.Fprintf(writer, `{"error": "invalid method"}`)
//	}
//
//	data, err := io.ReadAll(req.Body)
//	if err != nil {
//		//fmt.Fprintf(writer, `{"error": "reading request body failed"}`)
//		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
//		return
//	}
//
//	var lReq userservice.LoginRequest
//	err = json.Unmarshal(data, &lReq)
//
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
//		return
//	}
//	authSvc := authservice.New(Jwt_SignKey, AccessTokenExpireDuration, RefreshTokenExpireDuration,
//		AccessTokenSubject, RefreshTokenSubject)
//
//	mysqlRepo := mysql.New()
//	userSvc := userservice.NewService(authSvc, mysqlRepo)
//
//	resp, err := userSvc.Login(lReq)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
//		return
//	}
//	data, err = json.Marshal(resp)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
//		return
//	}
//
//	writer.Write(data)
//}
//
//func userPeofileHandler(writer http.ResponseWriter, req *http.Request) {
//	if req.Method != http.MethodGet {
//		fmt.Fprintf(writer, `{"error": "invalid method"}`)
//	}
//
//	authSvc := authservice.New(Jwt_SignKey, AccessTokenExpireDuration, RefreshTokenExpireDuration,
//		AccessTokenSubject, RefreshTokenSubject)
//
//	authToken := req.Header.Get("Authorization")
//	claims, err := authSvc.ParseToken(authToken)
//	if err != nil {
//		fmt.Fprintf(writer, `{"error": "token is not valid"}`)
//	}
//	mysqlRepo := mysql.New()
//	userSvc := userservice.NewService(authSvc, mysqlRepo)
//
//	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
//
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
//		return
//	}
//
//	data, err := json.Marshal(resp)
//	if err != nil {
//		writer.Write([]byte(fmt.Sprintf("error:%s", err.Error())))
//		return
//	}
//
//	writer.Write(data)
//}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.NewService(authSvc, MysqlRepo)
	return authSvc, userSvc

}
