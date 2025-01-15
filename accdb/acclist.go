package accdb

import (
	"gorm.io/gorm"
)

type AccountInfo struct {
	// gorm.Model
	AccID   int    `gorm:"column:acc_id;primaryKey;autoIncrement"`
	AccName string `gorm:"column:acc_name;size:255;not null"`
}

const TblNameAccList = "t_acclist"

func CreateAccListTable(db *gorm.DB) error {
	return db.Table(TblNameAccList).AutoMigrate(&AccountInfo{})
}

func DropAccListTable(db *gorm.DB) error {
	return db.Migrator().DropTable(TblNameAccList)
}

func GetAccountByID(db *gorm.DB, id int) (*AccountInfo, error) {
	var account AccountInfo
	err := db.Table(TblNameAccList).Where("acc_id = ?", id).First(&account).Error
	return &account, err
}

func SearchAccountsByName(db *gorm.DB, name string) ([]AccountInfo, error) {
	var accounts []AccountInfo
	err := db.Table(TblNameAccList).Where("acc_name LIKE ?", "%"+name+"%").Find(&accounts).Error
	return accounts, err
}

func CreateAccount(db *gorm.DB, account *AccountInfo) error {
	return db.Table(TblNameAccList).Create(account).Error
}

func DeleteAccount(db *gorm.DB, id int) error {
	return db.Table(TblNameAccList).Where("acc_id = ?", id).Delete(&AccountInfo{}).Error
}

func UpdateAccount(db *gorm.DB, account *AccountInfo) error {
	return db.Table(TblNameAccList).Where("acc_id = ?", account.AccID).Updates(account).Error
}

/*
select * from t_acclist where acc_id=?
select * from "TBNAME_ACCLIST
select * from t_acclist where acc_name like '%'||?||'%'
insert into t_acclist(acc_name) values(:acc_name)
delete from t_acclist where acc_id=?
update t_acclist set acc_name=? where acc_id=?
*/
