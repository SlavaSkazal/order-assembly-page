package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"orderAssembly/models"
	"strings"
)

type Database struct {
	dbSql *sql.DB
}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error: %w %w", errors.New("failed open db"), err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error: %w %w", errors.New("failed ping db"), err)
	}
	return &Database{dbSql: db}, nil
}

func (db *Database) CreateTables() error {
	err := db.makeQuery("CREATE TABLE IF NOT EXISTS products (product_id INTEGER PRIMARY KEY, name VARCHAR(200), type VARCHAR(100) CHECK( type IN ('phones','notebooks','wristwatchs', 'system units', 'microphones','tv')), category VARCHAR(100) CHECK( category IN ('electronics','other')), price INTEGER,vendor_code VARCHAR(50), quantity INTEGER);")
	if err != nil {
		return err
	}
	err = db.makeQuery(" CREATE TABLE IF NOT EXISTS orders (order_id INTEGER PRIMARY KEY, number INTEGER, date_creation DATE, date_delivery DATE, user_id INTEGER, user_name VARCHAR(200), sum INTEGER, done BOOLEAN DEFAULT 0);")
	if err != nil {
		return err
	}
	err = db.makeQuery("CREATE TABLE IF NOT EXISTS racks(rack_id  INTEGER PRIMARY KEY, name VARCHAR(100));")
	if err != nil {
		return err
	}
	err = db.makeQuery("CREATE TABLE IF NOT EXISTS orders_products (order_id INTEGER NOT NULL, product_id  INTEGER NOT NULL, product_name VARCHAR(200), price INTEGER, sum_product INTEGER, product_quantity INTEGER, PRIMARY KEY(order_id, product_id), CONSTRAINT fk_order_id FOREIGN KEY(order_id) REFERENCES orders(order_id), CONSTRAINT fk_product_id FOREIGN KEY(product_id) REFERENCES products(product_id));")
	if err != nil {
		return err
	}
	err = db.makeQuery("CREATE TABLE IF NOT EXISTS racks_products (rack_id INTEGER NOT NULL, product_id INTEGER NOT NULL, main_rack INTEGER, PRIMARY KEY(rack_id, product_id), CONSTRAINT fk_rack_id FOREIGN KEY(rack_id) REFERENCES racks(rack_id), CONSTRAINT fk_product_id FOREIGN KEY(product_id) REFERENCES products(product_id));")
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) CreateRecords() error {
	err := db.makeQuery("INSERT INTO products (name, type, category, price, vendor_code, quantity) VALUES ('Ноутбук Asus', 'notebooks', 'electronics', 30000, '126734', 10), ('Телевизор Philips', 'tv', 'electronics', 20000, '563334', 12), ('Телефон Samsung', 'phones', 'electronics', 10000, '566777', 10), ('Системный блок Sony', 'system units', 'electronics', 15000, '866734', 15), ('Часы Xiaomi', 'wristwatchs', 'electronics', 4000, '666734', 20), ('Микрофон JBH', 'microphones', 'electronics', 5000, '566234', 5);")
	if err != nil {
		return err
	}
	err = db.makeQuery("INSERT INTO orders (number, date_creation, user_id, user_name) VALUES (10, '2022-02-02', 4,'Nosov Anton'), (11, '2022-02-01', 3,'Asin Oleg'), (14, '2022-02-03', 4,'Nosov Anton'), (15, '2022-03-02', 2,'Kolesov Artem'), (16, '2022-04-11', 5,'Ivanov Oleg');")
	if err != nil {
		return err
	}
	err = db.makeQuery("INSERT INTO racks (name) VALUES ('А'), ('Б'), ('В'), ('Ж'), ('З'), ('А');")
	if err != nil {
		return err
	}
	err = db.makeQuery("INSERT INTO orders_products (order_id, product_id, price, product_quantity, sum_product, product_name) VALUES (1, 1, 30000, 2, 60000, 'Ноутбук Asus'), (2, 2, 20000, 3, 60000, 'Телевизор Philips'), (3, 1, 30000, 3, 90000, 'Ноутбук Asus'), (1, 3, 10000, 1, 10000, 'Телефон Samsung'), (3, 4, 15000, 4, 60000, 'Системный блок Sony'), (4, 5, 4000, 1, 4000, 'Часы Xiaomi'), (1, 6, 5000, 1, 5000, 'Микрофон JBH');")
	if err != nil {
		return err
	}
	err = db.makeQuery("INSERT INTO racks_products (rack_id, product_id , main_rack) VALUES (1, 1, 1), (1, 2, 1), (2, 3, 1), (4, 4, 1), (4, 5, 1), (4, 6, 1), (3, 3, 0), (5, 3, 0), (1, 5, 0);")
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) PrintAssemblyPage(orderNumbers []int) error {
	questionMarks := make([]string, len(orderNumbers))
	for i := range orderNumbers {
		questionMarks[i] = "?"
	}
	query := "SELECT orders_products.product_name, orders_products.product_id, orders_products.product_quantity, orders.number as order_number, orders.order_id, racks.name as rack_name, racks_products.main_rack, racks.rack_id FROM orders INNER JOIN orders_products ON orders.order_id = orders_products.order_id INNER JOIN racks_products ON orders_products.product_id = racks_products.product_id INNER JOIN racks ON racks_products.rack_id = racks.rack_id WHERE orders.number IN (" +
		strings.Join(questionMarks, ",") +
		");"

	args := make([]interface{}, len(orderNumbers))
	for i, v := range orderNumbers {
		args[i] = v
	}

	rows, err := db.dbSql.Query(query, args...)
	if err != nil {
		return err
	}
	printRows(rows, orderNumbers)
	return nil
}

