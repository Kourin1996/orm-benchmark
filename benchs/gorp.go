package benchs

import (
	"fmt"

	"database/sql"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

var dbmap *gorp.DbMap

func init() {
	st := NewSuite("gorp")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, GorpInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, GorpInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, GorpUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, GorpRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, GorpReadSlice)

		db, err := sql.Open("postgres", ORM_SOURCE)
		checkErr(err)
		d := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
		d.AddTableWithName(Model{}, "models").SetKeys(true, "Id")
		dbmap = d
	}
}

func initGorpDB() {
	err := dbmap.DropTableIfExists(Model{})
	checkErr(err)

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err)
}

func GorpInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initGorpDB()
		m = NewModel()
	})
	for i := 0; i < b.N; i++ {
		m.Id = 0
		d := dbmap.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GorpInsertMulti(b *B) {
	panic(fmt.Errorf("Not support multi insert"))
}

func GorpUpdate(b *B) {
	var m *Model
	wrapExecute(b, (func() {
		initGorpDB()
		m = NewModel()

		d := dbmap.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}))

	for i := 0; i < b.N; i++ {
		_, d := dbmap.Update(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GorpRead(b *B) {
	var m *Model
	wrapExecute(b, (func() {
		initGorpDB()
		m = NewModel()

		d := dbmap.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}))

	for i := 0; i < b.N; i++ {
		d := dbmap.SelectOne(m, "select * from models where id = $1", m.Id)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GorpReadSlice(b *B) {
	var m *Model
	wrapExecute(b, (func() {
		initGorpDB()
		m = NewModel()

		for i := 0; i < 100; i++ {
			m.Id = 0
			d := dbmap.Insert(m)
			if d != nil {
				fmt.Println(d.Error())
				b.FailNow()
			}
		}
	}))

	for i := 0; i < b.N; i++ {
		var models []Model
		_, d := dbmap.Select(&models, "select * from models where id > $1 limit 100", m.Id)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}
