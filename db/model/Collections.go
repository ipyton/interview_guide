package model

type BookmarkCollectionModel struct {
	UserId         string `gorm:"column:user_id" json:"user_id"`
	CollectionID   int    `gorm:"column:collection_id" json:"collection_id"`
	CollectionName string `gorm:"column:collection_name" json:"collection_name"`
	Count          int    `gorm:"column:count" json:"count"`
	Extra1         string `gorm:"column:extra1;omitempty" json:"extra1,omitempty"` // 额外字段1
	Extra2         string `gorm:"column:extra2;omitempty" json:"extra2,omitempty"` // 额外字段2
	Extra3         string `gorm:"column:extra3;omitempty" json:"extra3,omitempty"` // 额外字段3
}

type BookmarkItemModel struct {
	CollectionID  int    `gorm:"column:collection_id" json:"collection_id"`
	ResourceID    int    `gorm:"column:resource_id" json:"resource_id"`
	Type          string `gorm:"column:type" json:"type"`
	ResourceTitle string `gorm:"column:resource_title" json:"resource_title"`
	Content       string `gorm:"column:content" json:"content"`
	Extra1        string `gorm:"column:extra1;omitempty" json:"extra1,omitempty"` // 额外字段1
	Extra2        string `gorm:"column:extra2;omitempty" json:"extra2,omitempty"` // 额外字段2
	Extra3        string `gorm:"column:extra3;omitempty" json:"extra3,omitempty"` // 额外字段3
}
