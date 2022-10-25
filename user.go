package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/caarlos0/httperr"
)

func (api *API) GetUser(w http.ResponseWriter, req *http.Request) error {
	id := req.URL.Query().Get("userId")
	if id == "" {
		return httperr.Errorf(http.StatusBadRequest, "")
	}

	rows, err := api.db.Query(fmt.Sprintf(`SELECT id,name,email FROM users WHERE id = '%s' limit 1`, id))
	if err != nil {
		log.Print("Error: ", err)
		return err
	}

	users, err := api.readUsers(rows)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(
		"Content-Type", "application/json")
	return json.NewEncoder(w).Encode(users[0])
}

func (api *API) Users(w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "GET":
		return api.GetUsers(w, req)
	default:
		return httperr.Errorf(http.StatusMethodNotAllowed, "")
	}
}

func (api *API) User(w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "GET":
		return api.GetUser(w, req)
	case "POST":
		return api.UpdateUser(w, req)
	case "PUT":
		return api.CreateUser(w, req)
	case "DELETE":
		return api.DeleteUser(w, req)
	default:
		return httperr.Errorf(http.StatusMethodNotAllowed, "")
	}
}

type User struct {
	ID       int    `json:"id" db:"ID"`
	Name     string `json:"name" db:"Name"`
	Email    string `json:"email" db:"Email"`
	Password string `json:"password" db:"Password"`
	Role     string `json:"role,omitempty" db:"Role"`
}

func (api *API) GetUsers(w http.ResponseWriter, req *http.Request) error {
	if req.URL.Query().Get("isAdmin") != "1" {
		return httperr.Errorf(http.StatusForbidden, "")
	}

	rows, err := api.db.Query("SELECT id,name,email FROM users")
	if err != nil {
		return err
	}

	users, err := api.readUsers(rows)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(users)
}

func (api *API) readUsers(rows *sql.Rows) ([]User, error) {
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (api *API) CreateUser(w http.ResponseWriter, req *http.Request) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return httperr.Errorf(http.StatusBadRequest, "")
	}

	defer req.Body.Close()

	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		return httperr.Wrap(err, http.StatusBadRequest)
	}

	if user.Name == "" || user.Email == "" {
		return httperr.Wrap(err, http.StatusBadRequest)
	}

	id, err := NewUser(api.db, user.Name, user.Email, Hash(user.Password))
	if err != nil {
		return err
	}

	// Update the user with the new id and clear the password
	user.ID = id
	user.Password = ""

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (api *API) UpdateUser(w http.ResponseWriter, req *http.Request) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return httperr.Errorf(http.StatusMethodNotAllowed, "")
	}

	defer req.Body.Close()

	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		return httperr.Wrap(err, http.StatusMethodNotAllowed)
	}

	if user.ID == 0 || user.Name == "" || user.Email == "" {
		return httperr.Errorf(http.StatusMethodNotAllowed, "")
	}

	_, err = api.db.Exec(` UPDATE users SET name = ?, email = ? WHERE id = ? `, user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(user)
}

func (api *API) DeleteUser(w http.ResponseWriter, req *http.Request) error {
	id := req.URL.Query().Get("userId")
	if id == "" {
		return httperr.Errorf(http.StatusBadRequest, "")
	}

	q := fmt.Sprintf("DELETE FROM users WHERE id = '%s'", id)
	_, err := api.db.Exec(q)
	if err != nil {
		return err
	}

	if _, err := strconv.Atoi(id); err != nil {
		fmt.Print(hacked)
		err = nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
