package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	proto_user "git.ve.home/nicolasc/linotte/services/user/proto"
)

type AuthenticationController struct {
	user proto_user.UserClient
}

func NewAuthenticationController(user proto_user.UserClient) *AuthenticationController {
	return &AuthenticationController{user}
}

func (controller *AuthenticationController) Login(w http.ResponseWriter, r *http.Request) {
	var temp struct {
		Username string
		Password string
	}

	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &temp)

	reply, err := controller.user.Login(
		context.Background(),
		&proto_user.LoginRequest{
			Login:    temp.Username,
			Password: temp.Password,
		},
	)

	if err == nil && reply.Success {
		json.NewEncoder(w).Encode(reply)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	helpers.WriteError(w, "Login failed", err)
}
