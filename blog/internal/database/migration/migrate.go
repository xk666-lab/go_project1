package migration

import (
	articleModels "blog/internal/modules/article/models"
	userModels "blog/internal/modules/user/models"
	"blog/pkg/database"
	"log"
)

func Migrate() {
	db := database.Connection()
	log.Println("开始数据库迁移...")

	// 删除现有的表（如果存在）
	log.Println("删除现有的表...")
	db.Migrator().DropTable(&userModels.User{}, &articleModels.Article{})

	// 创建新表
	log.Println("创建新表...")
	err := db.AutoMigrate(&userModels.User{}, &articleModels.Article{})
	if err != nil {
		log.Fatalf("迁移失败: %v", err)
		return
	}

	log.Println("数据库迁移完成")
}
