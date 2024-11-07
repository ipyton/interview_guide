package dao

import (
	"errors"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

type CollectionInterfaceImpl struct {
	CollectionInterface
}

//var Imp CounterInterface = &CounterInterfaceImp{}

func (dao *CollectionInterfaceImpl) IsResourceCollected(userId string, resourceId int) (bool, error) {
	return true, errors.New("not implement")
}

func (dao *CollectionInterfaceImpl) GetCollections() (*[]model.BookmarkCollectionModel, error) {
	var collections *[]model.BookmarkCollectionModel
	cli := db.Get()
	err := cli.Table("book_collection").Find(&collections).Error
	return collections, err
}

func (dao *CollectionInterfaceImpl) GetItemsInCollection(userId string, collectionID int) ([]*model.BookmarkItemModel, error) {
	var items []*model.BookmarkItemModel
	cli := db.Get()
	err := cli.Table("collectionItem").Where("collection_id = ?", collectionID).Find(&items).Error
	return items, err
}

func (dao *CollectionInterfaceImpl) DeleteBookMarkItem(userId string, collectionId int, questionId int) error {
	cli := db.Get()
	return cli.Table("collectionItem").Where("resource_id = ?", questionId).Delete(nil).Error
}

func (dao *CollectionInterfaceImpl) AddBookMarkItem(item *model.BookmarkItemModel) error {
	cli := db.Get()
	return cli.Table("collectionItem").Create(item).Error
}

func (dao *CollectionInterfaceImpl) AddBookMarkCollection(collection *model.BookmarkCollectionModel) error {
	cli := db.Get()
	return cli.Table("book_collection").Create(collection).Error
}

func (dao *CollectionInterfaceImpl) DeleteBookMarkCollection(user_id string, collectionID int) error {
	cli := db.Get()
	return cli.Table("book_collection").Where("collection_id = ?", collectionID).Delete(nil).Error
}
