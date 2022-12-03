package service

import (
	"github.com/gin-gonic/gin"
	"nat-penetration/helper"
	"net/http"
)

func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	serverConf, err := helper.GetServerConf()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Server Conf Analyze Error: " + err.Error(),
		})
		return
	}
	if username != serverConf.Web.Username || password != serverConf.Web.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}
	// 生成 token
	token, err := helper.GenerateToken(username)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Token Generate Failed Error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Login Success",
		"data": gin.H{
			"token": token,
		},
	})
}
