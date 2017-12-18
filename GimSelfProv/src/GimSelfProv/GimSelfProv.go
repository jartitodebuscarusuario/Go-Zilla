package main

import (
	"flag"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "net/http"
    "log"
)

func main() {
	
  var CONN_HOST string
  var CONN_TYPE string
  var CONN_PORT string
  var CONN_HTTP_PORT string
    
  flag.StringVar(&CONN_HOST, "host", "::1", "Hostname to bind (default: agracia7)")
  flag.StringVar(&CONN_PORT, "port", "8888", "Port to bind (default: 8888)")
  flag.StringVar(&CONN_HTTP_PORT, "http_port", "9090", "Http port to bind (default: 9090)")
  flag.StringVar(&CONN_TYPE, "type", "tcp6", "Connection type (default: tcp)")
  flag.BoolVar(&noprint, "noprint", true, "Supress output (default: false)")
  flag.IntVar(&nvm, "nvm", 1, "Machines per inf (default: 1)")
  flag.StringVar(&encoder, "encoder", "json", "Message encoding type (default: json)")
    
  flag.Parse()
  
  dir, derr := filepath.Abs(filepath.Dir(os.Args[0]))
  if derr != nil {
	  log.Fatal(derr)
  }
    
  //Initialize infmap
  infmap = &Infmap{ Data: map[string]*Infdata{} }
  //infmap.Data = map[string]*Infdata{}
	
  //Initialize evalsp (how we evaluate sp)
  evalsp = map[string]int{
	"cpu": 0,
	"mem": 1,
  }
    
  //Load default conf from file
  defconf = readConf(dir + "\\config.json")
	
  fmt.Println(defconf)
	
  go func() {

    // Listen for incoming connections.
    tcpAddr, err := net.ResolveTCPAddr(CONN_TYPE, "[" + CONN_HOST + "]:" + CONN_PORT)
    //l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    l, err := net.ListenTCP(CONN_TYPE, tcpAddr)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        //atomic.AddInt32(&timeCounter,1)
        defer conn.Close()
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
  }()
  
  http.HandleFunc("/", regInfid) // set router
  err := http.ListenAndServe(":" + CONN_HTTP_PORT, nil) // set listen port
  if err != nil {
        log.Fatal("ListenAndServe: ", err)
  }
}