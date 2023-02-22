package metrics

import "time"

type CollectResources struct {
	NodeStatus MessageType
	PodStatus  MessageType
	ConfigMap  MessageType
	Secret     MessageType
	QueryNode  MessageType
	UpdateNode MessageType
	PodDelete  MessageType

	ReceiveCo ReceiveType //upstream的接受情况
	SendError SendType    //downstream发送情况
	SendOK    SendType    //downstream发送情况

}

type MessageType struct {
	GetMessageErrorCount int           //累计处理消息中，获取相关数据的错误数
	OperateErrorCount    int           //累计处理消息中，执行操作的错误数
	ChanLen              int           //当前队列长度
	ReceiveCount         int           //收到的消息数
	ConsumeTime          time.Duration //收到的消息数，进行处理的耗时
}

// SendType 发送消息 -- downstream ==统计发送数量
type SendType struct {
	SyncPodCount  int
	SyncConfigMap int
	SyncSecret    int
	SyncEdgeNodes int
}

// ReceiveType 接受数据  -- upstream
type ReceiveType struct {
	ErrorCount  int
	TotalCount  int
	ConsumeTime time.Duration
}
