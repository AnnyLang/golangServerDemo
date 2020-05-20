package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func CheckLogin(username string) string {
	coll := authCollection
	res := make(map[string]string, 0)
	errFind := coll.Find(bson.M{"humanName": username}).One(&res)
	if errFind != nil {
		fmt.Println("find err:", errFind)
	}
	fmt.Println("find result is :", res["humanName"])
	fmt.Println("result is:",toString(res))
	return res["humanName"]
}

func toString(data map[string]string)string{
	var txt=""
	if data!=nil{
		for key,value:=range data{
			txt+=key+"="+value+","
		}
	}
	return txt
}