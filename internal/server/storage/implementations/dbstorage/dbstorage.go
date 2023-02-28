package storage

import (
	"context"
	"database/sql"
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

//-d=postgres://postgres:mypassword@localhost:5432/yandexxx

const selectLineMetricsTable = `SELECT id,type,delta,value FROM public.metrics2 WHERE id=$1;`
const selectAllMetricsTable = `SELECT id,type,delta,value FROM public.metrics2;`

const createOrUpdateIfExistsMetricsTable = `
	INSERT INTO public.metrics2 (id, type, delta,value) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE 
  	SET delta = $3, 
      	value = $4;`

const createMetricsTable = `create table public.metrics2
	(	id varchar(40) not null primary key,
		type integer not null,
		delta bigint,
		value double precision
	);`

const checkIfMetricsTableExists = `SELECT 'public.metrics2'::regclass;`

var message = []string{
	0:  "unable to connect to database",
	1:  "table metrics does not exists, a try to create:",
	2:  "server: db metrics table creation response text:",
	3:  "server: db metrics table existence check response text:",
	4:  "server: getMetrics: Unknown metric type value",
	5:  "server: NullValue type not valid",
	6:  "nil pointer in mvList",
	7:  "undefined type in type switch dbStorage",
	8:  "server sendBatch tag:",
	9:  "server: unable to rollback, error fatal",
	10: "server: unable to commit, trying again",
	11: "server: unable to prepare statement, fatal error",
	12: "server: unable to get all metrics before saving them all, fatal error",
}

type dbMetrics struct {
	id    sql.NullString
	_type sql.NullInt64
	delta sql.NullInt64
	value sql.NullFloat64
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
		log.Fatalf(mess+": %v\n", err)
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

func (s DBStorage) GetMetric(ctx context.Context, name string, MType string) (mv mVal.MetricValue, err error) {
	if !s.connectDb(ctx) {
		return nil, errors.New(message[0])
	}
	defer s.conn.Close(ctx)

	d := dbMetrics{id: sql.NullString{String: name, Valid: true}}

	switch MType {
	case "counter":
		d._type = sql.NullInt64{Int64: 1, Valid: true}
	case "gauge":
		d._type = sql.NullInt64{Int64: 2, Valid: true}
	default:
		log.Fatalf(message[4])
	}

	row := s.conn.QueryRow(ctx, selectLineMetricsTable, &d.id)

	err = row.Scan(&d.id, &d._type, &d.delta, &d.value)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, err
	}

	switch d._type.Int64 {
	case 1:
		{
			if d.delta.Valid {
				mv = mVal.NewInt(d.delta.Int64)
			}
		}
	case 2:
		{
			if d.value.Valid {
				mv = mVal.NewFloat(d.value.Float64)
			}
		}
	default:
		log.Fatalf(message[4])
	}
	return mv, nil
}
func (s DBStorage) SaveMetric(ctx context.Context, name string, mv *mVal.MetricValue) (err error) {
	var m mVal.MetricValue
	if mv == nil {
		return errors.New(message[6])
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
		return errors.New(message[7])
	}
	conn, err := s.conn.Exec(ctx, createOrUpdateIfExistsMetricsTable, name, _type, delta, value)
	if err != nil {
		log.Println(conn)
	}
	return err
}

// GetAllMetrics Restore data from database to mem storage
func (s DBStorage) GetAllMetrics(ctx context.Context) (mvList *metricsjson.MetricsMapType, err error) {
	if !s.connectDb(ctx) {
		return
	}

	rows, err := s.conn.Query(ctx, selectAllMetricsTable)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	m := make(metricsjson.MetricsMapType)
	emptyList := make(metricsjson.MetricsMapType)
	for rows.Next() {
		d := dbMetrics{}
		err = rows.Scan(&d.id, &d._type, &d.delta, &d.value)
		if err != nil {
			return nil, err
		}
		var mv mVal.MetricValue
		switch d._type.Int64 {
		case 1:
			{
				if !d.delta.Valid {
					return &emptyList, errors.New(message[5])
				}
				mv = mVal.NewInt(d.delta.Int64)
			}
		case 2:
			{
				if !d.value.Valid {
					return &emptyList, errors.New(message[5])
				}
				mv = mVal.NewFloat(d.value.Float64)
			}
		default:
			log.Fatalf(message[4])
		}
		if !d.id.Valid {
			return &emptyList, errors.New(message[5])
		}
		m[d.id.String] = mv
	}
	return &m, nil
}

// SaveAllMetrics Park data to database
func (s DBStorage) SaveAllMetrics(ctx context.Context, mvList *metricsjson.MetricsMapType) (err error) {
	log.Println("DBStorage SaveAllMetrics invoked")
	if mvList == nil {
		return errors.New(message[6])
	}

	if !s.connectDb(ctx) {
		return errors.New(message[0])
	}
	mvL := *mvList

	//Нужно прочитать все значения метрик counter для суммирования при записи
	var currentMetricsList metricsjson.MetricsMapType = make(metricsjson.MetricsMapType)
	data, err := s.GetAllMetrics(ctx)
	if err != nil {
		return errors.New(message[12])
	}
	if data != nil {
		currentMetricsList = *data
	}

	logFatalf(message[11], err)

	batch := &pgx.Batch{}
	for k, v := range mvL {
		var d dbMetrics
		switch value := v.(type) {
		case *mVal.GaugeValue:
			d = dbMetrics{
				id:    sql.NullString{String: k, Valid: true},
				_type: sql.NullInt64{Int64: 2, Valid: true},
				value: sql.NullFloat64{Float64: value.GetInternalValue().(float64), Valid: true},
				delta: sql.NullInt64{},
			}
		case *mVal.CounterValue:
			//Нужно прочитать есть ли значение метрики в базе и прибавить к текущему
			var counter int64
			if v := currentMetricsList[k]; v != nil {
				counter = v.GetInternalValue().(int64)
			}

			d = dbMetrics{
				id:    sql.NullString{String: k, Valid: true},
				_type: sql.NullInt64{Int64: 1, Valid: true},
				value: sql.NullFloat64{},
				delta: sql.NullInt64{Int64: value.GetInternalValue().(int64) + counter, Valid: true},
			}
		default:
			return errors.New(message[7])
		}
		batch.Queue(createOrUpdateIfExistsMetricsTable, d.id, d._type, d.delta, d.value)
	}

	batchResults := s.conn.SendBatch(ctx, batch)

	log.Println(mvL)
	for range mvL {
		tag, err := batchResults.Exec()
		if err != nil {
			return err
		}
		log.Println(message[8] + tag.String())
	}

	defer batchResults.Close()

	logFatalf(message[10], err)

	log.Print("data saved")
	return nil
}
