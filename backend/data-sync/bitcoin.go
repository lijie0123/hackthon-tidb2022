package main

import "cloud.google.com/go/bigquery"

var BlockConverter = Converter{
	InputSql:  "SELECT * FROM `bigquery-public-data.crypto_bitcoin.blocks` WHERE timestamp_month<\"2020-01-01\"",
	BatchSize: 100,
	OutputSchemes: []TableDesc{
		{Table: "bitcoin_block", Fields: []string{
			"hash",
			"size",
			"stripped_size",
			"weight",
			"number",
			"version",
			"merkle_root",
			"timestamp",
			"nonce",
			"bits",
			"coinbase_param",
			"transaction_count",
		},
		},
	},
	ParseOne: func(values []bigquery.Value) (RecordsOfTable, error) {
		return map[string]Records{
			"bitcoin_block": {{
				"hash":              values[0],
				"size":              values[1],
				"stripped_size":     values[2],
				"weight":            values[3],
				"number":            values[4],
				"version":           values[5],
				"merkle_root":       values[6],
				"timestamp":         values[7],
				"nonce":             values[9],
				"bits":              values[10],
				"coinbase_param":    values[11],
				"transaction_count": values[12],
			},
			},
		}, nil
	},
}

var TransactionConverter = Converter{
	InputSql:  "SELECT * FROM `bigquery-public-data.crypto_bitcoin.transactions` WHERE block_timestamp_month >= '2022-08-01' and block_timestamp_month <= '2022-09-01';",
	BatchSize: 100,
	OutputSchemes: []TableDesc{
		{
			Table: "bitcoin_transaction", Fields: []string{
				"hash",
				"size",
				"virtual_size",
				"version",
				"lock_time",
				"block_hash",
				"block_number",
				"block_timestamp",
				"input_count",
				"output_count",
				"input_value",
				"output_value",
				"is_coinbase",
				"fee",
			},
		},
		{
			Table: "bitcoin_transaction_input", Fields: []string{
				"transaction_hash",
				"index",
				"spent_transaction_hash",
				"spent_output_index",
				"script_asm",
				"script_hex",
				"sequence",
				"required_signatures",
				"type",
				"addresses",
				"value",
			},
		},
		{
			Table: "bitcoin_transaction_output", Fields: []string{
				"transaction_hash",
				"index",
				"script_asm",
				"script_hex",
				"required_signatures",
				"type",
				"addresses",
				"value",
			},
		},
	},
	ParseOne: func(values []bigquery.Value) (RecordsOfTable, error) {
		var inputs Records
		for _, ivalues := range values[15].([]bigquery.Value) {
			iv := ivalues.([]bigquery.Value)
			inputs = append(inputs, map[string]any{
				"transaction_hash":       values[0],
				"index":                  iv[0],
				"spent_transaction_hash": iv[1],
				"spent_output_index":     iv[2],
				"script_asm":             iv[3],
				"script_hex":             iv[4],
				"sequence":               iv[5],
				"required_signatures":    iv[6],
				"type":                   iv[7],
				"addresses":              strs2json(iv[8]),
				"value":                  rat2float(iv[9], values[0]),
			})
		}

		var outputs Records
		for _, ivalues := range values[16].([]bigquery.Value) {
			iv := ivalues.([]bigquery.Value)
			outputs = append(outputs, map[string]any{
				"transaction_hash":    values[0],
				"index":               iv[0],
				"script_asm":          iv[1],
				"script_hex":          iv[2],
				"required_signatures": iv[3],
				"type":                iv[4],
				"addresses":           strs2json(iv[5]),
				"value":               rat2float(iv[6], values[0]),
			})
		}

		return map[string]Records{
			"bitcoin_transaction": {{
				"hash":            values[0],
				"size":            values[1],
				"virtual_size":    values[2],
				"version":         values[3],
				"lock_time":       values[4],
				"block_hash":      values[5],
				"block_number":    values[6],
				"block_timestamp": values[7],
				"input_count":     values[9],
				"output_count":    values[10],
				"input_value":     rat2float(values[11], values[0]),
				"output_value":    rat2float(values[12], values[0]),
				"is_coinbase":     values[13],
				"fee":             rat2float(values[14], values[0]),
			},
			},
			"bitcoin_transaction_input":  inputs,
			"bitcoin_transaction_output": outputs,
		}, nil
	},
}
