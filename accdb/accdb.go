package accdb

import (
	"accgo/utils/logs"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AccDB struct {
	db *gorm.DB
}

type PropKV struct {
	PropID int
	Key    string
	Val    string
}

func Start(dbpath string) (*AccDB, error) {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		logs.Fatal(err).Str("db", dbpath).Msg("open fail")
		return nil, err
	}

	return &AccDB{db}, nil
}

func (ad *AccDB) Init() error {
	if err := CreateAccListTable(ad.db); err != nil {
		logs.Error(err).Msg("CreateAccListTable fail")
		return err
	}
	if err := CreateItemsTable(ad.db); err != nil {
		logs.Error(err).Msg("CreateItemsTable fail")
		return err
	}
	return nil
}

func (ad *AccDB) NewAcc(accName string, props ...PropKV) error {
	acc := &AccountInfo{
		AccName: accName,
	}
	if err := CreateAccount(ad.db, acc); err != nil {
		logs.Error(err).Str("name", accName).Msg("CreateAccount fail")
		return err
	}
	logs.Debug().Str("name", accName).Int("accid", acc.AccID).Msg("CreateAccount succ")
	return ad.AddProps(acc.AccID, props...)
}

func (ad *AccDB) AddProps(accid int, props ...PropKV) error {
	var rtnErr error = nil

	item := &Item{
		Accid: accid,
	}

	for _, prop := range props {
		item.Propid = 0
		item.Prop = prop.Key
		item.Value = prop.Val
		if err := CreateItem(ad.db, item); err != nil {
			rtnErr = err
			logs.Error(err).Int("accid", accid).Str("key", item.Prop).Str("val", item.Value).Msg("CreateItem fail")
		} else {
			logs.Debug().Int("accid", accid).Str("key", item.Prop).Str("val", item.Value).Int("propid", item.Propid).Msg("CreateItem succ")
		}
	}
	return rtnErr
}

func (ad *AccDB) DelAcc(accid int, propids ...int) error {
	if len(propids) > 0 { // 传入了propid, 说明只是删除一个kv
		var rtnErr error = nil
		for _, propid := range propids {
			if err := DeleteItemByPropID(ad.db, propid); err != nil {
				rtnErr = err
				logs.Error(err).Int("accid", accid).Int("propid", propid).Msg("DeleteItemByPropID fail")
			}
		}
		return rtnErr
	} // else， 没有传 propid， 只传了accid，说明要删除指定acc 及其下所有的 prop

	if err := DeleteAccount(ad.db, accid); err != nil {
		logs.Error(err).Int("accid", accid).Msg("DeleteAccount fail")
		return err
	}

	if err := DeleteItemsByAccID(ad.db, accid); err != nil {
		logs.Error(err).Int("accid", accid).Msg("DeleteItemsByAccID fail")
		return err
	}
	return nil
}

func (ad *AccDB) UpdateAccName(accid int, accName string) error {
	acc := &AccountInfo{
		AccID:   accid,
		AccName: accName,
	}
	if err := UpdateAccount(ad.db, acc); err != nil {
		logs.Error(err).Int("accid", accid).Str("name", accName).Msg("UpdateAccount fail")
		return err
	}
	logs.Debug().Int("accid", accid).Str("name", acc.AccName).Msg("UpdateAccount succ")
	return nil
}

func (ad *AccDB) UpdateItems(accid int, props ...PropKV) error {
	// UpdateItem(db *gorm.DB, item *Item)
	var rtnErr error = nil

	item := &Item{
		Accid: accid,
	}

	for _, prop := range props {
		item.Propid = prop.PropID
		item.Prop = prop.Key
		item.Value = prop.Val

		if err := UpdateItem(ad.db, item); err != nil {
			rtnErr = err
			logs.Error(err).Int("accid", accid).Str("key", item.Prop).Str("val", item.Value).Msg("UpdateItem fail")
		} else {
			logs.Debug().Int("accid", accid).Str("key", item.Prop).Str("val", item.Value).Int("propid", item.Propid).Msg("UpdateItem succ")
		}

	}
	return rtnErr
}

type AccProps struct {
	AccID   int
	AccName string
	Props   []PropKV
}

func (ad *AccDB) QueryItem(accid int) (*AccProps, error) {
	// GetAccountByID(db *gorm.DB, id int) (*AccountInfo, error)
	// GetItemsByAccID(db *gorm.DB, accid int) ([]Item, error)

	aps := &AccProps{AccID: accid}

	acc, err := GetAccountByID(ad.db, accid)
	if err != nil {
		logs.Error(err).Int("accid", accid).Msg("GetAccountByID fail")
		return nil, err
	}
	aps.AccName = acc.AccName
	// logs.Info().Int("accid", accid).Str("name", aps.AccName).Msg("GetAccountByID succ")

	items, err := GetItemsByAccID(ad.db, accid)
	if err != nil {
		logs.Error(err).Int("accid", accid).Str("name", aps.AccName).Msg("GetItemsByAccID fail")
		return aps, err
	}

	nItem := len(items)
	// logs.Info().Int("accid", accid).Str("name", aps.AccName).Int("nItem", nItem).Msg("GetItemsByAccID succ")
	if nItem == 0 {
		return aps, nil
	}

	aps.Props = make([]PropKV, nItem)
	for i, item := range items {
		aps.Props[i].PropID = item.Propid
		aps.Props[i].Key = item.Prop
		aps.Props[i].Val = item.Value
		// logs.Info().Int("accid", accid).Str("name", aps.AccName).Int("i", i).Str("key", item.Prop).Msg("found item")
	}

	return aps, nil
}

func (ad *AccDB) SearchAcc(accName string, cb func(*AccProps)) error {
	// SearchAccountsByName(db *gorm.DB, name string) ([]AccountInfo, error)
	accs, err := SearchAccountsByName(ad.db, accName)
	if err != nil {
		logs.Error(err).Str("name", accName).Msg("SearchAccountsByName fail")
		return err
	}

	for _, acc := range accs {
		if aps, err := ad.QueryItem(acc.AccID); err == nil {
			cb(aps)
		}
	}
	return nil
}
