package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// Upload 上传图像
// FileController.go

// Upload 函数用于处理图像上传的请求。
func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		// 如果获取文件失败，返回服务器内部错误信息。
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "格式错误",
		})
		return
	}

	// 获取文件名和扩展名。
	filename := header.Filename
	ext := path.Ext(filename)

	// 使用当前时间作为新文件名的一部分，以避免文件名冲突。
	name := "image_" + time.Now().Format("20060102150405")
	// 构建新的文件名，包括扩展名。
	newFilename := name + ext

	// 创建保存文件的路径，并打开文件准备写入。
	out, err := os.Create("static/images/" + newFilename)
	if err != nil {
		// 如果创建文件失败，返回服务器内部错误信息。
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建错误",
		})
		return
	}
	defer out.Close() // 确保在函数结束时关闭文件。

	// 将上传的文件内容复制到新创建的文件中。
	_, err = io.Copy(out, file)
	if err != nil {
		// 如果文件复制失败，返回服务器内部错误信息。
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "复制错误",
		})
		return
	}

	// 如果上传成功，返回状态码200和文件路径。
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		// 返回的文件路径是相对于静态资源目录的路径。
		"data": gin.H{"filePath": "/images/" + newFilename},
		"msg":  "上传成功",
	})
}

// RichEditorUpload 上传富文本编辑器中的图像
func RichEditorUpload(c *gin.Context) {
	fromData, _ := c.MultipartForm()
	files := fromData.File["wangeditor-uploaded-image"]
	var url []string
	for _, file := range files {
		ext := path.Ext(file.Filename)
		name := "image_" + time.Now().Format("20060102150405")
		newFilename := name + ext
		dst := path.Join("./static/images", newFilename)
		fileurl := "/images/" + newFilename
		url = append(url, fileurl)
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"errno":   1,
				"message": "上传失败",
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"errno": 0,
		"data": gin.H{
			"url": url[0],
		},
	})
}
