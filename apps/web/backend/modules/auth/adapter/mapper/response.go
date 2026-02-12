package mapper

import (
	"github.com/kubrickcode/specvital/apps/web/backend/internal/api"
	"github.com/kubrickcode/specvital/apps/web/backend/modules/auth/domain/entity"
)

func ToUserInfo(user *entity.User) api.UserInfo {
	return api.UserInfo{
		AvatarURL: user.AvatarURL,
		ID:        user.ID,
		Login:     user.Username,
		Name:      nil,
	}
}

func ToLoginResponse(authURL string) api.LoginResponse {
	return api.LoginResponse{
		AuthURL: authURL,
	}
}
