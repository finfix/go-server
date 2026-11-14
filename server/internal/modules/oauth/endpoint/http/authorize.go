package http

import (
	"net/url"

	"server/internal/modules/oauth/model"
)

// authorizeReq собирает параметры авторизации из query или полей формы
func authorizeReq(values url.Values) model.AuthorizeReq {
	return model.AuthorizeReq{
		ClientID:            parseUUID(values.Get("client_id")),
		RedirectURI:         values.Get("redirect_uri"),
		ResponseType:        values.Get("response_type"),
		CodeChallenge:       values.Get("code_challenge"),
		CodeChallengeMethod: values.Get("code_challenge_method"),
		Resource:            values.Get("resource"),
		State:               values.Get("state"),
	}
}
