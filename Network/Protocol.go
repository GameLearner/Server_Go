package Network

import (
    "Server/PBProto"
    "github.com/golang/protobuf/proto"
    "bytes"
    "encoding/binary"
    "fmt"
)

const (
    LoginID = iota;
)

//ProtocolInterface protocol interface
type ProtocolInterface interface {
    Marshal() ([]byte, error)
    UnMarshal(data []byte) (int, error)
}

//Protocol protocol 
type Protocol struct {
    ID int32;
    Packet interface{}
}

var packetMap map[int32]interface{}

func init()  {
    packetMap = make(map[int32]interface{})
    packetMap[LoginID] = new(PBProto.Login);   
}

func NewPBPacket(id int32) interface{} {
    packet, ok := packetMap[id];
    if !ok {
        return nil;
    }
    ms, _ := packet.(proto.Message);
    
    return proto.Clone(ms)
}

//NewProtocol create Protocol
func NewProtocol(id int32) *Protocol {

    protocol := new(Protocol);
    protocol.ID = id;
    protocol.Packet = NewPBPacket(id);
    
    return protocol;
}

//Marshal convert the protocol to bytes
func (protocol *Protocol)Marshal() ([]byte, error) {
    buff := new(bytes.Buffer);
    binary.Write(buff, binary.BigEndian, protocol.ID);
    ms, ok := protocol.Packet.(proto.Message);
    if !ok {
		return nil, fmt.Errorf("Protocol error not valid protobuff");
	}
    data, err := proto.Marshal(ms)
    if nil != err {
        return nil, fmt.Errorf("Packet Marshal Error");
    }
    buff.Write(data);
    return buff.Bytes(), nil
}

//UnMarshal Protocol's Unmarshal UnSerialize protocol from bytes
func (protocol *Protocol) UnMarshal(data []byte) (int, error) {
    if(len(data) < 4) {
        //不完整的数据 待下次再读
        return 0, fmt.Errorf("incomplete data."); 
    }
    
    idSplit := data[:4]
    packSplit := data[4:]
    
    buf := bytes.NewReader(idSplit);
    err := binary.Read(buf, binary.LittleEndian, &protocol.ID);
    if nil != err {
        return 0, fmt.Errorf("Packet Id UnMarshal Error");
    }
    
    protocol.Packet = NewPBPacket(protocol.ID)
    
    if nil == protocol.Packet {
        return -1, fmt.Errorf("Invalid packetId, need Close client??? id = %d, data = %v", protocol.ID, data)
    }
    
    ms, _ := protocol.Packet.(proto.Message);
    

    err = proto.Unmarshal(packSplit, ms)
    
    if nil != err {
        return 0, fmt.Errorf("PBMessage Unmarshal Error!!! incomplete packet or need close client")
    }
    
    msLen := proto.Size(ms)

    packetLen := msLen + 4

    return packetLen, nil
}

