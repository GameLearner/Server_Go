package Network

import (
    "net"
    "fmt"
	"time"
    "Server/Util"
    "bytes"
	"sync/atomic"
)


//SendBuffSize send buff size.
const (
	MAXSENDNUM = 128
	MAXRECVNUM = 512

	SENDBUFFSIZE = 1024 * 64
    RECVBUFFSIZE = 1024 * 64
)

//Session : net connection.
type Session struct {
	ID          int
	ServerID    int
	conn        net.Conn
	sPacketBuff chan interface{} // sync send buff
	rPacketBuff chan interface{} // sync recv buff
	recvCh      chan int         // sync recv
	recvData	[]byte
	validFlag	int32			 // 
}

//NewSession :Create new Session.
func NewSession(conn net.Conn) (*Session, error) {
	session := new(Session)
	session.ID = Util.GetUniqID()
	session.conn = conn

	session.sPacketBuff = make(chan interface{}, MAXSENDNUM)
    session.rPacketBuff = make(chan interface{}, MAXRECVNUM)
    session.recvCh = make(chan int)
	
	session.validFlag = -1;
    
	return session, nil
}

//Run loop while conn valid, process input and output
//
func (session *Session) Run() {
	session.validFlag = 1;
	//recv
	go func ()  {
		session.recv();
	}()
	//send
	go func ()  {
		session.send();
	}();
	//process
	for {
		session.handle();
	}
}

func (session *Session) handle() {

}

//recv do read data from socket
func (session *Session) recv() {
	defer session.Close();
	
	for{
		num, err := 
	}
}

//Close close session
func (session *Session) Close()  {
	if atomic.CompareAndSwapInt32(&session.validFlag, 1, -1) {
		session.conn.Close();
	}
}

func (session *Session) doSend(data []byte) (error)  {
	_, err := session.conn.Write(data)
	return err;
}

func (session *Session) doSendBuff(buff *bytes.Buffer) (error)  {
	err := session.doSend(buff.Bytes())
	buff.Reset();
	return err;
}

//send when buff overflow or time out do send to the socket
//reduce the call times of system send
func (session *Session) send() {
	defer session.Close();
	
    var buff bytes.Buffer;
	tickCh := time.Tick(20 * time.Millisecond)
	for {
		select {
			case packet := <-session.sPacketBuff:
			{
				data, err := MarshalPacket(packet)
				if nil != err {
					fmt.Println("error not protobuf packet")
					continue
				}
				if (buff.Len() + len(data)) >= SENDBUFFSIZE {
					err = session.doSendBuff(&buff)
				} 
                buff.Write(data);
                 
			}
			case <-tickCh:
			{
				if buff.Len() > 0 {
					err := session.doSendBuff(&buff)
					if nil != err {
						
					}
				}
			}
		}
	}
}

//SendPacket send packet to the send buff of session wait send to do really send, 
//if buff overflow return err
func (session *Session) SendPacket(packet interface{}) error {
	select {
	case session.sPacketBuff <- packet:
		{

		}
	default:
		{
			return fmt.Errorf("session Id = %d, send buff overflow ", session.ID)
		}
	}
	return nil
}

//MarshalPacket use protobuf's Marshaler interface to Marshal packet
func MarshalPacket(packet interface{}) ([]byte, error) {
	if ms, ok := packet.(ProtocolInterface); ok {
		return ms.Marshal()
	}
	return nil, fmt.Errorf("protocol error not proto buffer")
}
