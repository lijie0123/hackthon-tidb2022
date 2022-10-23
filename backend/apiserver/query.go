package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type QueryService struct {
	db     *sql.DB
	logger *log.Logger
}

type QueryReq struct {
	Sql   string `json:"sql"`
	Args  []any  `json:"args"`
	Limit int    `json:"limit"`
}

type QueryRes struct {
	Schema []Scheme `json:"schema"`
	Rows   []any    `json:"rows"`
}

type Scheme struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *QueryService) Query(ctx context.Context, req QueryReq) (*QueryRes, error) {
	sqlStr := strings.TrimSpace(strings.ToLower(req.Sql))
	if !(strings.HasPrefix(sqlStr, "select") ||
		strings.HasPrefix(sqlStr, "show") ||
		strings.HasPrefix(sqlStr, "desc")) {
		return nil, fmt.Errorf("%w: only DQL is supported", BadRequest)
	}
	rows, err := s.db.QueryContext(ctx, req.Sql, req.Args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", BadRequest, err.Error())
	}
	defer rows.Close()
	rt := QueryRes{}
	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	for _, v := range cols {
		rt.Schema = append(rt.Schema, Scheme{Name: v.Name(), Type: v.DatabaseTypeName()})
	}
	var limit = 100
	if req.Limit != 0 {
		limit = req.Limit
	}
	count := len(rt.Schema)
	for i := 0; i < limit; i++ {
		succ := rows.Next()
		if !succ {
			break
		}

		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		row, err := convertValues(values)
		if err != nil {
			return nil, err
		}
		rt.Rows = append(rt.Rows, row)
	}
	return &rt, nil
}

func convertValues(values []any) (rt []any, err error) {
	for _, v := range values {

		x := v.([]byte)

		if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
			rt = append(rt, nx)
		} else if b, ok := strconv.ParseBool(string(x)); ok == nil {
			rt = append(rt, b)
		} else if "string" == fmt.Sprintf("%T", string(x)) {
			rt = append(rt, string(x))
		} else {
			return nil, fmt.Errorf("Failed on if for type %T of %v", x, x)
		}
	}
	return
}
