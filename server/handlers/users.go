package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"server/db"
	errs "server/errors"
)

type ReducedUserDetail struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func ListUsersHandler(sqlDB *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userEntities, err := db.GetUsers(sqlDB, req.Context())
		if err != nil {
			errs.APIError(w, err)
			return
		}

		var users []ReducedUserDetail
		for _, userEntity := range userEntities {
			user := ReducedUserDetail{
				UserId:    userEntity.UserId,
				FirstName: userEntity.FirstName,
				LastName:  userEntity.LastName,
			}
			users = append(users, user)
		}

		resBody, err := json.Marshal(users)
		if err != nil {
			errs.APIError(w, jsonMarshallingError(err))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

type UserDetail struct {
	UserId    string   `json:"userId"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Age       int32    `json:"age"`
	Aliases   []string `json:"aliases"`
}

func RetrieveUserHandler(sqlDB *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		userEntity, err := db.GetUserById(sqlDB, req.Context(), userId)
		if err != nil {
			errs.APIError(w, err)
			return
		}

		user := UserDetail{
			UserId:    userEntity.UserId,
			Email:     userEntity.Email,
			FirstName: userEntity.FirstName,
			LastName:  userEntity.LastName,
			Age:       userEntity.Age,
			Aliases:   userEntity.Aliases,
		}

		resBody, err := json.Marshal(user)
		if err != nil {
			errs.APIError(w, jsonMarshallingError(err))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func CreateUserHandler(sqlDB *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var reqBody db.UpdateUserPayload
		err := json.NewDecoder(req.Body).Decode(&reqBody)
		if err != nil {
			errs.APIError(w, jsonDecodeError(err))
			return
		}

		userEntity, err := db.CreateUser(sqlDB, req.Context(), reqBody)
		if err != nil {
			errs.APIError(w, err)
			return
		}

		newUser := UserDetail{
			UserId:    userEntity.UserId,
			Email:     userEntity.Email,
			FirstName: userEntity.FirstName,
			LastName:  userEntity.LastName,
			Age:       userEntity.Age,
			Aliases:   userEntity.Aliases,
		}

		resBody, err := json.Marshal(newUser)
		if err != nil {
			errs.APIError(w, jsonMarshallingError(err))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func UpdateUserHandler(sqlDB *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		var reqBody db.UpdateUserPayload
		err := json.NewDecoder(req.Body).Decode(&reqBody)
		if err != nil {
			errs.APIError(w, jsonDecodeError(err))
			return
		}

		userEntity, err := db.UpdateUser(sqlDB, req.Context(), userId, reqBody)
		if err != nil {
			errs.APIError(w, err)
			return
		}

		resBody, err := json.Marshal(UserDetail{
			UserId:    userEntity.UserId,
			Email:     userEntity.Email,
			FirstName: userEntity.FirstName,
			LastName:  userEntity.LastName,
			Age:       userEntity.Age,
			Aliases:   userEntity.Aliases,
		})
		if err != nil {
			errs.APIError(w, jsonMarshallingError(err))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func DeleteUserHandler(sqlDB *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		db.DeleteUser(sqlDB, req.Context(), userId)

		w.WriteHeader(http.StatusNoContent)
	}

	return http.HandlerFunc(fn)
}

func jsonDecodeError(err error) error {
	return fmt.Errorf("failed to decode json: %v", err)
}

func jsonMarshallingError(err error) error {
	return fmt.Errorf("error serializing json: %v", err)
}
