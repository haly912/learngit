// AuthMiddleware.go
package middleware

import (
	"blog_server/common"
	"blog_server/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware 是一个 Gin 中间件，用于验证请求中的 JWT Token。
func AuthMiddleware() gin.HandlerFunc {
	// 返回一个 Gin 的 HandlerFunc，用于中间件的实际处理。
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段。
		tokenString := c.Request.Header.Get("Authorization")

		// 如果 Authorization 字段为空，返回 401 状态码和错误信息。
		if tokenString == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort() // 终止后续处理器的执行。
			return
		}

		// 如果 Authorization 不合法（不包含 "Bearer" 前缀或长度不足），返回 401 状态码和错误信息。
		if tokenString == "" || len(tokenString) < 7 || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		// 提取 Authorization 字符串中的 token 部分，即 "Bearer" 后面的内容。
		tokenString = tokenString[7:]

		// 使用 common 包中的 ParseToken 函数解析 token。
		token, claims, err := common.ParseToken(tokenString)

		// 如果 token 解析失败或 token 无效，返回 401 状态码和错误信息。
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		// 从解析后的 claims 中获取 userId。
		userId := claims.UserId

		// 获取数据库连接实例。
		DB := common.GetDB()

		// 查询数据库，根据 userId 获取用户信息。
		var user model.User
		DB.Where("id = ?", userId).First(&user)

		// 将查询到的用户信息存储到 Gin 上下文中，以便后续处理函数可以访问。
		c.Set("user", user)

		// 执行后续的处理函数。
		c.Next()
	}
}
