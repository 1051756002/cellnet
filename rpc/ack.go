package rpc

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
)

type Response interface {
	Feedback(interface{})
}

type response struct {
	ses cellnet.Session
	req *coredef.RemoteCallREQ
}

func (self *response) Feedback(msg interface{}) {

	pkt := cellnet.BuildPacket(msg)

	self.ses.Send(&coredef.RemoteCallACK{
		MsgID:  proto.Uint32(pkt.MsgID),
		Data:   pkt.Data,
		CallID: proto.Int64(self.req.GetCallID()),
	})
}

func (self *response) ContextID() int {
	return int(self.req.GetMsgID())
}

func InstallServer(p cellnet.Peer) {

	// 服务端
	socket.RegisterSessionMessage(p, coredef.RemoteCallREQ{}, func(ses cellnet.Session, content interface{}) {
		msg := content.(*coredef.RemoteCallREQ)

		// 客户端发过来的请求消息注入到回调中
		p.CallData(&response{
			ses: ses,
			req: msg,
		})

	})

}

// 注册连接消息
func RegisterMessage(eq cellnet.EventQueue, msgIns interface{}, userHandler func(Response, interface{})) {

	msgType := reflect.TypeOf(msgIns)

	msgName := msgType.String()

	msgID := cellnet.Name2ID(msgName)

	// 将消息注册到mapper中, 提供反射用
	socket.MapNameID(msgName, msgID)

	eq.RegisterCallback(msgID, func(data interface{}) {

		if ev, ok := data.(*response); ok {

			rawMsg, err := cellnet.ParsePacket(&cellnet.Packet{
				MsgID: ev.req.GetMsgID(),
				Data:  ev.req.Data,
			}, msgType)

			if err != nil {
				log.Printf("[cellnet] unmarshaling error:\n", err)
				return
			}

			userHandler(ev, rawMsg)

		}

	})

}
