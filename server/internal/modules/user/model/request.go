package model

import (
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
	"time"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	userRepoModel "server/internal/modules/user/repository/model"
)

type CreateReq struct {
	ID              uuid.UUID
	Name            string
	Email           string
	PasswordHash    []byte
	PasswordSalt    []byte
	TimeCreate      time.Time
	DefaultCurrency string
	IsAdmin         bool
}

type GetUserReq struct {
	Necessary necessary.NecessaryUserInformation
	AccessToken string
}

// ProtoGetUserReq wrapper for proto request
type ProtoGetUserReq struct {
	*proto.GetUserRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetUserReq) ConvertToModel() (GetUserReq, error) {
	if p.GetUserRequest == nil {
		return GetUserReq{}, errors.BadRequest.New("GetUserRequest is required")
	}

	return GetUserReq{
		AccessToken: p.AccessToken,
	}, nil
}

type GetUsersReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID
	Emails    []string
}

type UpdateUserReq struct {
	Necessary         necessary.NecessaryUserInformation
	Name              *string `json:"name"`
	Email             *string `json:"email"`
	Password          *string `json:"password"`
	OldPassword       *string `json:"oldPassword"`
	DefaultCurrency   *string `json:"defaultCurrency"`
	NotificationToken *string `json:"notificationToken"`
}

func (s UpdateUserReq) ConvertToRepoModel() userRepoModel.UpdateUserReq {
	return userRepoModel.UpdateUserReq{
		ID:              s.Necessary.UserID,
		Name:            s.Name,
		Email:           s.Email,
		PasswordHash:    nil,
		PasswordSalt:    nil,
		DefaultCurrency: s.DefaultCurrency,
	}
}

// ProtoUpdateUserReq wrapper for proto request
type ProtoUpdateUserReq struct {
	*proto.UpdateUserRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdateUserReq) ConvertToModel() (UpdateUserReq, error) {
	if p.UpdateUserRequest == nil {
		return UpdateUserReq{}, errors.BadRequest.New("UpdateUserRequest is required")
	}

	return UpdateUserReq{
		Name:              p.Name,
		Email:             p.Email,
		Password:          p.Password,
		OldPassword:       p.OldPassword,
		DefaultCurrency:   p.DefaultCurrency,
		NotificationToken: p.NotificationToken,
	}, nil
}
