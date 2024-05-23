package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
)

func dynadotDomainList(c *gin.Context) {
	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=list_domain"

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	res := gjson.Get(string(body), "ListDomainInfoResponse.Status")
	if res.String() == "success" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res": "success",
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"data": gin.H{
				"res": "failed",
			},
		})
	}

}

func dynadotSearchDomain(c *gin.Context) {

	domain := c.PostForm("domain")

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=search&domain0=" + domain + "&show_price=1&currency=USD"
	println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	res := gjson.Get(string(body), "SearchResponse.SearchResults.0.Available")
	fmt.Println(res)
	if res.String() == "yes" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"message": "域名可被注册",
				"res":     "success",
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"data": gin.H{
				"res":     "failed",
				"message": "域名不可被注册",
			},
		})
	}

}

func dynadotBuyDomain(c *gin.Context) {
	domain := c.PostForm("domain")

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=register&domain=" + domain + "&duration=1&currency=USD"
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	res := gjson.Get(string(body), "RegisterResponse.Status")
	fmt.Println(res)

	if res.String() == "success" {
		c.JSON(http.StatusOK, gin.H{
            "code": 20000,
            "data": gin.H{
                "message": "域名购买成功",
                "res":     "success",
            },
        })
	}else {
		c.JSON(http.StatusOK, gin.H{
            "code": 40001,
            "data": gin.H{
                "message": "域名购买失败",
                "res":     "failed",
				"reason": res.String(),
            },
        })
	}
}
