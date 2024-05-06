package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func userLogin(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 现在可以使用 loginReq.Username 和 loginReq.Password 获取 JSON 中的值
	fmt.Printf("Username: %s, Password: %s, Role: %sn", user.Username, user.Password, user.Role)

	// 进行登录验证逻辑...
	dbUser := new(User)
	if err := DB.Where("username = ?", user.Username).First(dbUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 60000, "message": "用户不存在"})
		return
	}
	if dbUser.CheckPassword(user.Password, user.PasswordHash) {
		// 密码验证成功，生成 token 并返回
		token := "admin-token" // 这里应该是生成的 token
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"token": token,
			},
		})
	} else {
		// 密码验证失败
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    60000,
			"message": "账号密码错误",
		})
	}
}

func userInfo(c *gin.Context) {
	token := c.Query("token")
	fmt.Println(token)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"roles":        []string{"admin"},
			"introduction": "Hello iam a admin",
			"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			"name":         "Super Admin YQ",
		},
	})
}

func userLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": "success",
	})
}

func userCreate(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 加密密码
	user.PasswordHash = user.SetPassword(user.Password)

	// 创建用户
	if err := user.CreateUser(DB,user); err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "detail": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": "User created successfully",
	})
}

func userDelete(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// 使用用户名来删除用户
	result := DB.Where("username = ?", username).Delete(&User{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "detail": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": "User deleted successfully",
	})
}