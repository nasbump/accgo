package main

import (
	"accgo/accdb"
	"accgo/utils/encrypts"
	"accgo/utils/logs"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type kvSlice []accdb.PropKV

// 实现 flag.Value 接口
func (s *kvSlice) String() string {
	var sb strings.Builder
	for _, kv := range *s {
		sb.WriteString(kv.Key)
		sb.WriteString(":")
		sb.WriteString(kv.Val)
		sb.WriteString(",")
	}
	return sb.String()
}

func (s *kvSlice) Set(value string) error {
	kv := strings.SplitN(value, ":", 2)

	p := accdb.PropKV{
		Key: kv[0],
		Val: kv[1],
	}

	*s = append(*s, p)
	return nil
}

var dbToken string

func main() {
	var forSearch, forQuery, forAdd, forDelete, forUpdate bool
	var name string
	var props kvSlice
	var accid, propid int

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.BoolVar(&forSearch, "search", false, "search -n name")
	fs.BoolVar(&forQuery, "query", false, "query  -i accid")
	fs.BoolVar(&forAdd, "add", false, "add -n name -v key:value... or add -i accid -v key:value...")
	fs.BoolVar(&forDelete, "delete", false, "delete -i accid [-p propid]...")
	fs.BoolVar(&forUpdate, "update", false, "update -i accid -n new_name or  update -i accid -p propid -v key:value")

	fs.StringVar(&name, "n", "", "account name")
	fs.StringVar(&dbToken, "t", "", "token for encrypt")
	fs.Var(&props, "v", "key:value for prop")

	fs.IntVar(&accid, "i", -1, "account id")
	fs.IntVar(&propid, "p", -1, "prop id")

	var dbpath string
	var debug bool
	fs.StringVar(&dbpath, "db", "accgo.db", "db filepath")
	fs.BoolVar(&debug, "debug", false, "enable debug mode")

	fs.Parse(os.Args[1:])

	if dbToken == "" {
		fs.Usage()
		return
	}

	if debug {
		logs.LogsInit(0)
	} else {
		logs.LogsInit(7) // disable logger
	}

	ad, err := accdb.Start(dbpath)
	if err != nil {
		logs.Catch(err).Send()
		return
	}

	if err := ad.Init(); err != nil {
		logs.Catch(err).Send()
		return
	}

	for i, prop := range props {
		enc, err := encrypts.GoAhead(prop.Val, dbToken)
		if err != nil {
			log.Println("invalid token:", dbToken)
			return
		}
		props[i].Val = enc
	}

	if forSearch {
		ad.SearchAcc(name, showAccProps)
		return
	}
	if forQuery {
		if ap, _ := ad.QueryItem(accid); ap != nil {
			showAccProps(ap)
		}
		return
	}
	if forAdd {
		if accid > 0 { // 有accid，则添加属性
			ad.AddProps(accid, props...)
		} else { // 创建新的accid，并添加属性
			ad.NewAcc(name, props...)
		}
		return
	}
	if forDelete {
		propids := []int{}
		if propid > 0 {
			propids = append(propids, propid)
		}
		ad.DelAcc(accid, propids...)
		return
	}
	if forUpdate {
		if name == "" { // 更新accid下的属性
			if len(props) > 0 {
				props[0].PropID = propid
			}
			ad.UpdateItems(accid, props...)
		} else { // 更新accid的名称
			ad.UpdateAccName(accid, name)
		}
		return
	}

	fs.Usage()
}

func showAccProps(ap *accdb.AccProps) {
	fmt.Printf("accid: %d, accname: %s\n", ap.AccID, ap.AccName)
	for _, prop := range ap.Props {
		dec, err := encrypts.ComeBack(prop.Val, dbToken)
		if err != nil {
			log.Println("invalid token:", dbToken)
			// os.Exit(1)
			dec = prop.Val
		}

		fmt.Printf("  -> prop: %d, k: %s, v: %s\n", prop.PropID, prop.Key, dec)
	}
}
