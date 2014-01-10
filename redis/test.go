package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	total_data := 100000
	con_count := 20

	ch := make(chan int, 1)
	dbPool.InitPool(con_count, init_redis)

	ts := time.Now()
	for i := 0; i < con_count; i++ {
		go loopWrite(total_data/con_count, ch)
	}

	endCount := 0
	for {
		<-ch
		endCount++
		if endCount >= con_count {
			break
		}
	}
	te := time.Now()

	dura := te.Sub(ts)

	fmt.Println(dura)
}

func init_redis() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	return c, err
}

func loopWrite(count int, ch chan int) {
	for i := 0; i < count; i++ {
		//	fmt.Print(".")
		rndWrite()
	}
	ch <- 1
}

func rndWrite() {
	conn := dbPool.GetConnection()
	defer dbPool.ReleaseConnection(conn)

	redisCon, ok := conn.(redis.Conn)
	if !ok {
		fmt.Println("convert error")
		return
	}

	key := rand.Int()
	data := rand.Int()

	strKey := strconv.Itoa(key)
	strData := strconv.Itoa(data)

	/*_, err := redisCon.Do("SET", strKey, strData)
	if err != nil {
		fmt.Println(err)
	}
	*/
}
