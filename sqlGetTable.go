package sqlGetTable

import (
	"fmt"
	"time"
	"database/sql"
	"encoding/json"
	_"github.com/lib/pq"
)

type Conn struct {
	Port	string
	Host	string
	User	string
	Pass	string
	Name	string
	SSL	string
	Type	string
}

func GetAll(q string, DB *Conn) ([]string, error) {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		DB.Host, DB.Port, DB.User, DB.Pass, DB.Name, DB.SSL)
	db, err := sql.Open(DB.Type, dbinfo)
	if (err != nil) {
		fmt.Println(err)
	}
	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	jsonMap := make([]string, 0)
	record := make(map[string]interface{})
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, _ := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if (ok) {
				v = string(b)
			} else {
				v = val
			}
			_ = v
		}
		for i, col := range values {
			switch t := col.(type) {
			default:
				fmt.Printf("Unexpected type %T\n", t)
			case nil:
				record[columns[i]] = ""
			case bool:
				record[columns[i]] = col.(bool)
			case int:
				record[columns[i]] = col.(int)
			case int16:
				record[columns[i]] = col.(int16)
			case int32:
				record[columns[i]] = col.(int32)
			case int64:
				record[columns[i]] = col.(int64)
			case float64:
				record[columns[i]] = col.(float64)
			case string:
				record[columns[i]] = col.(string)
			case []byte:   // -- all cases go HERE!
				record[columns[i]] = string(col.([]byte))
			case time.Time:
//				record[columns[i]] = string(col.(string))
			}
		}
		kkk, err := json.Marshal(record)
		if err != nil {
			fmt.Println(err)
		}
		jsonMap = append(jsonMap, string(kkk))
	}
	db.Close()
	for _, GG := range jsonMap {
		var data map[string]interface{}
		err = json.Unmarshal([]byte(GG), &data)
	}
	return jsonMap, err
}
