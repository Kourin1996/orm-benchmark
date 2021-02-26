package benchs

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/naoina/genmai"
)

var genmaidb *genmai.DB

func init() {
	st := NewSuite("genmai")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, GenmaiInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, GenmaiInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, GenmaiUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, GenmaiRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, GenmaiReadSlice)

		db, err := genmai.New(&genmai.PostgresDialect{}, ORM_SOURCE)
		checkErr(err)
		genmaidb = db
	}
}

type GenmaiModel struct {
	Id      int `db:"pk" column:"tbl_id"`
	Name    string
	Title   string
	Fax     string
	Web     string
	Age     int
	Right   bool
	Counter int64
}

func NewGenmaiModel() *GenmaiModel {
	m := new(GenmaiModel)
	m.Name = "Orm Benchmark"
	m.Title = "Just a Benchmark for fun"
	m.Fax = "99909990"
	m.Web = "http://blog.milkpod29.me"
	m.Age = 100
	m.Right = true
	m.Counter = 1000

	return m
}

func initGenmaiDB() {
	genmaidb.DropTable(&GenmaiModel{})

	err := genmaidb.CreateTable(&GenmaiModel{})
	checkErr(err)
}

func GenmaiInsert(b *B) {
	var m *GenmaiModel
	wrapExecute(b, func() {
		initGenmaiDB()
		m = NewGenmaiModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = i
		_, d := genmaidb.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GenmaiInsertMulti(b *B) {
	var ms []GenmaiModel
	wrapExecute(b, func() {
		initDB3()
		ms = make([]GenmaiModel, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, *NewGenmaiModel())
		}
	})

	for i := 0; i < b.N; i++ {
		_, d := genmaidb.Insert(ms)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GenmaiUpdate(b *B) {
	var m *GenmaiModel
	wrapExecute(b, func() {
		initGenmaiDB()
		m = NewGenmaiModel()

		_, d := genmaidb.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		_, d := genmaidb.Update(*m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GenmaiRead(b *B) {
	var m *GenmaiModel
	wrapExecute(b, func() {
		initGenmaiDB()
		m = NewGenmaiModel()

		_, d := genmaidb.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		var results []GenmaiModel
		d := genmaidb.Select(&results, genmaidb.Limit(1))
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func GenmaiReadSlice(b *B) {
	var m *GenmaiModel
	wrapExecute(b, func() {
		initGenmaiDB()
		m = NewGenmaiModel()
		m.Id = 0
		_, d := genmaidb.Insert(m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		var models []GenmaiModel

		d := genmaidb.Select(&models, genmaidb.Where("tbl_id", ">", 0).Limit(100))
		if d != nil {
			fmt.Println(d)
			b.FailNow()
		}
	}
}
