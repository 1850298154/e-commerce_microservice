package main

import (
	"2501YTC/app/user/biz/model"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudwego/kitex/pkg/klog"
	"os"
	"sync"
)

const (
	ConsumerGroupName = "UserConsumerGroup"
)

type Consumer struct {
	ctx      context.Context
	consumer rocketmq.PushConsumer
	done     chan struct{}
	once     sync.Once
}

func NewConsumer(ctx context.Context) *Consumer {
	namesrvAddr, _ := primitive.NewNamesrvAddr(Endpoint)
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer(namesrvAddr),
		consumer.WithGroupName(ConsumerGroupName),
	)

	return &Consumer{
		ctx:      ctx,
		consumer: c,
		done:     make(chan struct{}),
	}
}

func (c *Consumer) Close() {
	c.once.Do(func() {
		close(c.done)
		if c.consumer != nil {
			_ = c.consumer.Shutdown()
		}
	})
}

func (c *Consumer) processUserLoginLog(msg []byte) (uint, error) {
	log, err := model.FromJSON(msg)
	if err != nil {
		klog.Error("parse login log failed: ", err)
		return 0, err
	}

	fmt.Println(log)

	var q model.UserLoginLogQuery
	logId, err := q.Insert(log)
	if err != nil {
		klog.Error("insert login log failed: ", err)
		return 0, err
	}

	return logId, nil
}

// 订阅消息
func (c *Consumer) Consume(topic string) error {
	err := c.consumer.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			logId, err := c.processUserLoginLog(msg.Body)
			if err != nil {
				fmt.Printf("process messages error: %s\n", err)
				return consumer.ConsumeRetryLater, err
			}
			klog.Infof("Insert login log with id: %s\n", logId)
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	err = c.consumer.Start()
	if err != nil {
		fmt.Printf("start consumer error: %s", err.Error())
		os.Exit(1)
	}
	klog.CtxInfof(c.ctx, "Consumer start successfully, listening topic: %s", topic)
	<-c.done
	klog.CtxInfof(c.ctx, "Consumer Stopped!")
	return nil
}