func printRows(rows *sql.Rows, orderNumbers []int) {
	var productName, rackName string
	var productID, productQuantity, orderNumber, orderID, mainRack, rackID int
	racks := map[int]*models.Rack{}
	products := map[int]*models.Product{}
	orders := map[int]*models.Order{}

	for rows.Next() {
		rows.Scan(&productName, &productID, &productQuantity, &orderNumber, &orderID, &rackName, &mainRack, &rackID)
		_, ok := racks[rackID]
		if ok {
			racks[rackID].Products[productID] = true
		} else {
			rack := &models.Rack{rackID, rackName, map[int]bool{productID: true}}
			racks[rackID] = rack
		}

		_, ok = products[productID]
		if ok {
			if mainRack == 1 {
				products[productID].MainRack = rackID

			} else {
				products[productID].AdditionalRacks[rackID] = true
			}
			products[productID].Orders[orderID] = true
		} else {
			prod := &models.Product{ID: productID, Name: productName, Orders: map[int]bool{orderID: true}, AdditionalRacks: map[int]bool{}}
			if mainRack == 1 {
				prod.MainRack = rackID
			} else {
				prod.AdditionalRacks[rackID] = true
			}
			products[productID] = prod
		}

		_, ok = orders[orderID]
		if ok {
			orders[orderID].ProdQauntyty[productID] = productQuantity
		} else {
			ord := &models.Order{orderID, orderNumber, map[int]int{productID: productQuantity}}
			orders[orderID] = ord
		}
	}

	fmt.Print("=+=+=+=\nСтраница сборки заказов ")
	for i := 0; i < len(orderNumbers); i++ {
		fmt.Print(orderNumbers[i])
		if i+1 != len(orderNumbers) {
			fmt.Print(",")
		}
	}
	fmt.Print("\n\n")

	for _, rack := range racks {
		rackAlreadyPrint := false
		for prod := range rack.Products {
			if products[prod].MainRack == rack.ID {
				if !rackAlreadyPrint {
					fmt.Printf("===Стеллаж %s\n", rack.Name)
					rackAlreadyPrint = true
				}

				adRacks := make([]string, len(products[prod].AdditionalRacks))
				i := 0
				for adRack := range products[prod].AdditionalRacks {
					adRacks[i] = racks[adRack].Name
					i++
				}

				for ord := range products[prod].Orders {
					fmt.Printf("%s (id=%d)\n", products[prod].Name, prod)
					fmt.Printf("заказ %d, %d шт\n", orders[ord].Number, orders[ord].ProdQauntyty[prod])
					if len(adRacks) > 0 {
						fmt.Printf("доп стеллаж: %s\n", strings.Join(adRacks, ","))
					}
					fmt.Print("\n")
				}
			}
		}
	}
}

func (db *Database) makeQuery(strQuery string) error {
	stat, err := db.dbSql.Prepare(strQuery)
	if err != nil {
		return err
	}
	_, err = stat.Exec()
	if err != nil {
		return err
	}
	return nil
}
