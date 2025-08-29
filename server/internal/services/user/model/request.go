package model

import (
	"time"

	"github.com/google/uuid"

	userRepoModel "server/internal/services/user/repository/model"
	"server/internal/utils/necessary"
)

type CreateReq struct {
	Name            string
	Email           string
	PasswordHash    []byte
	PasswordSalt    []byte
	TimeCreate      time.Time
	DefaultCurrency string
}

type GetUsersReq struct {
	Necessary necessary.NecessaryUserInformation
	IDs       []uuid.UUID
	Emails    []string
}

func (s GetUsersReq) SetNecessary(information necessary.NecessaryUserInformation) any {
	s.Necessary = information
	return s
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

func (s UpdateUserReq) SetNecessary(information necessary.NecessaryUserInformation) any {
	s.Necessary = information
	return s
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
