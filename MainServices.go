package main

import (
	"fmt"
	"encoding/json"
	// "context"
	"gopkg.in/mgo.v2/bson"
)

//how to use LINQ to return map object
func CheckLogin(username string) string {
	coll := authCollection
	res := make(map[string]string, 0)
	errFind := coll.Find(bson.M{"humanName": username}).One(&res)
	if errFind != nil {
		fmt.Println("find err:", errFind)
	}
	fmt.Println("find result is :", res["humanName"])
	fmt.Println("result is:", toString(res))
	CheckLoginObj(username)
	return res["humanName"]
}

//how to use LINQ to return entity object
func CheckLoginObj(username string) []AuthUserInfo {
	coll := authCollection
	res := make([]AuthUserInfo, 0)
	errFind := coll.Find(bson.M{"humanName": username}).All(&res)
	if errFind != nil {
		fmt.Println("find err:", errFind)
	}
	if len(res)>0{
		//for i,item 		==>for array 
		//for key,value		==>for map
		//for i=0;i<5;i++ 	==>for both
		//if you do not need i or key,you can use this formate :>for _,value{}
		for _,item :=range res{
			if len(item.HumanGroup)>0{
				//func json.Marshal return two values:>[value,err] 
				jsonStr,err:=json.Marshal(item)
				if err!=nil{
					fmt.Println(err)
				}
				fmt.Println("item is not nil",string(jsonStr))
			}
			fmt.Println("parseTo:" + parseToString(res[0]))
		}
	}
	return res
}

func login(username string,password string)AuthUserInfo{
	fmt.Println("username="+username+","+"password="+password)
	coll:=authCollection
	res:=AuthUserInfo{}
	errFind:=coll.Find(bson.M{"humanName":username,"passwd":password}).One(&res)
	if errFind!=nil{
		fmt.Println(errFind)
	}
	fmt.Println(res.HumanName+","+res.Passwd)
	return res
}

func insertNewUser(user AuthUserNew){
    docD := bson.D{
        {"name", "Project"},
        {"description", "Project Tasks"},
    }
    	coll:=authCollection
		err:=coll.Insert(docD)
		if err!=nil{
			fmt.Println(err)
		}
    
}

type Category struct {
    Id          bson.ObjectId `bson:"_id,omitempty"`
    Name        string
    Description string
}

//map[keyType]valueType
//map[string]struct same as map(string,object) in java
func toString(data map[string]string) string {
	var txt = ""
	if data != nil {
		for key, value := range data {
			txt += key + "=" + value + ","
		}
	}
	return txt
}

func parseToString(userinfo AuthUserInfo) string {
	var txt = ""
	if len(userinfo.HumanName) > 0 {
		txt += "username=" + userinfo.HumanName + ","
		txt += "password=" + userinfo.Passwd + ","
		txt += "description=" + userinfo.Info
	}
	return txt
}
