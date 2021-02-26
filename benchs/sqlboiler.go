package benchs

import (
	"context"
	"fmt"

	"database/sql"

	"github.com/Kourin1996/orm-benchmark/models"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var sbDb *sql.DB

func init() {
	st := NewSuite("sqlboiler")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, SqlboilerInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, SqlboilerInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, SqlboilerUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, SqlboilerRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, SqlboilerReadSlice)

		db, err := sql.Open("postgres", ORM_SOURCE)
		checkErr(err)

		boil.SetDB(db)
		sbDb = db
	}
}

func NewSqlboilerModel() *models.Model {
	m := new(models.Model)
	m.Name = null.StringFrom("Orm Benchmark")
	m.Title = null.StringFrom("Just a Benchmark for fun")
	m.Fax = null.StringFrom("99909990")
	m.Web = null.StringFrom("http://blog.milkpod29.me")
	m.Age = null.IntFrom(100)
	m.Right = null.BoolFrom(true)
	m.Counter = null.Int64From(1000)

	return m
}

func SqlboilerInsert(b *B) {
	var m *models.Model
	wrapExecute(b, func() {
		initDB()
		m = NewSqlboilerModel()
	})
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		m.ID = 0
		d := m.Insert(ctx, sbDb, boil.Infer())
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func SqlboilerInsertMulti(b *B) {
	panic(fmt.Errorf("Not support multi insert"))
}

func SqlboilerUpdate(b *B) {
	var m *models.Model
	wrapExecute(b, (func() {
		initDB()
		m = NewSqlboilerModel()

		ctx := context.Background()
		d := m.Insert(ctx, sbDb, boil.Infer())
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}))

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, d := m.Update(ctx, sbDb, boil.Infer())
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func SqlboilerRead(b *B) {
	var m *models.Model
	wrapExecute(b, (func() {
		initDB()
		m = NewSqlboilerModel()

		ctx := context.Background()
		d := m.Insert(ctx, sbDb, boil.Infer())
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}))

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, d := models.FindModel(ctx, sbDb, 1)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func SqlboilerReadSlice(b *B) {
	var m *models.Model
	wrapExecute(b, (func() {
		initDB()
		m = NewSqlboilerModel()

		ctx := context.Background()
		for i := 0; i < 100; i++ {
			m.ID = 0
			d := m.Insert(ctx, sbDb, boil.Infer())
			if d != nil {
				fmt.Println(d.Error())
				b.FailNow()
			}
		}
	}))

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, d := models.Models(qm.Where("id > $1", 0), qm.Limit(100)).All(ctx, sbDb)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}
