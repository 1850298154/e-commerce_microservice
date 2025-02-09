package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"2501YTC/app/order/biz/model"
	"2501YTC/app/order/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

// Consumer 消费者结构体
type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	done       chan struct{}
	maxRetries int
	DB         *gorm.DB
	orderQuery model.OrderQuery
}

// NewConsumer 创建消费者
func NewConsumer(db *gorm.DB) (*Consumer, error) {
	channel, err := RabbitMQConn.Channel()
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		conn:       RabbitMQConn,
		channel:    channel,
		done:       make(chan struct{}),
		maxRetries: conf.GetConf().RabbitMQ.MaxRetries,
		DB:         db,
		orderQuery: model.NewOrderQuery(context.Background(), db),
	}

	klog.Info("RabbitMQ Consumer 初始化成功")
	return consumer, nil
}

// handleOrderWithOptimisticLock 使用乐观锁处理订单
func (c *Consumer) handleOrderWithOptimisticLock(orderMsg OrderMessage) error {
	var err error
	klog.Infof("正在处理订单: %d", orderMsg.OrderID)

	for i := 0; i < c.maxRetries; i++ {
		version, orderState, err := c.orderQuery.GetOrderVersionAndState(orderMsg.OrderID)
		if err != nil {
			klog.Errorf("获取订单版本号失败: %v", err)
			return err
		}
		// 如果订单状态不是已下单，不处理 -> 已被其他消费者处理过了 || 订单已取消、已完成
		if orderState != model.OrderStatePlaced {
			klog.Infof("订单 %d 状态不是已下单，不处理", orderMsg.OrderID)
			return nil
		}
		err = c.orderQuery.CancelOrderWithVersion(orderMsg.OrderID, model.CancelTypeTimeout, int32(time.Now().Unix()), version)
		if err == nil {
			klog.Infof("订单 %d 处理成功", orderMsg.OrderID)
			return nil
		}
		// 如果是乐观锁冲突，继续重试
		klog.Warnf("乐观锁冲突，正在重试 (%d/%d)", i+1, c.maxRetries)
	}

	return fmt.Errorf("达到最大重试次数，处理订单失败: %v", err)
}

// Consume 消费者消费消息
func (c *Consumer) Consume() error {
	// 设置预取计数，控制消费者同时处理的消息数量
	err := c.channel.Qos(1, 0, false)
	if err != nil {
		klog.Errorf("设置RabbitMQ Consumer预取计数失败: %v", err)
		return err
	}

	msgs, err := c.channel.Consume(
		DeadLetterQueue,
		"",    // 消费者标识
		false, // 自动确认
		false, // 独占
		false, // 非阻塞
		false, // 等待服务器确认
		nil,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Consumer start failed: %v", err)
		return err
	}

	go func() {
		for msg := range msgs {
			var orderMsg OrderMessage
			if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
				klog.Errorf("Consumer解析订单消息失败: %v", err)
				_ = msg.Nack(false, false)
				continue
			}

			// 使用乐观锁处理订单
			err := c.handleOrderWithOptimisticLock(orderMsg)
			if err != nil {
				klog.Errorf("Consumer处理订单失败: %v", err)
				if err == gorm.ErrRecordNotFound {
					_ = msg.Ack(false) // 订单不存在，直接确认
					continue
				}
				_ = msg.Nack(false, true) // 重新入队
				continue
			}

			_ = msg.Ack(false)
		}
	}()

	klog.Infof("Consumer start successfully, listening dead letter queue: %s", DeadLetterQueue)
	<-c.done
	klog.Info("Consumer Stopped!")
	return nil
}

// Stop 停止消费者
func (c *Consumer) Stop() {
	close(c.done)
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
