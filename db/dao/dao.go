package dao

import (
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

const tableName = "Counters"

type CollectionInterfaceImpl struct {
	CollectionInterface
}

// ClearCounter 清除Counter
func (imp *CounterInterfaceImp) ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func (imp *CounterInterfaceImp) UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableName).Save(counter).Error
}

// GetCounter 查询Counter
func (imp *CounterInterfaceImp) GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableName).Where("id = ?", id).First(counter).Error

	return counter, err
}

func (dao *CollectionInterfaceImpl) GetBookCollections() ([]*model.BookmarkCollectionModel, error) {
	var collections []*model.BookmarkCollectionModel
	cli := db.Get()
	err := cli.Table("book_collection").Find(&collections).Error
	return collections, err
}

func (dao *CollectionInterfaceImpl) GetBookMarkItems(collectionID int) ([]*model.BookmarkItemModel, error) {
	var items []*model.BookmarkItemModel
	cli := db.Get()
	err := cli.Table("collectionItem").Where("collection_id = ?", collectionID).Find(&items).Error
	return items, err
}

func (dao *CollectionInterfaceImpl) DeleteBookMarkItem(resourceID string) error {
	cli := db.Get()
	return cli.Table("collectionItem").Where("resource_id = ?", resourceID).Delete(nil).Error
}

func (dao *CollectionInterfaceImpl) AddBookMarkItem(item *model.BookmarkItemModel) error {
	cli := db.Get()
	return cli.Table("collectionItem").Create(item).Error
}

func (dao *CollectionInterfaceImpl) AddBookMarkCollection(collection *model.BookmarkCollectionModel) error {
	cli := db.Get()
	return cli.Table("book_collection").Create(collection).Error
}

func (dao *CollectionInterfaceImpl) DeleteBookMarkCollection(collectionID int) error {
	cli := db.Get()
	return cli.Table("book_collection").Where("collection_id = ?", collectionID).Delete(nil).Error
}
