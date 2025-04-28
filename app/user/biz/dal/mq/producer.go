package main

import (
	"2501YTC/app/user/biz/model"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/cloudwego/kitex/pkg/klog"
	"os"
	"strconv"
)

const (
	Topic             = "LoginLogTopic"
	ProducerGroupName = "UserProducer"
	Endpoint          = "127.0.0.1:9876"
)

type Producer struct {
	ctx      context.Context
	producer rocketmq.Producer
}

func NewProducer(ctx context.Context) *Producer {
	namesrvAddr, _ := primitive.NewNamesrvAddr(Endpoint)
	p, _ := rocketmq.NewProducer(
		producer.WithNameServer(namesrvAddr),
		producer.WithRetry(2),
		//producer.WithGroupName(ProducerGroupName)
	)

	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	return &Producer{
		ctx:      ctx,
		producer: p,
	}
}

func (p *Producer) Close() {
	_ = p.producer.Shutdown()
}

func (p *Producer) SendMessage(log *model.UserLoginLog) error {
	logmsg, _ := log.ToJSON()
	msgbody := primitive.NewMessage(Topic, logmsg)

	if log.LoginStatus == model.LoginFailed {
		msgbody.WithTag("fail")
	} else {
		msgbody.WithTag("success")
	}
	userId := strconv.Itoa(int(log.UserID))
	msgbody.WithKeys([]string{userId})
	err := p.producer.SendOneWay(p.ctx, msgbody)
	if err != nil {
		klog.Errorf("send user login log error: %s", err.Error())
		return err
	}
	klog.Infof("send user login log success")
	return nil
}
