package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	producer := NewProducer(ctx)

	//log := model.UserLoginLog{
	//	UserID:      12345,
	//	LoginTime:   time.Now(),
	//	IPAddress:   "192.168.1.1",
	//	DeviceType:  "Mobile",
	//	LoginStatus: 1,
	//	Location:    "New York, USA",
	//}
	//err := producer.SendMessage(&log)
	//if err != nil {
	//	fmt.Println(err)
	//}

	consumer := NewConsumer(ctx)
	err := consumer.Consume("LoginLogTopic")
	if err != nil {
		fmt.Println(err)
	}

	defer consumer.Close()
	defer producer.Close()

}
