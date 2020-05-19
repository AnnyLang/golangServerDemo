package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	r.HandleFunc(path+server_api_level_login, UserAuth).Methods("GET")
	http.ListenAndServe(default_server_port,
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE", "PATCH"}),
			handlers.AllowedHeaders(strings.Split(AllowedHeaders, ",")),
			handlers.AllowedOrigins([]string{"*"}))(r))
}

//
func UserAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start server to find user info by username")
	defer SendToHellCommonHead("UserAuth", ErrorResultFill(1, "request failure "), http.StatusNoContent, w)
	query := r.URL.Query()
	userName := query.Get("username")
	res := CheckLogin(userName)
	c := CommonDataFeed{ErrorCode: 0, Data: res}
	SendBackCommonHead(c.CommonDataPrintOff(), http.StatusAccepted, w)
}
