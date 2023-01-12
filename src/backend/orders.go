package backend

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Total        int         `json:"total"`
	Status       string      `json:"status"`
	Items        []orderItem `json:"items"`
}

type orderItem struct {
	OrderId   int `json:"order_id"`
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func getOrders(db *sql.DB) (orders []order, err error) {
	rows, err := db.Query("select ID,customerName,total,status from orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ID, &o.CustomerName, &o.Total, &o.Status); err != nil {
			return nil, err
		}
		o.getOrderItems(db)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return
}

func (o *order) getOrder(db *sql.DB) error {
	db.QueryRow("select ID,customerName,total,status from orders where id=?", o.ID).
		Scan(&o.ID, &o.CustomerName, &o.Total, &o.Status)
	err := o.getOrderItems(db)
	if err != nil {
		return err
	}
	return nil
}

func (o *order) getOrderItems(db *sql.DB) error {
	rows, err := db.Query("Select * from order_items where order_id=?", o.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	orderItems := []orderItem{}
	for rows.Next() {
		var oi orderItem
		if err := rows.Scan(&oi.OrderId, &oi.ProductId, &oi.Quantity); err != nil {
			return err
		}
		orderItems = append(orderItems, oi)
	}
	o.Items = orderItems
	return nil
}

func (o *order) createOrder(db *sql.DB) error {
	res, err := db.Exec(
		"INSERT INTO Orders(customerName,total,status) VALUES(?,?,?)",
		o.CustomerName,
		o.Total,
		o.Status,
	)
	if err != nil {
		fmt.Println("insert err")
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	o.ID = int(id)

	if err != nil {
		return err
	}
	return nil
}

func (oi *orderItem) createOrderItem(db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO order_items (order_id,product_id,quantity) VALUES(?,?,?)",
		oi.OrderId,
		oi.ProductId,
		oi.Quantity,
	)
	if err != nil {
		return err
	}
	return nil
}
