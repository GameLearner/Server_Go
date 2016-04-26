package main

import (
    "net"
    "log"
    "Server/Network"
)

var sessionMap map[int]Network.Session;

func main() {
    
    listner, err := net.Listen("tcp", ":9999");
    if nil != err {
        log.Fatalln(err);
    }
    
    var sid int;
    sessionMap = make(map[int]Network.Session);
    
    for {
        conn, err := listner.Accept();
        if nil != err {
            log.Println(err);
            continue;
        }
        log.Println("new connection " + conn.RemoteAddr().String())
        
        var session Network.Session;

        sid++;
        
        sessionMap[session.ID] = session;
        
        go handleConnect(conn);        
    }
}


func handleConnect(conn net.Conn) {
    defer conn.Close();
    for{
        
    }
}