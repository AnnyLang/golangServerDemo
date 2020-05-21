package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// "gopkg.in/mgo.v2/bson"
)

var dataStore DataStore

func InitDataStoreHandlers() {
	dataStore.session = ConnectToDB()
}

func deferCommon(name string) {
	err := recover()
	if err != nil {
		fmt.Println(name)
		fmt.Println("handler-error:", err)
	}
}

func checkPath(path string) string {
	if path != "" && path[0:1] != "/" {
		path = "/" + path
	}
	return path
}

func handlersMain(path string) {
	defer deferCommon("HandlersMain")
	path = checkPath(path)
	fmt.Println("server_port=" + default_server_port)
	r := mux.NewRouter().StrictSlash(true)
	//handleFunc(apiRoot,methodName).Methods(requestType)
	r.HandleFunc(path+server_api_level_manager_get, UserAuth).Methods("GET")
	r.HandleFunc(path+server_api_level_login, loginCheck).Methods("POST")
	//
	r.HandleFunc(path+server_api_level_manager_add, insertOne).Methods("POST")
	//start the server
	http.ListenAndServe(default_server_port,
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE", "PATCH"}),
			handlers.AllowedHeaders(strings.Split(AllowedHeaders, ",")),
			handlers.AllowedOrigins([]string{"*"}))(r))
}

//get method
func UserAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start server to find user info by username")
	defer SendToHellCommonHead("UserAuth", ErrorResultFill(1, "request failure "), http.StatusNoContent, w)
	query := r.URL.Query()
	userName := query.Get("username")
	res := CheckLogin(userName)
	c := CommonDataFeed{ErrorCode: 0, Data: res}
	SendBackCommonHead(c.CommonDataPrintOff(), http.StatusAccepted, w)
}

func loginCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method(post) loginCheck")
	defer SendToHellCommonHead("loginCheck", ErrorResultFill(1, "request failure "), http.StatusNoContent, w)
	query := r.URL.Query()
	username := query.Get("username")
	password := query.Get("password")
	res := login(username, password)
	userInfo, _ := json.Marshal(res)
	// c:=CommonResult{Code:0,Data:res}
	CommonWriteBack(w, userInfo)
}

func insertOne(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method(post) insertOne")
	defer SendToHellCommonHead("insertOne", ErrorResultFill(1, "request failure"), http.StatusNoContent, w)
	query := r.URL.Query()
	jsonStr, _ := json.Marshal(query)
	if query != nil {
		fmt.Println("insertOne info is :" + string(jsonStr))
		isLive, _ := strconv.Atoi(query.Get("isLive"))
		level, _ := strconv.Atoi(query.Get("authLevel"))

		userInfo := AuthUserNew{
			// Id:          bson.NewObjectId(),
			AuthLevel:   level,
			HumanName:   query.Get("humanName"),
			HumanGroup:  query.Get("humanGroup"),
			Passwd:      query.Get("password"),
			Info:        query.Get("info"),
			DeadOrAlive: isLive,
		}
		fmt.Println(userInfo)
		insertNewUser(userInfo)
	}
}
