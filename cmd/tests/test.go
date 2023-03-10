package main

import (
	"database/sql"
	"log"

	metricvalueI "github.com/alphaonly/harvester/internal/server/metricvaluei"
	db "github.com/alphaonly/harvester/internal/server/storage/implementations/dbstorage"
	storage "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"golang.org/x/net/context"
)

func main() {

	var s storage.Storage
	var mv metricvalueI.MetricValue

	mv = metricvalueI.NewInt(300)
	dbURL := "postgres://postgres:mypassword@localhost:5432/yandex"
	s = db.NewDBStorage(context.Background(), dbURL)
	log.Print(s)

	//err := s.SaveMetric(context.Background(), "Poll233Counter", &mv)
	mv, err := s.GetMetric(context.Background(), "Poll233Counter", "counter")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(mv)
	ss := sql.Stmt{}
	_, err = ss.Query()
	if err != nil {
		log.Fatal(err) //
	}

}
