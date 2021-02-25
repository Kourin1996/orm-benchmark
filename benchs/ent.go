package benchs

import (
	"context"
	"fmt"
	"log"

	"github.com/Kourin1996/orm-benchmark/ent"
	"github.com/Kourin1996/orm-benchmark/ent/model"
)

var entdb *ent.Client

func init() {
	st := NewSuite("ent")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, EntInsert)
		st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, EntInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, EntUpdate)
		st.AddBenchmark("Read", 4000*ORM_MULTI, EntRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, EntReadSlice)

		client, err := ent.Open("postgres", ORM_SOURCE)
		if err != nil {
			log.Fatalf("Error open mysql ent client: %v\n", err)
		}

		entdb = client
	}
}

func EntInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := entdb.Model.Create().SetName(m.Name).SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetRight(m.Right).SetCounter(m.Counter).Save(ctx)

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func EntInsertMulti(b *B) {
	wrapExecute(b, func() {
		initDB()
	})

	for i := 0; i < b.N; i++ {
		bulk := make([]*ent.ModelCreate, 100)
		for i := 0; i < 100; i++ {
			m := NewModel()
			bulk[i] = entdb.Model.Create().SetName(m.Name).SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetRight(m.Right).SetCounter(m.Counter)
		}

		ctx := context.Background()
		_, err := entdb.Model.CreateBulk(bulk...).Save(ctx)

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func EntUpdate(b *B) {
	var m *Model
	var entModel *ent.Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()

		ctx := context.Background()
		em, err := entdb.Model.Create().SetName(m.Name).SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetRight(m.Right).SetCounter(m.Counter).Save(ctx)
		entModel = em

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := entModel.Update().SetName(m.Name).SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetRight(m.Right).SetCounter(m.Counter).Save(ctx)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func EntRead(b *B) {
	var m *Model

	wrapExecute(b, func() {
		initDB()
		m = NewModel()

		ctx := context.Background()
		_, err := entdb.Model.Create().SetName(m.Name).SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetRight(m.Right).SetCounter(m.Counter).Save(ctx)

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := entdb.Model.Query().First(ctx)

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func EntReadSlice(b *B) {
	wrapExecute(b, func() {
		initDB()

		bulk := make([]*ent.ModelCreate, 100)
		for i := 0; i < 100; i++ {
			m := NewModel()
			bulk[i] = entdb.Model.Create().SetName(m.Name).SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetRight(m.Right).SetCounter(m.Counter)
		}

		ctx := context.Background()
		_, err := entdb.Model.CreateBulk(bulk...).Save(ctx)

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := entdb.Model.Query().Where(
			model.IDGT(0),
		).Limit(100).All(ctx)

		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}
