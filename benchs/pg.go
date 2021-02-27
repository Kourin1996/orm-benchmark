package benchs

import (
	"fmt"

	pg "github.com/go-pg/pg/v10"
)

var pgdb *pg.DB

func init() {
	st := NewSuite("pg")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, PgInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, PgInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, PgUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, PgRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, PgReadSlice)

		db := pg.Connect(&pg.Options{
			Addr:     "localhost:5432",
			User:     "postgres",
			Password: "postgres",
			Database: "test",
		})
		pgdb = db
	}
}

func PgInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := pgdb.Model(m).Insert(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgInsertMulti(b *B) {
	var ms []interface{}
	wrapExecute(b, func() {
		initDB()
	})

	for i := 0; i < b.N; i++ {
		ms = make([]interface{}, 0, 100)
		for i := 0; i < 100; i++ {
			var m interface{} = NewModel()
			ms = append(ms, m)
		}
		_, err := pgdb.Model(ms...).Insert()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := pgdb.Model(m).Insert(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		_, err := pgdb.Model(m).WherePK().Update()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := pgdb.Model(m).Insert(); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		err := pgdb.Model(m).WherePK().Select()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func PgReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < 100; i++ {
			m.Id = 0
			if _, err := pgdb.Model(m).Insert(); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})

	for i := 0; i < b.N; i++ {
		var models []*Model
		err := pgdb.Model(&models).Where("id > ?", 0).Limit(100).Select()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
