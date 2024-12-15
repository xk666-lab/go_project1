package responses

import (
	userModel "blog/internal/modules/user/models"
	"fmt"
)

type User struct {
	ID    uint
	Image string
	Name  string
	Email string
}

type Users struct {
	Data []User
}

func ToUser(user userModel.User) User {
	// 如果用户有上传头像，使用上传的头像
	avatarPath := user.Avatar
	if avatarPath == "" {
		// 如果没有上传头像，使用默认的头像生成服务
		avatarPath = fmt.Sprintf("https://ui-avatars.com/api/?name=%s", user.Name)
	} else {
		// 确保头像路径以"/"开头
		if avatarPath[0] != '/' {
			avatarPath = "/" + avatarPath
		}
	}

	return User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Image: avatarPath,
	}
}
