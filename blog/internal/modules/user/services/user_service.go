package services

import (
	userModel "blog/internal/modules/user/models"
	UserRepository "blog/internal/modules/user/repositories"
	"blog/internal/modules/user/requests/auth"
	UserResponse "blog/internal/modules/user/responses"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type UserService struct {
	userRepository UserRepository.UserRepositoryInterface
}

func New() *UserService {
	return &UserService{
		userRepository: UserRepository.New(),
	}
}

func (userService *UserService) Create(request auth.RegisterRequest) (UserResponse.User, error) {
	var response UserResponse.User
	var user userModel.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return response, errors.New("error hashing the password")
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Password = string(hashedPassword)

	// 如果上传了头像，处理头像上传
	if request.Avatar != nil {
		log.Printf("开始处理头像上传: %s", request.Avatar.Filename)
		
		// 确保上传���录存在
		uploadDir := "assets/uploads/avatars"
		if err := os.MkdirAll(uploadDir, 0777); err != nil {
			log.Printf("创建目录失败: %v", err)
			return response, fmt.Errorf("创建上传目录失败: %v", err)
		}

		// 生成文件名
		ext := filepath.Ext(request.Avatar.Filename)
		filename := fmt.Sprintf("%s%s", user.Email, ext) // 使用邮箱作为文件名
		avatarPath := fmt.Sprintf("/assets/uploads/avatars/%s", filename)
		fullPath := filepath.Join(uploadDir, filename)

		log.Printf("准备保存头像到: %s", fullPath)

		// 检查目录权限
		if err := os.Chmod(uploadDir, 0777); err != nil {
			log.Printf("修改目录权限失败: %v", err)
		}

		// 打开上传的文件
		src, err := request.Avatar.Open()
		if err != nil {
			log.Printf("打开上传文件失败: %v", err)
			return response, fmt.Errorf("打开上传文件失败: %v", err)
		}
		defer src.Close()

		// 创建目标文件
		dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("创建目标文件失败: %v", err)
			return response, fmt.Errorf("创建目标文件失败: %v", err)
		}
		defer dst.Close()

		// 复制文件内容
		if _, err = io.Copy(dst, src); err != nil {
			log.Printf("复制文件失败: %v", err)
			return response, fmt.Errorf("复制文件失败: %v", err)
		}

		// 检查文件是否成功创建
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			log.Printf("文件未成功创建: %v", err)
			return response, fmt.Errorf("文件未成功创建: %v", err)
		}

		log.Printf("头像上传成功: %s", avatarPath)
		user.Avatar = avatarPath
	}

	newUser := userService.userRepository.Create(user)

	if newUser.ID == 0 {
		return response, errors.New("创建用户失败")
	}

	return UserResponse.ToUser(newUser), nil
}

func (userService *UserService) CheckUserExists(email string) bool {
	user := userService.userRepository.FindByEmail(email)
	return user.ID != 0
}

func (userService *UserService) HandleUserLogin(request auth.LoginRequest) (UserResponse.User, error) {
	var response UserResponse.User

	user := userService.userRepository.FindByEmail(request.Email)

	if user.ID == 0 {
		return response, errors.New("用户不存在")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return response, errors.New("密码错误")
	}

	return UserResponse.ToUser(user), nil
}
