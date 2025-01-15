package accdb

import (
	"gorm.io/gorm"
)

type Item struct {
	Propid int    `gorm:"column:propid;primaryKey;autoIncrement"`
	Accid  int    `gorm:"column:accid;not null"`
	Prop   string `gorm:"column:prop;not null"`
	Value  string `gorm:"column:value;not null"`
}

const TblNameItems = "t_items"

func CreateItemsTable(db *gorm.DB) error {
	return db.Table(TblNameItems).AutoMigrate(&Item{})
}

func DropItemsTable(db *gorm.DB) error {
	return db.Migrator().DropTable(TblNameItems)
}

func GetItemsByAccID(db *gorm.DB, accid int) ([]Item, error) {
	var items []Item
	err := db.Table(TblNameItems).Where("accid = ?", accid).Find(&items).Error
	return items, err
}

func CountItemsByAccID(db *gorm.DB, accid int) (int64, error) {
	var count int64
	err := db.Table(TblNameItems).Where("accid = ?", accid).Count(&count).Error
	return count, err
}

func CreateItem(db *gorm.DB, item *Item) error {
	return db.Table(TblNameItems).Create(item).Error
}

func DeleteItemByPropID(db *gorm.DB, propid int) error {
	return db.Table(TblNameItems).Where("propid = ?", propid).Delete(&Item{}).Error
}

func DeleteItemsByAccID(db *gorm.DB, accid int) error {
	return db.Table(TblNameItems).Where("accid = ?", accid).Delete(&Item{}).Error
}

func UpdateItem(db *gorm.DB, item *Item) error {
	return db.Table(TblNameItems).Where("propid = ?", item.Propid).Updates(item).Error
}
