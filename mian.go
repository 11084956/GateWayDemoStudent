package main

import (
	"github.com/gin-gonic/gin"
	"image/color"
	"os"
)
import qrcode "github.com/skip2/go-qrcode"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/delete", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"form_id": c.PostForm("form_id"),
			"id":      c.Query("id"),
		})
	})

	r.GET("qr_code", func(ctx *gin.Context) {
		var png []byte
		png, err := qrcode.Encode("https://img-blog.csdnimg.cn/2019122612311680.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3poYW5neHc4NzIxOTY=,size_16,color_FFFFFF,t_70", qrcode.Medium, 1080)
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg": "生成二维码失败",
				"err": err.Error(),
			})

			return
		}

		//w.Header().Set("content-type","application/json")
		ctx.Header("content-type", "image/png")
		ctx.JSON(200, gin.H{
			"msg": "success",
			"png": png,
		})
	})

	r.GET("qr_code_file", func(ctx *gin.Context) {
		isFile("./1214.png")
		err := qrcode.WriteColorFile("https://cupcake.nilssonlee.se/wp-content/uploads/2020/07/IMG_1432-scaled.jpg", qrcode.Highest, 256, color.Black, color.White, "1214.png")
		if err != nil {
			ctx.JSON(200, gin.H{
				"msg": "二维码创建失败",
				"err": err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"msg": "success",
			"img": nil,
		})
	})

	_ = r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func isFile(fileName string) {
	file, err := os.Open("./" + fileName)

	if err != nil && os.IsNotExist(err) {
		file, _ = os.Create("./" + fileName)
	}

	defer func() { _ = file.Close() }()
}
