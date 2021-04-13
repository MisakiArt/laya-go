package main

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"log"
)

func main()  {
	fmt.Println(123)
	connect, err := sql.Open("clickhouse", "tcp://cc-uf6cquovu3bypn059.ads.rds.aliyuncs.com:3306?username=dev_click_house&password=OtG5R4B3jbDbqIIR&database=dev_jing_ch_dw_s_1&debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}
	rows, err := connect.Query("select jing_uuid,mid from jing_channel_identify_actions where omni_user_uuid  = 'pwrTfQkjNNawWzW5exqSBo'")
	defer rows.Close()
	for rows.Next() {
		var (
			jing_uuid               string
			mid           uint8
		)
		if err := rows.Scan(&jing_uuid, &mid); err != nil {
			log.Fatal(err)
		}

		log.Printf("jing_uuid: %s, mid: %d", jing_uuid,mid)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	//if _, err := connect.Exec("DROP TABLE example"); err != nil {
	//	log.Fatal(err)
	//}
}