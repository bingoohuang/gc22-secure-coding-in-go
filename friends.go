package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/httperr"
	"github.com/jmoiron/sqlx"
)

func (api *API) GetFriends(w http.ResponseWriter, req *http.Request) error {
	userID := req.URL.Query().Get("userId")

	friends, err := GetFriends(userID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	data, err := json.Marshal(friends)
	if err != nil {
		log.Fatal(err)
		return err
	}

	w.Write(data)
	return nil
}

func GetFriends(userId string) ([]*User, error) {
	q := fmt.Sprintf(`SELECT users.* FROM users JOIN friends ON users.ID = friends.FriendId WHERE friends.UserId = '%s'; `, userId)

	return Query[*User](q)
}

func Query[T any](query string) ([]T, error) {
	dbx, err := sqlx.Open("sqlite", db)
	if err != nil {
		return nil, err
	}

	var results []T
	err = dbx.Select(&results, query)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (api *API) Friends(w http.ResponseWriter, req *http.Request) error {
	log.Print("Users")
	switch req.Method {
	case "GET":
		return api.GetFriends(w, req)
	default:
		return httperr.Errorf(http.StatusMethodNotAllowed, "")
	}
}

func (api *API) Friend(w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "PUT":
		return api.AddFriend(w, req)
	default:
		return httperr.Errorf(http.StatusMethodNotAllowed, "")
	}
}

func (api *API) AddFriend(w http.ResponseWriter, req *http.Request) error {
	userID := req.URL.Query().Get("userId")
	friendID := req.URL.Query().Get("friendId")

	log.Printf("Adding friend [%s] to database for user [%s]", friendID, userID)
	err := AddFriend(api.db, userID, friendID)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
