package imgo

import (
	"context"
	"net"
	"time"
)

type Server interface {
	SetAcceptor(Acceptor)               // 用于设置一个Acceptor
	SetMessageListener(MessageListener) // 设置一个消息监听器
	SetStateListener(StateListener)     // 用于设置一个状态监听器，将连接断开的事件上报给业务层。
	SetReadWait(time.Duration)          // 用于设置连接读超时，用于控制心跳逻辑。
	SetChannelMap(ChannelMap)           // 设置一个连接管理器，Server在内部会自动管理连接的生命周期。

	Start() error
	Push(string, []byte) error
	Shutdown(context.Context) error
}

// 连接接收器
type Acceptor interface {
	// 让上层业务处理握手相关工作，比如在网关层可以在这个接收者中实现对用户身份的认证。
	// 返回两个参数 channelID: 唯一通道标识；error: 如果返回一个error连接就会被断开。
	Accept(Conn, time.Duration) (string, error)
}

// 对net.Conn进行封装，把读与写的操作封装到连接中。
// ReadFrame 与 WriteFrame，完成对websocket/tcp两种协议的封包与拆包逻辑的包装。
type Conn interface {
	net.Conn
	ReadFrame() (Frame, error)
	WriteFrame(OpCode, []byte) error
	Flush() error
}

// 监听消息
type MessageListener interface {
	// 收到消息回调
	Receive(Agent, []byte)
}

// 消息发送方
type Agent interface {
	ID() string        // 返回连接的channelID
	Push([]byte) error // 用于上层业务返回消息
}

// 状态监听器
type StateListener interface {
	// 连接断开回调
	Disconnect(string) error
}

// 客户端
type Client interface {
	ID() string
	Name() string

	Connect(string) error // 连接服务器。
	SetDialer(Dialer)     // 设置一个拨号器，这个方法会在Connect中被调用，完成连接的建立和握手。
	Send([]byte) error    // 发送消息到服务端。
	Read() (Frame, error) // 读取一帧数据，返回Frame。
	Close()               // 断开连接，退出。
}

type Dialer interface {
	DialAndHandshake(DialerContext) (net.Conn, error)
}

// 拨号信息
type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}

// OpCode
type OpCode byte

// OpCode 类型
const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)

// 解决底层封包与拆包问题
type Frame interface {
	SetOpCode(OpCode)
	GetOpCode() OpCode
	SetPayload([]byte)
	GetPayLoad() []byte
}
