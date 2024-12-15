package repositories

import (
	userModel "blog/internal/modules/user/models"
	"blog/pkg/database"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func New() *UserRepository {
	return &UserRepository{
		DB: database.Connection(),
	}
}

func (userRepository *UserRepository) Create(user userModel.User) userModel.User {
	// 创建用户
	result := userRepository.DB.Create(&user)
	if result.Error != nil {
		return userModel.User{}
	}

	// 重新查询用户以获取完整信息
	var newUser userModel.User
	userRepository.DB.First(&newUser, user.ID)
	return newUser
}

func (userRepository *UserRepository) FindByEmail(email string) userModel.User {
	var user userModel.User
	userRepository.DB.First(&user, "email = ?", email)
	return user
}

func (userRepository *UserRepository) FindByID(id int) userModel.User {
	var user userModel.User
	userRepository.DB.First(&user, "id = ?", id)
	return user
}
