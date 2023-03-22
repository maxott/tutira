package relalg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/maxott/tutiro/core"
)

// const DB_DRIVER_DEF = "postgres"
// const DB_URL_DEF = "postgres://localhost:5432/blue_growth?sslmode=disable"

type psqAdapter struct {
	conn *pgx.Conn //*sql.DB
}

var typeInfo = pgtype.NewConnInfo()

func PostgresAdapter(url string) (DbAdapterI, error) {
	dbContext := context.Background()
	db, err := pgx.Connect(dbContext, url)
	if err != nil {
		return nil, err
	}
	// check the connection
	if err = db.Ping(dbContext); err != nil {
		return nil, err
	}
	return &psqAdapter{db}, nil
}

func (db *psqAdapter) Query(ctx context.Context, query string, args ...interface{}) (res core.TermI, err error) {
	rows, err := db.conn.Query(ctx, query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	var scanners []interface{}
	var mappers []func() core.TermI
	scanners, mappers, err = db.buildScanners(&rows)
	if err != nil {
		return
	}
	colCount := len(mappers)
	lrows := []core.TermI{}
	for rows.Next() {
		if err := rows.Scan(scanners...); err != nil {
			log.Fatal(err)
		}
		terms := make([]core.TermI, colCount)
		for i, m := range mappers {
			terms[i] = m()
		}
		lrows = append(lrows, core.List(terms))
	}
	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return
	}
	res = core.List(lrows)
	return
}

func (db *psqAdapter) buildScanners(rows *pgx.Rows) (scanners []interface{}, mappers []func() core.TermI, err error) {
	fieldDescr := (*rows).FieldDescriptions()
	scanners = make([]interface{}, len(fieldDescr))
	mappers = make([]func() core.TermI, len(fieldDescr))
	for i, fd := range fieldDescr {
		finfo, _ := typeInfo.DataTypeForOID(fd.DataTypeOID)
		switch finfo.Name {
		case "float4", "float8":
			var v float64
			scanners[i] = &v
			mappers[i] = func() core.TermI { return core.Float64(v) }
		case "int4", "int8":
			var v int64
			scanners[i] = &v
			mappers[i] = func() core.TermI { return core.Int64(v) }
		case "varchar":
			var v string
			scanners[i] = &v
			mappers[i] = func() core.TermI { return core.String(v) }
		case "timestamp":
			var v time.Time
			scanners[i] = &v
			mappers[i] = func() core.TermI { return core.String(v.Format(time.RFC3339)) }
		case "uuid":
			var v uuid.UUID
			scanners[i] = &v
			mappers[i] = func() core.TermI { return core.String(v.String()) }
		default:
			fmt.Printf("Unsupported postgres type '%s'\n", finfo.Name)
		}
		///fmt.Printf("Field: %d: %s - %s - %v\n", i, string(fd.Name), finfo.Name, finfo)
	}
	return
}
