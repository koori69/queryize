# Queryize

## Feature
- convert struct to string , support unexported  field (ptr only), e.g parameters of GET request
- convert map to struct , support unexported  field (ptr only), e.g url.Values to struct

## Installation

Install:

```shell
go get -u github.com/koori69/queryize
```

Import

```go
import "github.com/koori69/queryize"
```

## Quickstart

```go
type User struct {
	Name  string  `query:"name"`
	Age   *int    `query:"age"`
	Sex   int     `query:"sex"`
	DepID *string `query:"dep_id"`
	self  string  `query:"self"`
}

func example() error {
	s := "name=Tome&age=12&sex=0&dep_id=math&self=AD"
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

```

