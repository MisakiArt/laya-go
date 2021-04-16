package main

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"log"
	"sync"
)

func main()  {
	fmt.Println("start")
	connect, err := sql.Open("clickhouse", "tcp://cc-uf6cquovu3bypn059.ads.rds.aliyuncs.com:3306?username=dev_click_house&password=OtG5R4B3jbDbqIIR&database=dev_jing_ch_dw_s_1&debug=false")
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

	partitionArray :=[...]string{"202104","202103","202102","202101"}
	log.Println(len(partitionArray))
	ch :=make(chan string,len(partitionArray))
	var wg sync.WaitGroup
	for _,v :=range partitionArray {
		wg.Add(1)
		go optimizeTable(connect,v,ch,&wg)
	}
	//close(ch)
	for i:=0 ; i<4 ;i++ {
		log.Println(<-ch)
	}
	wg.Wait()
}

func optimizeTable(connect *sql.DB, partition string,ch chan<- string,wg *sync.WaitGroup) {
	defer wg.Done()
	sqlString := "optimize table dev_jing_ch_dw_s_1.local_jing_channel_identify_actions     ON CLUSTER default PARTITION  '"+partition+"' final  DEDUPLICATE"
	log.Println(sqlString)
	_,err := connect.Exec(sqlString)
	if err != nil {
		log.Fatal(err)
		ch <- partition+err.Error()
	} else  {
		ch <- partition+" optimize success"
	}

}