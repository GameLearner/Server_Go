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
    
    for {
        conn, err := listner.Accept();
        if nil != err {
            log.Println(err);
            continue;
        }
        log.Println("new connection " + conn.RemoteAddr().String())
        
        session, _ := Network.NewSession(conn)

        go session.Run();
        
        sessionMap[session.ID] = session;
        
    }
}

