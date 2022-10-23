package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	ctx := context.Background()
	googleConn := getGoogleConn(ctx)
	defer googleConn.Close()
	tidbConn := getTidbConn()
	defer tidbConn.Close()
	switch os.Args[1] {
	case "btctx":
		ctr := migrateBitCoinTx(os.Args[2])
		DoMigrate(ctx, ctr, googleConn, tidbConn)
	case "btcblk":
		ctr := migrateBitCoinBlock(os.Args[2])
		DoMigrate(ctx, ctr, googleConn, tidbConn)
	case "cbt":
		var since string

		if strings.HasPrefix(os.Args[2], "now-") {
			duration, err := time.ParseDuration(strings.TrimPrefix(os.Args[2], "now-"))
			if err != nil {
				panic(err)
			}
			since = time.Now().Add(-duration).UTC().Format(time.RFC3339)
		} else {
			since = os.Args[2]
		}
		ctr1 := migrateCBitCoinBlock(since)
		ctr2 := migrateCBitCoinTx(since)
		wg.Add(2)
		go func() {
			defer wg.Done()
			DoMigrate(ctx, ctr1, googleConn, tidbConn)
		}()
		go func() {
			defer wg.Done()
			DoMigrate(ctx, ctr2, googleConn, tidbConn)
		}()
		wg.Wait()
	default:
		panic("bad arg")
	}

}

func migrateBitCoinTx(string2 string) Converter {
	proto := TransactionConverter
	proto.InputSql = "SELECT * FROM `bigquery-public-data.crypto_bitcoin.transactions` WHERE " + string2 + ";"
	return proto
}

func migrateBitCoinBlock(string2 string) Converter {
	proto := BlockConverter
	proto.InputSql = "SELECT * FROM `bigquery-public-data.crypto_bitcoin.blocks` WHERE " + string2 + ";"
	return proto
}

func migrateCBitCoinBlock(now string) Converter {
	proto := BlockConverter
	proto.InputSql = fmt.Sprintf("SELECT * FROM `bigquery-public-data.crypto_bitcoin.blocks` WHERE timestamp_month>='2022-10-01' and timestamp >= '%s'", now)
	proto.BatchSize = 1
	proto.IgnoreNoInsert = true
	return proto
}

func migrateCBitCoinTx(now string) Converter {
	proto := TransactionConverter
	proto.InputSql = fmt.Sprintf("SELECT * FROM `bigquery-public-data.crypto_bitcoin.transactions` WHERE block_timestamp_month='2022-10-01' and block_timestamp >= '%s'", now)
	proto.BatchSize = 10
	proto.IgnoreNoInsert = true
	return proto
}

func DoMigrate(ctx context.Context, converter Converter, googleConn *bigquery.Client, tidbConn *sqlx.DB) {
	q := googleConn.Query(converter.InputSql)
	readIter, err := q.Read(ctx)
	if err != nil {
		panic(err)
	}

	var counter int = 0
	var inserted = map[string]int{}

	for eof := false; !eof; {
		var batch = RecordsOfTable{}
		for i := 0; i < converter.BatchSize; i++ {
			var values []bigquery.Value
			err := readIter.Next(&values)
			if err == iterator.Done { // from "google.golang.org/api/iterator"
				eof = true
				break
			}
			if err != nil {
				if err != nil {
					panic(err)
				}
			}
			tableOfOne, err := converter.ParseOne(values)
			if err != nil {
				panic(err)
			}
			batch.Merge(tableOfOne)
			counter++
		}
		for _, table := range converter.OutputSchemes {
			if len(batch[table.Table]) == 0 {
				continue
			}
			affect, err := tidbConn.NamedExecContext(ctx, GenTableInsertSql(table), batch[table.Table])
			if err != nil {
				panic(err)
			}
			affected, err := affect.RowsAffected()
			inserted[table.Table] = inserted[table.Table] + int(affected)
			if affected > 0 || !converter.IgnoreNoInsert {
				fmt.Printf("source offset %d, inserted %d to %s\n", counter, affected, table.Table)
			}
		}
	}
	fmt.Printf("\n--- source has %d, inserted: %v\n", counter, inserted)
}

type Converter struct {
	InputSql       string
	BatchSize      int
	OutputSchemes  []TableDesc
	ParseOne       func([]bigquery.Value) (RecordsOfTable, error)
	IgnoreNoInsert bool
}

type TableDesc struct {
	Table  string
	Fields []string
}

type Records []map[string]any
type RecordsOfTable map[string]Records

func (o RecordsOfTable) Merge(r RecordsOfTable) {
	for k, v := range r {
		o[k] = append(o[k], v...)
	}
}

func GenTableInsertSql(table TableDesc) string {
	var quotedFields []string
	var scFields []string
	for _, field := range table.Fields {
		quotedFields = append(quotedFields, "`"+field+"`")
		scFields = append(scFields, ":"+field)
	}
	fieldsStr := strings.Join(quotedFields, ",")
	holderStr := strings.Join(scFields, ",")
	return fmt.Sprintf("insert ignore into %s (%s) values (%s)", table.Table, fieldsStr, holderStr)
}

func getGoogleConn(ctx context.Context) *bigquery.Client {
	authFile := os.Getenv("GAUTH")
	if authFile == "" {
		panic("no auth file")
	}
	gProject := os.Getenv("GPROJ")
	if gProject == "" {
		panic("no google project")
	}
	file, err := os.Open(authFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	credentialsJSON, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	c, err := bigquery.NewClient(ctx, gProject,
		option.WithCredentialsJSON(credentialsJSON),
	)
	if err != nil {
		panic(err)
	}
	return c
}

func getTidbConn() *sqlx.DB {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic(fmt.Errorf("please export DSN=<>"))
	}
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func rat2float(v any, hash any) any {
	if v == nil {
		return nil
	}
	fee, exact := (v).(*big.Rat).Float64()
	if !exact {
		fmt.Printf("convert fee not exact: %v, hash %s\n", fee, hash)
	}
	return fee
}

func strs2json(value bigquery.Value) []byte {
	if value == nil {
		return []byte("null")
	}
	vs := value.([]bigquery.Value)
	var strs []string
	for _, i2 := range vs {
		strs = append(strs, i2.(string))
	}
	marshal, err := json.Marshal(strs)
	if err != nil {
		fmt.Printf("marshual json err \n")
	}
	return marshal
}
