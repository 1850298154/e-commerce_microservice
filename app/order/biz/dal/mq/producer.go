package mq

import (
	"encoding/json"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

const (
	DelayExchange      = "order.delay.exchange"
	DelayQueue         = "order.delay.queue"
	DeadLetterExchange = "order.dlx.exchange"
	DeadLetterQueue    = "order.dlx.queue"
)

var ProducerInstance *Producer // ProducerInstance 生产者实例

// OrderMessage 订单消息结构体
type OrderMessage struct {
	OrderID string `json:"order_id"`
}

// Producer 生产者结构体
type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewProducer 创建生产者实例
func NewProducer(orderTimeout int) (*Producer, error) {
	channel, err := RabbitMQConn.Channel()
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, err: %v", err)
		return nil, err
	}

	producer := &Producer{
		conn:    RabbitMQConn,
		channel: channel,
	}

	err = producer.initializeQueue(orderTimeout)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法初始化队列, err: %v", err)
		return nil, err
	}
	klog.Info("RabbitMQ Producer 初始化成功")
	return producer, nil
}

// initializeQueue 初始化交换机和队列
func (p *Producer) initializeQueue(orderTimeout int) error {
	// 声明死信交换机
	err := p.channel.ExchangeDeclare(
		DeadLetterExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法初始化死信交换机, err: %v", err)
		return err
	}

	// 声明死信队列
	_, err = p.channel.QueueDeclare(
		DeadLetterQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法初始化死信队列, err: %v", err)
		return err
	}

	// 绑定死信队列到死信交换机
	err = p.channel.QueueBind(
		DeadLetterQueue,
		DeadLetterQueue,
		DeadLetterExchange,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法绑定死信队列到死信交换机,err: %v", err)
		return err
	}

	// 声明延迟交换机
	err = p.channel.ExchangeDeclare(
		DelayExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法初始化延迟交换机,err: %v", err)
		return err
	}

	// 声明延迟队列，设置消息过期后转发到死信交换机
	args := amqp.Table{
		"x-dead-letter-exchange":    DeadLetterExchange,
		"x-dead-letter-routing-key": DeadLetterQueue,
		"x-message-ttl":             orderTimeout,
	}
	_, err = p.channel.QueueDeclare(
		DelayQueue,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法初始化延迟队列,err: %v", err)
		return err
	}

	// 绑定延迟队列到延迟交换机
	err = p.channel.QueueBind(
		DelayQueue,
		DelayQueue,
		DelayExchange,
		false,
		nil,
	)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 初始化失败, 无法绑定延迟队列到延迟交换机,err: %v", err)
		return err
	}
	return nil
}

// Stop 关闭连接
func (p *Producer) Stop() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
	klog.Info("RabbitMQ Producer 关闭成功")
}

// SendDelayMessage 发送延迟消息
func (p *Producer) SendDelayMessage(orderID string) error {
	message := OrderMessage{
		OrderID: orderID,
	}

	body, err := json.Marshal(message)
	if err != nil {
		klog.Errorf("RabbitMQ Producer 发送消息失败, err: %v", err)
		return err
	}

	return p.channel.Publish(
		DelayExchange,
		DelayQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
