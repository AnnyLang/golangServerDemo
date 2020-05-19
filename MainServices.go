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
	return res["humanName"]
}
