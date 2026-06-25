package mq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 消息队列
type MQ struct {
	conn *amqp.Connection
}

// MessageHandler 消费单条消息的处理函数。
// 返回 nil 表示处理成功；返回 error 表示失败（由库决定 Nack/重试策略）。
type MessageHandler func(msg amqp.Delivery) error

// New 创建 MQ 实例
func New(url string) (*MQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MQ{conn: conn}, nil
}

// Publish 快捷发布消息
func (m *MQ) Publish(ctx context.Context, queue string, body []byte) error {
	return m.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// RabbitMQ 本身不直接支持延迟队列，但可以通过使用带有 "x-delayed-message" 插件的 exchange 或基于消息的死信（TTL+DLX）机制实现延迟发送。
// 下面是一个通过设置消息属性（TTL+DLX）实现的延迟消息发布方法示例。

// PublishWithDelay 通过 TTL + 死信队列方式实现延迟消息投递（需确保死信队列已正确配置）
// 死信队列(DLX)配置通常包括以下步骤：
//
// 1. 创建一个“死信交换机”(Dead Letter Exchange, DLX)，可以是 direct、topic 类型，例如：`my-dlx`。
// 2. 创建一个“死信队列”，比如 `my-dlx-queue`，并绑定到上面的 DLX。
// 3. 在你的原业务队列（如 queue=`my-business-queue`）声明时，增加如下参数：
//    arguments := amqp.Table{
//        "x-dead-letter-exchange":    "my-dlx",     // 死信交换机名称
//        "x-dead-letter-routing-key": "my-dlx-key", // 死信路由 key（可选，通常与 DLX 绑定的 routing-key 保持一致，如果 DLX 是 direct/需要）。
//    }
//    这样，如果业务队列中的消息“过期”(TTL 到期) 或被拒绝（Nack 且 requeue=false），就会自动路由到上面的 DLX，然后进入死信队列。
//
// 4. 设置 TTL：
//    - 队列级 TTL：给 QueueDeclare 时加 "x-message-ttl": xxx（对整个队列所有消息生效）。
//    - 消息级 TTL：Publish 消息时给 Expiration 字段赋值（如上方 PublishWithDelay）。
//
// 5. 总体代码结构样例：
// ch.QueueDeclare("my-dlx-queue", true, false, false, false, nil)
// ch.ExchangeDeclare("my-dlx", "direct", true, false, false, false, nil)
// ch.QueueBind("my-dlx-queue", "my-dlx-key", "my-dlx", false, nil)
// ch.QueueDeclare("my-business-queue", true, false, false, false, amqp.Table{
//     "x-dead-letter-exchange":    "my-dlx",
//     "x-dead-letter-routing-key": "my-dlx-key",
// })
// 消息 PublishWithDelay 时传递 Expiration，实现消息单独延迟，延迟时间到即路由至死信交换机，由此达到延迟消费效果！
//
// 总结：
// - 死信队列负责接收“到期”或“被拒绝”的消息。
// - 原队列要用 "x-dead-letter-exchange" 指明死信交换机。
// - 死信交换机把死信分发到死信队列。
// - “延迟队列”效果通常是“业务队列 + 死信转移”完成。

func (m *MQ) PublishWithDelay(ctx context.Context, queue string, body []byte, delayMs int) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 设置消息的过期时间
	headers := amqp.Table{}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Expiration:  fmt.Sprintf("%d", delayMs), // 单位：毫秒
		Headers:     headers,
	}

	return ch.PublishWithContext(ctx, "", queue, false, false, msg)
}

// PublishWithContext 发布消息
func (m *MQ) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

// StartRabbitMQConsumer 启动 RabbitMQ 消费者
// nack 是否需要重新入列一次？
func (m *MQ) StartRabbitMQConsumer(nack bool, queue string, handler MessageHandler) error {
	ch, err := m.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err = ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}

	// 消费队列
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for d := range msgs {
		// 处理失败时重新入队（可按需改成死信队列策略）
		if err := handler(d); err != nil && nack {
			_ = d.Nack(false, true)
			continue
		}
		// 手动确认
		_ = d.Ack(false)
	}

	return nil
}
