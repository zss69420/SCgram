package service

import (
	"fmt"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {
	//select * from xxx where username and password
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("username", username))
	query.Must(elastic.NewTermQuery("password", password))

	//readfromES with input query in USER INDEX
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	var utype model.User
	//if it is utype, go to loop, cast item to model.user, if match, return true, if all no, return false
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if u.Password == password {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}
	return false, nil
}

//ES is not like MYSQL, which will return duplicate primary key error if username existed; so checke it by yourself here
func AddUser(user *model.User) (bool, error) {
	//check if this user has been registered before, if yes, return false, if no go ahead
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	//totalHits > 0 it means with the query it indeed has a search result, so it means this username has existed in ES before
	if searchResult.TotalHits() > 0 {
		return false, nil
	}

	//if no, save this user info to ES
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}
