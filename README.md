# ORM Benchmark

Forked from [yusaer/orm-benchmark](https://github.com/yusaer/orm-benchmark).  
Please refer this first.  
I upgraded versions and added benchmark for ent.

### Environment

* go version go1.6 linux/amd64

### PostgreSQL

* PostgreSQL 9.5 for Linux on x86_64

### ORMs

All package run in no-cache mode.

* [dbr](https://github.com/gocraft/dbr) (in preparation)
* [genmai](https://github.com/naoina/genmai) (in preparation)
* [gorm](https://github.com/jinzhu/gorm)
* [gorp](https://github.com/go-gorp/gorp) (in preparation)
* [pg](https://github.com/go-pg/pg)
* [beego/orm](https://github.com/astaxie/beego/tree/master/orm)
* [sqlx](https://github.com/jmoiron/sqlx) (in preparation)
* [xorm](https://github.com/xormplus/xorm)

* [ent](https://github.com/ent/ent)
	
### Run

```go
go get github.com/Kourin1996/orm-benchmark
# build
go install
# all
orm-benchmark -multi=20 -orm=all
# portion
orm-benchmark -multi=20 -orm=xorm
```
