package controller

import (
	"blog_server/common"
	"blog_server/model"
	"blog_server/response"
	"github.com/gin-gonic/gin"
)

// SearchCategory 查询分类
// controller/CategoryController.go
// SearchCategory 查询分类
// 此函数用于查询数据库中所有的分类信息。
func SearchCategory(c *gin.Context) {
	db := common.GetDB()            // 获取数据库连接
	var categories []model.Category // 定义一个切片用于存储查询到的分类信息
	// 使用Find方法查询所有分类信息，如果出错则返回错误信息
	if err := db.Find(&categories).Error; err != nil {
		response.Fail(c, nil, "查找失败") // 使用response包中的Fail函数返回错误信息
		return
	}
	// 如果查询成功，则使用response包中的Success函数返回分类信息和成功信息
	response.Success(c, gin.H{"categories": categories}, "查找成功")
}

// SearchCategoryName 查询分类名
// 此函数用于根据分类ID查询特定的分类名称。
func SearchCategoryName(c *gin.Context) {
	db := common.GetDB()        // 获取数据库连接
	var category model.Category // 定义一个结构体用于存储查询到的分类信息
	// 从请求的路径参数中获取分类ID
	categoryId := c.Params.ByName("id")
	// 使用Where和First方法根据分类ID查询分类信息，如果出错则返回错误信息
	if err := db.Where("id = ?", categoryId).First(&category).Error; err != nil {
		response.Fail(c, nil, "分类不存在") // 如果分类不存在则返回错误信息
		return
	}
	// 如果查询成功，则使用response包中的Success函数返回分类名称和成功信息
	response.Success(c, gin.H{"categoryName": category.CategoryName}, "查找成功")
}
