package main

import (
    "net"
    "log"
    "Server/Network"
    "time"
    "sync/atomic"
)

var sessionMap map[int]*Network.Session;
var count int32
func onSessionClose(session *Network.Session)  {
    atomic.AddInt32(&count, -1)
}


func main() {
    
    listner, err := net.Listen("tcp", ":9999");
    if nil != err {
        log.Fatalln(err);
    }
    
    sessionMap = make(map[int]*Network.Session);
    count = 0;

    tickCh := time.Tick(time.Second * 5)

    for {
        conn, err := listner.Accept();
        if nil != err {
            log.Println(err);
            continue;
        }
        count++
        go func ()  {
            //log.Printf("new connection %s, connections = %d", conn.RemoteAddr().String(), count)
            session, _ := Network.NewSession(conn, nil)
            session.Run();
        }()
        //sessionMap[session.ID] = session;
        select {
            case <-tickCh:
            {
                log.Printf("alive connections = %d", count)
            }
            default:
            {
                ;
            }
        }
    }
}

