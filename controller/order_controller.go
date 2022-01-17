package controller

import (
	"Assignment-2/service"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	fmt.Println(data)
	var orderedAtStr string
	_ = orderedAtStr
	if val, ok := data["orderedAt"].(string); ok {
		orderedAtStr = val
	}

	y, err := time.Parse("2006-01-02", orderedAtStr)
	if err != nil {
		fmt.Println(err)
	}

	items := []service.Item{}
	for _, v := range data["items"].([]interface{}) {
		item := v.(map[string]interface{})
		itemReq := []service.Item{
			{
				ItemCode:    item["itemCode"].(string),
				Description: item["description"].(string),
				Quantity:    int(item["quantity"].(float64)),
			},
		}
		items = append(items, itemReq...)
	}
	req := service.Order{
		CustomerName: data["customerName"].(string),
		OrderedAt:    y,
		Items:        items,
	}
	_ = req

	res := service.OrderService.CreateOrder(&req)
	c.JSON(201, res)

}

func GetOrder(c *gin.Context) {
	res := service.OrderService.GetOrder()
	c.JSON(201, res)
}

func UpdateOrder(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"errMsg": err.Error(),
		})
		return
	}
	fmt.Println(data)

	var orderedAtStr string
	_ = orderedAtStr
	if val, ok := data["orderedAt"].(string); ok {
		orderedAtStr = val
	}

	y, err := time.Parse("2006-01-02", orderedAtStr)
	if err != nil {
		fmt.Println(err)
	}

	items := []service.Item{}
	for _, v := range data["items"].([]interface{}) {
		item := v.(map[string]interface{})
		itemReq := []service.Item{
			{
				ItemId:      int(item["itemId"].(float64)),
				ItemCode:    item["itemCode"].(string),
				Description: item["description"].(string),
				Quantity:    int(item["quantity"].(float64)),
			},
		}
		items = append(items, itemReq...)
	}
	req := service.Order{
		CustomerName: data["customerName"].(string),
		OrderedAt:    y,
		Items:        items,
	}
	_ = req

	orderid := req.GetOrderParamId(c)

	req.OrderId = int(orderid)

	fmt.Println("item : ", req.Items)
	res := service.OrderService.UpdateOrder(&req)
	c.JSON(201, res)
}

func DeleteOrder(c *gin.Context) {
	var order service.Order
	orderid := order.GetOrderParamId(c)
	fmt.Println(orderid)
	res := service.OrderService.DeleteOrder(int(orderid))
	fmt.Println(res)
	c.String(201, res)
}
