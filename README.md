# kit

### RabbitMQ 使用方法

**实例化**
```
// url 如："amqp://admin:Bluetti2026.06.15@192.168.23.107:5672/"
mq, err := mq.New(c.RabbitMQ.Url) 
if err != nil {
	logx.Errorf("fail to new mq %s", err.Error())
}

// 创建消费者
go mq.StartRabbitMQConsumer(true, types.QueueName, ctx.Callback)
```

**消费者主体**
```
func DeliveryServiceConsumer(msg amqp.Delivery) error {
	fmt.Println(string(msg.Body))
	return nil
}
```

**生产者**
```
l.svcCtx.MQ.Publish(l.ctx, types.QueueName, []byte("test"))
```
