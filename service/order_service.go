package service

import (
	"Assignment-2/db"
	"fmt"
	"log"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Item struct {
	ItemId      int    `json:"itemId"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	OrderId     int    `json:"orderId"`
}

type Order struct {
	OrderId      int       `json:"orderId"`
	OrderedAt    time.Time `json:"orderedAt"`
	CustomerName string    `json:"customerName"`
	Items        []Item    `json:"items"`
}

const (
	CreateOrder  	 = `INSERT INTO orders (customer_name, ordered_at) VALUES($1, $2) RETURNING order_id, customer_name`
	CreateItem   	 = `INSERT INTO items (item_code, description, quantity, order_id) VALUES($1, $2, $3, $4) RETURNING item_id, item_code, description, quantity, order_id`
	GetAllOrder   	 = `SELECT * FROM orders`
	GetItemByOrderId = `SELECT * FROM items WHERE order_id = $1`
	GetOrderById     = `SELECT order_id, customer_name, ordered_at FROM orders WHERE order_id = $1`
	UpdateOrderById  = `UPDATE orders SET customer_name = $1, ordered_at = $2 WHERE order_id=$3 RETURNING order_id`
	UpdateItemById   = `UPDATE items SET item_code = $1, description = $2, quantity = $3, order_id = $4 WHERE item_id = $5 RETURNING item_code, description, quantity, order_id, item_id`
	DeleteOrder      = `DELETE FROM orders WHERE order_id = $1`
	DeleteItem       = `DELETE FROM items WHERE order_id = $1`
)

var OrderService orderService = &orderRepo{}

type orderService interface {
	CreateOrder(*Order) *Order
	DeleteOrder(int) string
	UpdateOrder(*Order) *Order
	GetOrder() *[]Order
}

type orderRepo struct{}

func (m *Order) GetOrderParamId(c *gin.Context) int64 {
	paramId := c.Param("id")
	fmt.Println("Ini adalah order id =>", paramId)
	ID, err := strconv.Atoi(paramId)

	if err != nil {
		return 0
	}

	return int64(ID)
}

func (m *orderRepo) CreateOrder(orderReq *Order) *Order {
	db := db.GetDB()
	row := db.QueryRow(CreateOrder, orderReq.CustomerName, orderReq.OrderedAt)
	var orderID = 0
	var customer_name = ""
	err := row.Scan(&orderID, &customer_name)
	if err != nil {
		log.Fatal(err)
	}

	items := []Item{}
	for _, itemReq := range orderReq.Items {
		row = db.QueryRow(CreateItem, itemReq.ItemCode, itemReq.Description, itemReq.Quantity, orderID)
		var itemRes Item
		err = row.Scan(&itemRes.ItemId, &itemRes.ItemCode, &itemRes.Description, &itemRes.Quantity, &itemRes.OrderId)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, itemRes)
	}

	orderReq.Items = items
	orderReq.OrderId = orderID
	orderReq.CustomerName = customer_name

	return orderReq
}

func (m *orderRepo) GetOrder() *[]Order {
	db := db.GetDB()
	row, err := db.Query(GetAllOrder)
	if err != nil {
		log.Fatal("fail to get data")
	}
	var dataOrder []Order

	for row.Next() {
		var ord Order
		if err := row.Scan(&ord.OrderId, &ord.CustomerName, &ord.OrderedAt); err != nil {
			log.Fatal("err")
		}

		row2, err := db.Query(GetItemByOrderId, &ord.OrderId)
		if err != nil {
			log.Fatal("Gagal query item")
		}

		var dataItem []Item
		for row2.Next() {
			var itm Item
			if err := row2.Scan(&itm.ItemId, &itm.ItemCode, &itm.Description, &itm.Quantity, &itm.OrderId); err != nil {
				log.Fatal(err)
			}
			dataItem = append(dataItem, itm)

		}
		ord.Items = append(ord.Items, dataItem...)
		dataOrder = append(dataOrder, ord)
	}

	if err = row.Err(); err != nil {
		log.Fatal("Error ")
	}

	return &dataOrder
}

func (m *orderRepo) UpdateOrder(orderReq *Order) *Order {
	db := db.GetDB()
	fmt.Println(orderReq)
	row := db.QueryRow(UpdateOrderById, orderReq.CustomerName, orderReq.OrderedAt, orderReq.OrderId)

	var orderdata Order
	err := row.Scan(&orderdata.OrderId)
	if err != nil {
		log.Fatal(err)
	}

	var resdata Order
	row2 := db.QueryRow(GetOrderById, orderReq.OrderId)
	if err := row2.Scan(&resdata.OrderId, &resdata.CustomerName, &resdata.OrderedAt); err != nil {
		log.Fatal(err)
	}

	items := []Item{}
	for _, itemReq := range orderReq.Items {
		row3 := db.QueryRow(UpdateItemById, itemReq.ItemCode, itemReq.Description, itemReq.Quantity, orderReq.OrderId, itemReq.ItemId)
		var itemRes Item
		err = row3.Scan(&itemRes.ItemCode, &itemRes.Description, &itemRes.Quantity, &itemRes.OrderId, &itemRes.ItemId)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, itemRes)
	}
	resdata.Items = append(resdata.Items, items...)
	return &resdata
}

func (m *orderRepo) DeleteOrder(data int) string {
	db := db.GetDB()
	var datastring string
	row, err := db.Exec(DeleteOrder, data)

	if err == nil {
		count, err := row.RowsAffected()
		if err == nil {
			datastring = fmt.Sprintf("%d row order affected. ", count)
		}

		row2, err := db.Exec(DeleteItem, data)
		if err == nil {
			count2, err := row2.RowsAffected()
			if err == nil {
				datastring = datastring + fmt.Sprintf("%d row item affected.", count2)
			}
		}
	}
	return datastring
}
