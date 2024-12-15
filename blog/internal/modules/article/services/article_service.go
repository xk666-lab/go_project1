package services

import (
	ArticleModel "blog/internal/modules/article/models"
	ArticleRepository "blog/internal/modules/article/repositories"
	"blog/internal/modules/article/requests/articles"
	ArticleResponse "blog/internal/modules/article/responses"
	UserResponse "blog/internal/modules/user/responses"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ArticleService struct {
	articleRepository ArticleRepository.ArticleRepositoryInterface
}

func New() *ArticleService {
	return &ArticleService{
		articleRepository: ArticleRepository.New(),
	}
}

func (articleService *ArticleService) GetFeaturedArticles() ArticleResponse.Articles {
	articles := articleService.articleRepository.List(4)

	return ArticleResponse.ToArticles(articles)
}

func (articleService *ArticleService) GetStoriesArticles() ArticleResponse.Articles {
	articles := articleService.articleRepository.List(6)

	return ArticleResponse.ToArticles(articles)
}

func (articleService *ArticleService) Find(id int) (ArticleResponse.Article, error) {
	var response ArticleResponse.Article

	article := articleService.articleRepository.Find(id)

	if article.ID == 0 {
		return response, errors.New("article not found")
	}

	return ArticleResponse.ToArticle(article), nil
}

func (articleService *ArticleService) StoreAsUser(request articles.StoreRequest, user UserResponse.User) (ArticleResponse.Article, error) {
	var article ArticleModel.Article
	var response ArticleResponse.Article
	
	article.Title = request.Title
	article.Content = request.Content
	article.UserID = user.ID

	// 先创建文章记录以获取ID
	newArticle := articleService.articleRepository.Create(article)

	if newArticle.ID == 0 {
		return response, errors.New("error in creating the article")
	}

	// 生成唯一的文件名（使用新创建的文章ID）
	ext := filepath.Ext(request.Image.Filename)
	filename := fmt.Sprintf("%d%s", newArticle.ID, ext)
	
	// 设置图片保存路径（相对于项目根目录）
	uploadDir := "assets/uploads"
	imagePath := fmt.Sprintf("/%s/%s", uploadDir, filename)
	
	// 更新文章的图片路径
	newArticle.Image = imagePath
	newArticle = articleService.articleRepository.Update(newArticle)

	// 打开上传的文件
	src, err := request.Image.Open()
	if err != nil {
		return response, errors.New("error opening uploaded file")
	}
	defer src.Close()

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return response, errors.New("error creating upload directory")
	}

	// 创建目标文件
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return response, errors.New("error creating destination file")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return response, errors.New("error copying file")
	}

	return ArticleResponse.ToArticle(newArticle), nil
}
