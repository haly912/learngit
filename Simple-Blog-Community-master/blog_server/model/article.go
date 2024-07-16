package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// model/article.go

// Article 定义了文章的数据模型，与数据库中的文章表相对应。
type Article struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"`   // 文章的唯一标识符，使用 UUID。
	UserId     uint      `json:"user_id" gorm:"not null"`                // 文章作者的用户 ID。
	CategoryId uint      `json:"category_id" gorm:"not null"`            // 文章所属分类的 ID。
	Title      string    `json:"title" gorm:"type:varchar(50);not null"` // 文章标题，最大长度为 50。
	Content    string    `json:"content" gorm:"type:text;not null"`      // 文章内容。
	HeadImage  string    `json:"head_image"`                             // 文章头图的链接或路径。
	CreatedAt  Time      `json:"created_at" gorm:"type:timestamp"`       // 文章创建时间。
	UpdatedAt  Time      `json:"updated_at" gorm:"type:timestamp"`       // 文章更新时间。

}

// ArticleInfo 定义了用于传输的文章信息，可能是用于 API 响应。
type ArticleInfo struct {
	ID         string `json:"id"`          // 文章 ID，作为字符串传输。
	CategoryId uint   `json:"category_id"` // 文章所属分类的 ID。
	Title      string `json:"title"`       // 文章标题。
	Content    string `json:"content"`     // 文章内容。
	HeadImage  string `json:"head_image"`  // 文章头图的链接或路径。
	CreatedAt  Time   `json:"created_at"`  // 文章创建时间。

}

// BeforeCreate 是 GORM 的钩子方法，在创建文章之前自动调用。
// 此方法用于在创建文章记录之前生成并设置文章的 ID。
func (a *Article) BeforeCreate(s *gorm.Scope) error {
	// 使用 uuid.NewV4() 生成一个全新的随机 UUID 并设置为文章的 ID。
	return s.SetColumn("ID", uuid.NewV4())
}

// BeforCreate 在创建文章之前自动调用，用于初始化文章的评论。
