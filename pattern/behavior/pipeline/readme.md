> 流水线工作模型在实际开发中非常常见，将一个工作拆分成多个步骤，根据每个步骤的工作量分配不同的人手去处理，从而提供生产效率。

流水线我们又称为pipeline
go语言与pipeline非常契合，每一个步骤根据工作量交给多个worker执行，步骤和步骤之间通过channel进行数据传递。

根据 channel 和 goroutine 的关系，出现两种常见的消费模型 扇入`fan-in` 扇出`fan-out`
在这个模型中，需要格外注意安全退出；核心关注一点：关闭通道前，保证不会再有消息发往次通道。