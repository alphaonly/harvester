package storage

import (
	"context"
	"errors"
	"log"

	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"
	mVal "github.com/alphaonly/harvester/internal/server/metricvalueInt"
	storage "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"github.com/jackc/pgx/v5"
)

//	type Storage interface {
//		GetMetric(ctx context.Context, name string) (mv *M.MetricValue, err error)
//		SaveMetric(ctx context.Context, name string, mv *M.MetricValue) (err error)
//		GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error)
//		SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error)
//	}

const createMetricsTable = `create table public.metrics
	(	id varchar(40) not null primary key,
		type integer not null,
		delta integer,
		value double precision
	);`

const checkIfMetricsTableExists = `SELECT 'public.metrics'::regclass;`
const insertLineIntoMetricsTable = `
	INSERT INTO public.metrics (id, type, delta, value)VALUES ($1, $2, $3, $4);`

var message = []string{
	"unable to connect to database",
	"table metrics does not exists, a try to create:",
	"server: db metrics table creation response text:",
	"server: db metrics table existence check response text:",
}

type DBStorage struct {
	dataBaseUrl string
	conn        *pgx.Conn
}

func NewDBStorage(ctx context.Context, dataBaseUrl string) storage.Storage {
	//get params
	s := DBStorage{dataBaseUrl: dataBaseUrl}
	//connect db
	var err error
	s.conn, err = pgx.Connect(ctx, s.dataBaseUrl)
	if err != nil {
		logFatalf(message[0], err)
		return nil
	}
	defer s.conn.Close(ctx)

	// check metrics table exists
	resp, err := s.conn.Exec(context.Background(), checkIfMetricsTableExists)
	if err != nil {
		log.Println(message[1] + err.Error())
		//create metrics Table
		resp, err = s.conn.Exec(context.Background(), createMetricsTable)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message[2] + resp.String())
	} else {
		log.Println(message[3] + resp.String())
	}
	return s
}

func logFatalf(mess string, err error) {
	if err != nil {
		log.Fatalf(mess+": "+" %v\n", err)
	}
}
func (s *DBStorage) connectDb(ctx context.Context) (ok bool) {
	ok = false
	var err error

	if s.conn == nil {
		s.conn, err = pgx.Connect(ctx, s.dataBaseUrl)
	} else {
		err = s.conn.Ping(ctx)
		if err != nil {
			s.conn, err = pgx.Connect(ctx, s.dataBaseUrl)
		}
	}
	logFatalf(message[0], err)
	ok = true

	return ok
}

func (s DBStorage) GetMetric(ctx context.Context, name string) (mv mVal.MetricValue, err error) {
	if !s.connectDb(ctx) {
		return
	}

	defer s.conn.Close(ctx)
	return nil, nil
}
func (s DBStorage) SaveMetric(ctx context.Context, name string, mv *mVal.MetricValue) (err error) {
	var m mVal.MetricValue
	if mv == nil {
		return errors.New("nil pointer in mv")
	}
	m = *mv
	if !s.connectDb(ctx) {
		return
	}
	var (
		_type int
		delta int64
		value float64
	)

	switch v := m.GetInternalValue().(type) {
	case int64:
		{
			_type = 1
			delta = v
		}
	case float64:
		{
			_type = 2
			value = v
		}
	default:
		return errors.New("undefined type in type switch dbStorage::SaveMetric")
	}

	s.conn.Exec(ctx, insertLineIntoMetricsTable,name,_type,delta,value)

	return nil
}

// Restore data from database to mem storage
func (s DBStorage) GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error) {
	if !s.connectDb(ctx) {
		return
	}

	return mvList, nil
}

// Park data to database
func (s DBStorage) SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error) {
	if !s.connectDb(ctx) {
		return
	}

	return nil
}
