package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"blog_server/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

// ArticleController 结构体用于处理文章相关的请求。
// 它实现了 IArticleController 接口，并包含一个 *gorm.DB 类型的字段用于数据库操作。
type ArticleController struct {
	DB *gorm.DB
}

// IArticleController 接口定义了文章控制器需要实现的一系列方法。
type IArticleController interface {
	Create(c *gin.Context) // 创建文章的方法
	Update(c *gin.Context) // 更新文章的方法
	Delete(c *gin.Context) // 删除文章的方法
	Show(c *gin.Context)   // 显示文章详情的方法
	List(c *gin.Context)   // 列出所有文章的方法

}

// Create 方法实现 IArticleController 接口的创建文章功能。
// 它首先解析请求体中的 JSON 数据，然后创建并保存新文章到数据库。
func (a ArticleController) Create(c *gin.Context) {
	// vo.CreateArticleRequest 是一个值对象，用于接收和验证创建文章请求的数据。
	var articleRequest vo.CreateArticleRequest

	// c.ShouldBindJSON 解析请求体中的 JSON 数据到 articleRequest 变量。
	// 如果解析失败，返回错误响应。
	if err := c.ShouldBindJSON(&articleRequest); err != nil {
		response.Fail(c, nil, "数据错误")
		return
	}

	// 从 Gin 上下文中获取当前登录的用户信息。
	// 这里假设 Gin 上下文中存储的用户信息的键是 "user"。
	user, _ := c.Get("user")

	// 创建一个 model.Article 类型的实例，并填充数据。
	// 这里假设 model.User 类型有一个 ID 字段。
	article := model.Article{
		UserId:     user.(model.User).ID,      // 从当前登录用户获取 ID。
		CategoryId: articleRequest.CategoryId, // 从请求数据获取分类 ID。
		Title:      articleRequest.Title,      // 从请求数据获取标题。
		Content:    articleRequest.Content,    // 从请求数据获取内容。
		HeadImage:  articleRequest.HeadImage,  // 从请求数据获取头图链接。
	}

	// 使用 GORM 创建并保存新文章到数据库。
	// 如果创建失败，返回错误响应。
	if err := a.DB.Create(&article).Error; err != nil {
		response.Fail(c, nil, "发布失败")
		return
	}

	// 如果创建成功，返回成功响应和新创建的文章 ID。
	response.Success(c, gin.H{"id": article.ID}, "发布成功")
}

// Update 方法实现 IArticleController 接口的更新文章功能。
// 它首先解析请求体中的 JSON 数据，然后根据文章 ID 查找并更新文章。
func (a ArticleController) Update(c *gin.Context) {
	var articleRequest vo.CreateArticleRequest
	if err := c.ShouldBindJSON(&articleRequest); err != nil {
		response.Fail(c, nil, "数据错误")
		return
	}
	articleId := c.Params.ByName("id") // 从请求的 URL 参数中获取文章 ID。
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(c, nil, "文章不存在")
		return
	}
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(c, nil, "登录用户不正确")
		return
	}
	if err := a.DB.Model(&article).Updates(articleRequest).Error; err != nil {
		response.Fail(c, nil, "修改失败")
		return
	}
	response.Success(c, nil, "修改成功")
}

// Delete 方法实现 IArticleController 接口的删除文章功能。
// 它根据文章 ID 查找并删除文章。
func (a ArticleController) Delete(c *gin.Context) {
	articleId := c.Params.ByName("id")
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(c, nil, "文章不存在")
		return
	}
	user, _ := c.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(c, nil, "登录用户不正确")
		return
	}
	if err := a.DB.Delete(&article).Error; err != nil {
		response.Fail(c, nil, "删除失败")
		return
	}
	response.Success(c, nil, "删除成功")
}

// Show 方法实现 IArticleController 接口的显示文章详情功能。
// 它根据文章 ID 查找并显示文章的详细信息。
func (a ArticleController) Show(c *gin.Context) {
	articleId := c.Params.ByName("id")
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(c, nil, "文章不存在")
		return
	}
	response.Success(c, gin.H{"article": article}, "查找成功")
}

// List 方法实现 IArticleController 接口的列出所有文章功能。
// 它可以根据关键词、分类 ID 和分页参数来过滤和列出文章。
func (a ArticleController) List(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	categoryId := c.DefaultQuery("categoryId", "0")
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	var query []string
	var args []interface{}

	if keyword != "" {
		query = append(query, "(title LIKE ? OR content LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}
	if categoryId != "0" {
		query = append(query, "category_id = ?")
		args = append(args, categoryId)
	}

	var article []model.ArticleInfo
	var count int
	var querystr string
	if len(query) > 0 {
		querystr = strings.Join(query, " AND ")
	}

	switch len(args) {
	case 0:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").
			Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Count(&count)
	case 1:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").
			Where(querystr, args[0]).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(querystr, args[0]).Count(&count)
	case 2:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").
			Where(querystr, args[0], args[1]).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(querystr, args[0], args[1]).Count(&count)
	case 3:
		a.DB.Table("articles").Select("id, category_id, title, LEFT(content,80) AS content, head_image, created_at").
			Where(querystr, args[0], args[1], args[2]).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&article)
		a.DB.Model(model.Article{}).Where(querystr, args[0], args[1], args[2]).Count(&count)
	}

	response.Success(c, gin.H{"article": article, "count": count}, "查找成功")
}
func (ac *ArticleController) ShowWithComments(c *gin.Context) {
	articleID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		response.Fail(c, nil, "无效的文章ID")
		return
	}

	var article model.Article
	if err := ac.DB.Preload("Comment").Where("id = ?", articleID).First(&article).Error; err != nil {
		response.Fail(c, nil, "文章未找到")
		return
	}

	response.Success(c, gin.H{"article": article}, "成功")
}

// NewArticleController 函数用于创建并初始化 ArticleController 实例。
// 它获取数据库连接，执行自动迁移，并返回一个实现了 IArticleController 接口的控制器实例。
func NewArticleController() IArticleController {
	db := common.GetDB()              // 从 common 包中获取数据库连接
	db.AutoMigrate(model.Article{})   // 使用 GORM 自动迁移 Article 模型
	return &ArticleController{DB: db} // 返回初始化好的 ArticleController 实例
}
