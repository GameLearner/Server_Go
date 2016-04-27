package main

import (
    "net"
    "log"
    "Server/Network"
)

var sessionMap map[int]*Network.Session;

func main() {
    
    listner, err := net.Listen("tcp", ":9999");
    if nil != err {
        log.Fatalln(err);
    }
    
    sessionMap = make(map[int]*Network.Session);
    count := 0;
    for {
        conn, err := listner.Accept();
        if nil != err {
            log.Println(err);
            continue;
        }
        count++
        go func ()  {
            log.Printf("new connection %s, connections = %d", conn.RemoteAddr().String(), count)

            session, _ := Network.NewSession(conn, nil)
            session.Run();
        }()
        
        //sessionMap[session.ID] = session;
        
    }
}

