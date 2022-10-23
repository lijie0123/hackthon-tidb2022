package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestQuery(t *testing.T) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	defer db.Close()
	require.NoError(t, err)
	s := QueryService{db: db, logger: nil}
	query, err := s.Query(context.TODO(), QueryReq{Sql: "select * from t1"})
	require.NoError(t, err)
	fmt.Printf("%v\n", query)
}
