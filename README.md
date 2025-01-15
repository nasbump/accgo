# accgo
命令行式的密码管理软件，极简程序员风格
- 轻量无依赖：使用sqlite做为数据库
- 自定义私钥，仅需记住这一把钥匙
- 数据存储使用ChaCha20-Poly1305加密

# 设计
- 账号表：t_acclist，存储账号名称，可以是网站名、APP名、银行名等，用于助记，便于搜索
```
create table if not exists t_acclist (
    acc_id integer primary key AUTOINCREMENT,
    acc_name  text not null)
```
- 属性表：t_items，以KV形式存储，纪录账号下的各项属性，如用户名、手机号、密码、邮箱、等任意KV对
```
create table if not exists t_items (
    propid integer primary key AUTOINCREMENT,
    accid integer not null,
    prop   text not null,
    value  text not null)
```

# 操作
### usage
```
usage of ./bin/accgo:
  -add
        add -n name -v key:value... or add -i accid -v key:value...
  -db string
        db filepath (default "accgo.db")
  -debug
        enable debug mode
  -delete
        delete -i accid [-p propid]...
  -i int
        account id (default -1)
  -n string
        account name
  -p int
        prop id (default -1)
  -query
        query  -i accid
  -search
        search -n name
  -t string
        token for encrypt
  -update
        update -i accid -n new_name or  update -i accid -p propid -v key:value
  -v value
        key:value for prop
```

### 创建账号密码
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" add -n "百度" -v "用户名:myuser" -v "密码:mypwd" -v "手机:13000001111" -v "备注:文小言可用"`
上述命令返回 accid，做为“百度”这个账号的唯一索引

### 增加账号属性KV
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" add -i "<ACCID>" -v "昵称:nickname"`
上述命令返回 propid，使用 ACCID 这个账号下的 "昵称"这个属性的唯一索引

### 删除账号及账号下所有属性KV
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" delete -i "<ACCID>"`

### 删除账号某个属性KV
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" delete -i "<ACCID>" -p "<PROPID>"`

### 修改账号名称
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" update -i "<ACCID>" -n "新名称"`

### 修改账号下某一属性KV
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" update -i "<ACCID>" -p "<PROPID>" -v "新的KEY:新的VALUE"`

### 查询账号
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" query -i "<ACCID>"`

### 搜索账号
` ./bin/accgo -t "<唯一密钥>" -db "<path-to>/mypass.db" search -n "账号关键词"`


