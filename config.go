package main

import (
	// "crypto/md5"
	// "encoding/hex"
	"encoding/json"
	"fmt"

	// "hash/fnv"
	// "io"
	// "log"
	// "math/rand"
	"net/http"
	// "reflect"
	"strconv"
	"strings"
	"time"

	// "github.com/go-redis/redis"
	// uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

const (
	version_name_of_app = "testGo.app._ver=0.0.1"
	ServerAppName       = "testGo"
	default_server_port = ":60990"          //this server port
	mongod_main_one     = "127.0.0.1:27017" //mongo db connect host and port

	//mongodb db name and collection name
	mongod_truck_db   = "test_mgdb"
	mongod_coll_auths = "auth_users"

	//api root level
	gateway_word           = "golang"
	server_api_root        = "test"
	server_api_level_login = "/auth"
	server_api_level_token = "/auth/token"

	server_api_level_manager     = "/mng"
	server_api_level_manager_add = "/mng/add"
	server_api_level_manager_del = "/mng/del"
	server_api_level_manager_get = "/mng/get"
)

const OpErr = "err"
const BracketLeft = "{"
const BracketRight = "}"
const SquareLeft = "["
const SquareRight = "]"

const FeedbackNameTrace = "\"TraceBackServer\""
const FeedbackNameRes = "\"res\""
const FeedbackNameError = "\"error\""
const data_insert_ok = "inserted"
const data_insert_failed = "failed to inserted"
const RequestGoodResult = "I_AM_A_POLAR_BEAR"
const RequestBadResult = "HANABI_ENCORE"

const AllowedHeaders = "Accept,Access-Control-Allow-Headers,POST,PUT,DELETE,GET,PATCH" +
		"Access-Control-Allow-Methods,Access-Control-Allow-Origin,Access-Control-Expose-Headers"

//classes of entity
type DataStore struct {
	session *mgo.Session
}

type AuthUserInfo struct {
	AauthLevel    string          `bson:"authLevel" json:"authLevel"`
	HumanName     string          `bson:"humanName" json:"humanName"`
	HumanGroup    string          `bson:"humanGroup" json:"humanGroup"`
	Passwd        string          `bson:"passwd" json:"passwd"`
	Info          string          `bson:"info" json:"info"`
	CustomContext []CustomContext `bson:"customContext" json:"customContext"`
	DeadOrAlive   string          `bson:"deadOrAlive" json:"deadOrAlive"`
	BirthTime     string          `bson:"birthTime" json:"birthTime"`
}

type CustomContext struct {
	ParamName  string `bson:"paramName" json:"paramName"`
	ParamType  string `bson:"paramType" json:"paramType"`
	ParamValue string `bson:"paramValue" json:"paramValue"`
	ShowOrNot  string `bson:"showOrNot" json:"showOrNot"`
}

type CommonDataFeed struct {
	ErrorCode int    `json:"error_code" bson:"error_code"`
	Data      string `json:"data" bson:"data"`
}

type ResultFeedback struct {
	Res   int    `json:"res" bson:"res"`
	Error string `json:"error" bson:"error"`
}

func (f ResultFeedback) ErrorFeedbacks() string {
	var b strings.Builder
	b.WriteString(BracketLeft)
	b.WriteString(FeedbackNameRes)
	b.WriteString(": ")
	b.WriteString(strconv.Itoa(f.Res))
	b.WriteString(", ")
	b.WriteString(FeedbackNameError)
	b.WriteString(": ")
	b.WriteString("\"")
	b.WriteString(f.Error)
	b.WriteString("\"")
	b.WriteString(BracketRight)
	return b.String()
}

func (c CommonDataFeed) CommonDataPrintOff() string {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("CommonDataPrintOff Error: ", err)
		}
	}()
	jb, _ := json.Marshal(c)
	return string(jb[:])
}

var authCollection *mgo.Collection

func InitDataHandlers() {
	db := &DataStore{session: ConnectToDB()}
	db.InitDataStoreCrazyHandlers()
}

func (db *DataStore) InitDataStoreCrazyHandlers() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("initiation failed because mongodb connection down:", err)
		}
	}()
	var dbSession = db.session
	authCollection = InitCreateColl(mongod_coll_auths, dbSession)
}

func InitCreateColl(collName string, session *mgo.Session) *mgo.Collection {
	collect := session.DB(mongod_truck_db).C(collName)
	if collect == nil {
		fmt.Println("collection initialize failure")
		collect := mgo.Collection{}
		collect.Name = collName
		cInfo := mgo.CollectionInfo{}
		collect.Create(&cInfo)
	}
	fmt.Println("create collection normally:", collName)
	return collect
}

func ConnectToDB() *mgo.Session {
	session, err := mgo.Dial(mongod_main_one)
	if err != nil {
		panic(err)
		fmt.Println("connect DB err:", err)
	}
	session.SetMode(mgo.Monotonic, true)
	fmt.Println("connect to mongoDB successfully")
	return session
}

func MongoClient() *mgo.Session {
	session, err := mgo.Dial(mongod_main_one)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

func MongoDBConnectionCheck() {
	fmt.Println("start check")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("mongodb connection check result is", err)
		} else {
			fmt.Println("status is successful")
		}
	}()
	time.Sleep(time.Second * 30)
	for {
		err := dataStore.session.Ping()
		if err != nil {
			fmt.Println(err)
			fmt.Println("init Data Handlers again")
			InitDataHandlers()
			InitDataStoreHandlers()
		}
		time.Sleep(time.Second * 30)
	}
	fmt.Println("finish check")
}

func main() {
	fmt.Println("Wellcome to the go demo test")
	InitDataStoreHandlers()
	InitDataHandlers()
	// handlersMain(server_api_root)
	go MongoDBConnectionCheck()
	handlersMain(server_api_root)
	// CheckLogin("test001")
}

func CheckUserInfo(username string) {
	fmt.Println("login start:", username)
}

//
//Send back with custom headers
func SendToHellCommonHead(name string, sendback string, status int, w http.ResponseWriter) {
	err := recover()
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(status)
		w.Write([]byte(sendback))
	}
}

//Send back with common headers. no custom.
func SendBackCommonHead(back string, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	w.Write([]byte(back))
}

//
func ErrorResultFill(result int, errorDesc string) string {
	res := new(ResultFeedback)
	res.Res = result
	res.Error = errorDesc
	return res.ErrorFeedbacks()
}
