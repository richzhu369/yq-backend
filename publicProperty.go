package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

func getPublicProperty(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 20000,
		"data": gin.H{
			"items": publicProperty,
		},
	})
}

// 这里改成同时传来字段和值，然后根据字段名称，更新值
func editPublicProperty(c *gin.Context) {
	//var publicProperty PublicProperty

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 20000, "message": err.Error()})
		return
	}
	key := gjson.GetBytes(jsonData, "key").String()
	value := gjson.GetBytes(jsonData, "newValue").String()
	fmt.Println(key, value)

	// 查找出来 value的值是哪个字段
	// 使用反射来确定字段类型并更新值
	elem := reflect.ValueOf(&publicProperty).Elem()
	fieldVal := elem.FieldByName(key)

	if fieldVal.IsValid() && fieldVal.CanSet() {
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intValue, err := strconv.ParseInt(value, 10, fieldVal.Type().Bits()); err == nil {
				fieldVal.SetInt(intValue)
			} else {
				log.Println(err.Error())
				c.JSON(http.StatusOK, gin.H{"code": 40001, "message": "Invalid integer value"})
				return
			}
		default:
			c.JSON(http.StatusOK, gin.H{"code": 40001, "message": "Unsupported field type"})
			return
		}
		// 更新数据库中的字段
		if err := DB.Model(&publicProperty).Where("id = ?", 1).Update(key, fieldVal.Interface()).Error; err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, gin.H{"code": 40001, "message": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 40001, "message": "Field not found or cannot be set"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "修改值成功",
	})
}
