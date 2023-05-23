package service

import (
	"fmt"
	"io/ioutil"
	"my_api_project/model"
	"net/http"
	"sync"
	"time"
)

// GetItems it is a function that makes a query and returns all the items in the database.
func (Db *PostgresDBHandler) GetItems() ([]model.Item, error) {
	items := make([]model.Item, 0)
	query := "SELECT * FROM items ORDER BY id"
	rows, err := Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
		if err != nil {
			fmt.Printf("Error extracting item: %v", err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// GetItemsPerPage it is a function that queries a database and returns a number of items per page.
func (Db *PostgresDBHandler) GetItemsPerPage(pages, itemsPerPage int) ([]model.Item, int, error) {

	//Calculate the initial index and item limit based on the current page and items per page.
	start := (pages - 1) * itemsPerPage

	//Get the total number of rows in the items table
	var count int
	query := "SELECT COUNT(*) FROM items"
	err := Db.QueryRow(query).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("Error querying the count in database: %w", err)
	}

	// Get the list of elements corresponding to the current page
	query = "SELECT * FROM items ORDER BY id OFFSET $1 LIMIT $2"
	rows, err := Db.Query(query, start, itemsPerPage)
	if err != nil {
		return nil, 0, fmt.Errorf("Error querying database: %w", err)
	}

	defer rows.Close()

	var newListItems []model.Item

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
		if err != nil {
			return nil, 0, fmt.Errorf("Error extracting item: %w", err)
		}
		newListItems = append(newListItems, item)
	}

	if len(newListItems) == 0 {
		return nil, 0, fmt.Errorf("No items found for page %d", pages)
	}
	return newListItems, count, nil
}

var m sync.Mutex

// GetItemId it is a function that performs a query by Item id and returns the item that matches.
func (Db *PostgresDBHandler) GetItemId(id int) (model.Item, error) {
	var item model.Item

	//Increase the count of views of the item
	m.Lock()
	Db.QueryRow("UPDATE items SET view_count = view_count + $1 WHERE id = $2 RETURNING *",
		1, id)
	m.Unlock()

	//Get item.
	query := "SELECT * FROM items WHERE id=$1"

	err := Db.QueryRow(query, id).Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
	if err != nil {
		return model.Item{}, fmt.Errorf("Error querying database: %w", err)
	}

	return item, nil
}

// GetItemName it is a function that performs a query by Item Name and returns the items that match.
func (Db *PostgresDBHandler) GetItemName(name string) ([]model.Item, error) {
	var items []model.Item
	var item model.Item

	query := "SELECT * FROM items WHERE customer_name ILIKE $1"
	rows, err := Db.Query(query, "%"+name+"%")
	if err != nil {
		return items, fmt.Errorf("Error querying database: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
		if err != nil {
			return items, fmt.Errorf("Error extracting item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// CreateItem is a function that creates an Item in the database.
func (Db *PostgresDBHandler) CreateItem(newItem model.Item) (model.Item, error) {

	details, err := getDetails()
	if err != nil {
		return model.Item{}, fmt.Errorf("Error getting details: %w", err)
	}

	newItem.TotalPrice = newItem.CalculatedTotalPrice()
	// Insertar el nuevo item en la base de datos
	query := `INSERT INTO items (customer_name, order_date, product, quantity, price,details,total_price ,view_count)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	row := Db.QueryRow(query, newItem.Customer_name, time.Now(), newItem.Product, newItem.Quantity, newItem.Price, details, newItem.TotalPrice, 0)

	err = row.Scan(&newItem.ID, &newItem.Customer_name, &newItem.Order_date, &newItem.Product, &newItem.Quantity, &newItem.Price, &newItem.Details, &newItem.TotalPrice, &newItem.ViewCount)
	if err != nil {
		return model.Item{}, fmt.Errorf("Error extracting item: %w", err)
	}

	return newItem, nil
}

// UpdateItem it is a function that updates an item by id.
func (Db *PostgresDBHandler) UpdateItem(id int, item model.Item) (model.Item, error) {
	var updatedItem model.Item

	details, err := getDetails()
	if err != nil {
		return model.Item{}, fmt.Errorf("Error getting details: %w", err)
	}
	item.TotalPrice = item.CalculatedTotalPrice()

	query := "UPDATE items SET customer_name = $1, order_date = $2, product = $3, quantity = $4, price = $5, details = $6, total_price=$7 WHERE id = $8 RETURNING *"
	row := Db.QueryRow(query, item.Customer_name, time.Now(), item.Product, item.Quantity, item.Price, details, item.TotalPrice, id)
	err = row.Scan(&updatedItem.ID, &updatedItem.Customer_name, &updatedItem.Order_date, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.Details, &updatedItem.TotalPrice, &updatedItem.ViewCount)
	if err != nil {
		return model.Item{}, fmt.Errorf("Error extracting item: %w", err)
	}

	return updatedItem, nil
}

// DeleteItem it is a function that updates an item by id.
func (Db *PostgresDBHandler) DeleteItem(id int) error {

	query := "DELETE FROM items WHERE id = ?"
	_, err := Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error querying database: %w", err)
	}

	return nil
}

// getDetails it is a function that consumes a service that returns a string to get the details of an item.
func getDetails() (string, error) {
	// Realizamos la petición a la API de loripsum
	url := fmt.Sprintf("http://loripsum.net/api/1/short")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Error getting details of the service: %w", err)
	}

	defer resp.Body.Close()

	// Read the API response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error in the body of the response of the service: %w", err)
	}

	data := string(body)

	return data, nil
}

/*
// UpdateItemDetails it is a function that updates the Details of an item.
func (Db *PostgresDBHandler) UpdateItemDetails(id int) (model.Item, error) {
	// Get the paragraph of the object from API
	detail, err := getDetails()
	if err != nil {
		fmt.Println(err)
		return model.Item{}, fmt.Errorf("Error getting details: %w", err)
	}

	// Update the "details" column in the "items" table
	var updatedItem model.Item
	query := "UPDATE items SET details=$1 WHERE id=$2 RETURNING *"

	row := Db.QueryRow(query, detail, id)
	err = row.Scan(&updatedItem.ID, &updatedItem.Customer_name, &updatedItem.Order_date, &updatedItem.Product, &updatedItem.Quantity, &updatedItem.Price, &updatedItem.Details, &updatedItem.TotalPrice, &updatedItem.ViewCount)
	if err != nil {
		fmt.Println(err)
		return model.Item{}, fmt.Errorf("Error extracting item: %w", err)
	}

	return updatedItem, nil
}


// UpdateTotalPriceItem it is a function that updates the total prices of the items
func UpdateTotalPriceItem(item model.Item) (model.Item, error) {

		totalPrice := item.CalculatedTotalPrice()
		query := "UPDATE items SET total_price=$1 WHERE id = $2 RETURNING *"
		row := Db.QueryRow(query, totalPrice, item.ID)
		err := row.Scan(&item.ID, &item.Customer_name, &item.Order_date, &item.Product, &item.Quantity, &item.Price, &item.Details, &item.TotalPrice, &item.ViewCount)
		if err != nil {
			return model.Item{}, fmt.Errorf("Error extracting item: %w", err)
		}

		return item, nil
	}

// UpdateViewCount it is a function that updates the view count of an item.
func UpdateViewCount(id int) {
	m.Lock()
	Db.QueryRow("UPDATE items SET view_count = view_count + $1 WHERE id = $2 RETURNING *",
		1, id)
	m.Unlock()
}
*/
