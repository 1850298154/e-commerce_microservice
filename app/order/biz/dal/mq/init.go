package mq

import (
	"fmt"

	"2501YTC/app/order/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

var (
	RabbitMQConn *amqp.Connection
)

func Init() {
	connString := fmt.Sprintf("%s://%s:%s@%s:%d/", conf.GetConf().RabbitMQ.MQ, conf.GetConf().RabbitMQ.User, conf.GetConf().RabbitMQ.Password, conf.GetConf().RabbitMQ.Host, conf.GetConf().RabbitMQ.Port)
	conn, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	RabbitMQConn = conn
	klog.Infof("RabbitMQ 初始化成功, DSN: %s", connString)
}
