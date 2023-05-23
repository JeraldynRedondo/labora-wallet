package service

import "my_api_project/model"

type ItemService struct {
	DbHandler model.DBHandler
}

func (s *ItemService) GetItems() ([]model.Item, error) {
	return s.DbHandler.GetItems()
}

func (s *ItemService) GetItemsPerPage(pages, itemsPerPage int) ([]model.Item, int, error) {
	return s.DbHandler.GetItemsPerPage(pages, itemsPerPage)
}

func (s *ItemService) GetItemId(id int) (model.Item, error) {
	return s.DbHandler.GetItemId(id)
}

func (s *ItemService) GetItemName(name string) ([]model.Item, error) {
	return s.DbHandler.GetItemName(name)
}

func (s *ItemService) CreateItem(newItem model.Item) (model.Item, error) {
	return s.DbHandler.CreateItem(newItem)
}

func (s *ItemService) UpdateItem(id int, item model.Item) (model.Item, error) {
	return s.DbHandler.UpdateItem(id, item)
}

func (s *ItemService) DeleteItem(id int) error {
	return s.DbHandler.DeleteItem(id)
}

/*
func (s *ItemService) UpdateItemDetails(id int) (model.Item, error) {
	return s.DbHandler.UpdateItemDetails(id)
}*/
