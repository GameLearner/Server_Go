package Network

import (
    "Server/PBProto"
    "github.com/golang/protobuf/proto"
    "bytes"
    "encoding/binary"
    "fmt"
)

const (
    TestID = iota;

)

//ProtocolInterface protocol interface
type ProtocolInterface interface {
    Marshal() ([]byte, error)
    UnMarshal(data []byte) error
}

//Protocol protocol 
type Protocol struct {
    ID int32;
    Packet interface{}
}

var packetMap map[int32]interface{}

func init()  {
    packetMap = make(map[int32]interface{})
    packetMap[TestID] = new(PBProto.Test);
    
}

//NewProtocol create Protocol
func NewProtocol(id int32) *Protocol {
    packet, ok := packetMap[id];
    if !ok {
        return nil;
    }
    protocol := new(Protocol);
    protocol.ID = id;
    protocol.Packet = packet;
    
    return protocol;
}

//Marshal convert the protocol to bytes
func (protocol *Protocol)Marshal() ([]byte, error) {
    buff := new(bytes.Buffer);
    binary.Write(buff, binary.BigEndian, protocol.ID);
    ms, ok := protocol.Packet.(proto.Marshaler);
    if !ok {
		return nil, fmt.Errorf("protocol error not valid protobuff");
	}
    data, err := ms.Marshal()
    if nil != err {
        return nil, fmt.Errorf("Packet Marshal Error");
    }
    buff.Write(data);
    return buff.Bytes(), nil
}

//UnMarshal Protocol's Unmarshal UnSerialize protocol from bytes
func (protocol *Protocol) UnMarshal(data []byte) error {
    if(len(data) < 4) {
        //不完整的数据 待下次再读
        return fmt.Errorf("incomplete data."); 
    }
    
    idSplit := data[:4]
    packSplit := data[4:]
    
    buf := bytes.NewReader(idSplit);
    err := binary.Read(buf, binary.LittleEndian, &protocol.ID);
    if nil != err {
        return fmt.Errorf("Packet Id UnMarshal Error");
    }
    
    ms, ok := protocol.Packet.(proto.Unmarshaler);
    if !ok {
        return fmt.Errorf("Packet data error");
    }
    ms.Unmarshal(packSplit);
    return nil
}

