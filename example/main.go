// Package main @author K·J Create at 2019-03-04 15:55
package main

import (
	"fmt"
	"net/url"
	"github.com/koori69/queryize"
)

type User struct {
	Name  string  `query:"name"`
	Age   *int    `query:"age"`
	Sex   int     `query:"sex"`
	DepID *string `query:"dep_id"`
	me    *string `query:"me"`
}

func main() {
	err := unmarshl()
	fmt.Println(err)
}

func unmarshl() error {
	s := "name=Tome&age=12&sex=0&dep_id=math&me=AD"
	var user User
	query, err := url.ParseQuery(s)
	if nil != err {
		return err
	}
	err = queryize.ConfigDefault.Unmarshal(query, &user)
	if nil != err {
		return err
	}
	fmt.Printf("%+v\n", user)
	str, err := queryize.ConfigDefault.Marshal(user)
	fmt.Println(str)
	str, err = queryize.ConfigDefault.Marshal(&user)
	fmt.Println(str)
	return nil
}
