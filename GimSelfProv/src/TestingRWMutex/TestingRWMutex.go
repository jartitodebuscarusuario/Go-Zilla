package main

import (
	"flag"
    "fmt"
    "net"
    "os"
    //"sync"
    //"sync/atomic"
    "time"
    //"strconv"
    //"math/rand"
    //"strconv"
    "encoding/json"
    //"encoding/gob"
    //"bytes"
    //"github.com/satori/go.uuid"
    "net/http"
    //"strings"
    "log"
    "io/ioutil"
)

//Global vars
var noprint bool
//var infmap Infmap
var infmap *Infmap
var evalsp map[string]int
var defconf map[string]interface{}
var nvm int
var encoder string
//var timeCounter int32
//var timeElapsed time.Duration

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
    
  //Initialize infmap
  infmap = &Infmap{ Data: map[string]*Infdata{} }
  //infmap.Data = map[string]*Infdata{}
	
  //Initialize evalsp (how we evaluate sp)
  evalsp = map[string]int{
	"cpu": 0,
	"mem": 1,
  }
    
  //Load default conf from file
  defconf = readConf("bin\\config.json")
	
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

func readConf(file string) map[string]interface{} {
	conf := make(map[string]interface{})
	dat, err := ioutil.ReadFile(file)
	var fconf interface{}
	err = json.Unmarshal(dat, &fconf)
		if err != nil {
		  fmt.Println("error:", err)
	}
	conf = fconf.(map[string]interface{})
	return conf
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  	
  //defer conn.Close()
  // Make a buffer to hold incoming data.
  //dec := gob.NewDecoder(conn)
  //var dec interface{}
  dec := json.NewDecoder(conn)
  
  //decType := dec.(type)
	
  //for dec.(*json.Decoder).More() {
  for dec.More() {
  message := &Tmessage{}
  //dec.(*json.Decoder).Decode(message)
  dec.Decode(message)
  // Write received data to disk
  //f.Write(request[:read_len])
  
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received " + message.Infid + " " + message.Vmid))

  //TEST
  start := time.Now()
  
  checkInfid(message.Infid)
  checkVmid(message.Infid, message.Vmid)
  addData2InfidVmid(message.Infid, message.Vmid, message.Data)
  evaluatesp(message.Infid)
  wprint("Added data", message.Data, "to vmid", message.Vmid, "in infid", message.Infid)
  
  wprint("timeCounter: ", time.Now().Sub(start))
  //TEST  
  
  }
  // Close the connection when you're done with it.
  conn.Close()
  
//  infmap.RLock()
//  defer infmap.RUnlock()
//  for key, value := range infmap.Data {
//  	fmt.Println("key: ",key," =>")
//    for key, value := range value.Data {
//      fmt.Println("key: ",key," =>")
//      for key, value := range value.Data {
//	    fmt.Println("Key:", key, "Value:", value)
//      }
//    }
//  }
//  fmt.Println("====================================")
  
}

//func readInf (infid string) (infdata map[string]*Vmdata) {
//	infmap.RLock()
//	infmap.Data[infid].RLock()
//	infdata = infmap.Data[infid].Data
//	infmap.Data[infid].RUnlock()
//	infmap.RUnlock()
//	return infdata
//}

func (mapinf *Infmap)infidRLock(idinf string) {
	mapinf.RLock()
	mapinf.Data[idinf].RLock()
}

func (mapinf *Infmap)infidRUnlock(idinf string) {
	mapinf.Data[idinf].RUnlock()
	mapinf.RUnlock()
}

func (mapinf *Infmap)AddVmidData(idinf string, idvm string, vmdata map[string]int) {
	mapinf.RLock()
	defer mapinf.RUnlock()
	mapinf.Data[idinf].RLock()
	defer mapinf.Data[idinf].RUnlock()
	mapinf.Data[idinf].Data[idvm].Lock()
	defer mapinf.Data[idinf].Data[idvm].Unlock()
	mapinf.Data[idinf].Data[idvm].Data = vmdata	
}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}