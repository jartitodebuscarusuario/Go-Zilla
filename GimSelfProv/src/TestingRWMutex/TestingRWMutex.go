package main

import (
	"flag"
    "fmt"
    "net"
    "os"
    "sync"
    //"time"
    "strconv"
    "math/rand"
    //"strconv"
    //"encoding/json"
    //"bytes"
    "github.com/satori/go.uuid"
    "encoding/gob"
)

//Type to store monitored data from remote instances
type Tmessage struct {
	Infid string
	Vmid string
	Data map[string]int
}

//Type to store monitored data from remote instances in global memory
type Vmdata struct {
	sync.RWMutex
	Data map[string]int
}

//Type to store remote instances in global memmory
type Infdata struct {
	sync.RWMutex
	Data map[string]*Vmdata
}

//Type to store remote infraestructures
type Infmap struct {
	sync.RWMutex
	Data map[string]*Infdata
}

//Global vars
var noprint bool
var infmap Infmap

func main() {

    var CONN_HOST string
    var CONN_TYPE string
    var CONN_PORT string
    
    var nvm int
    
    flag.StringVar(&CONN_HOST, "host", "agracia7", "Hostname to bind (default: agracia7)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Hostname to bind (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", false, "Supress output (default: false)")
    flag.IntVar(&nvm, "nvm", 1, "Machines per inf (default: 1)")
    
    flag.Parse()
    
    // Open file to write received data
	//f, err = os.Create("D:\\agracia\\Documents\\g\\perl\\GoReport_server.txt")

	const nTmessage = 1000000
	const ninfids = 100
	const nvmids = 5
	var prueba [nTmessage]*Tmessage
	var infids [ninfids]string
	var vmids [nvmids]string
	for i := 0; i < nvmids; i++ {
		vmids[i] = strconv.Itoa(rand.Intn(nvmids))
	}
	for i := 0; i < ninfids; i++ {
		infids[i] = fmt.Sprintf("%s", uuid.NewV1())
	}

	for i := 0; i < nTmessage; i++ {
		prueba[i] = &Tmessage{
			Infid: infids[rand.Intn(ninfids)],
			Vmid: vmids[rand.Intn(5)],
			Data: map[string]int{
				"cpu": rand.Intn(100),
				"mem": rand.Intn(100),
			},
		}
	}
	
	//Initialize infmap
	infmap.Data = map[string]*Infdata{}
	
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
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
        // Handle connections in a new goroutine.
        handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

  // Make a buffer to hold incoming data.
  dec := gob.NewDecoder(conn)
  message := &Tmessage{}
  dec.Decode(message)

  // Write received data to disk
  //f.Write(request[:read_len])
  
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received " + message.Infid + " " + message.Vmid))
  
  // Close the connection when you're done with it.
  conn.Close()
  
  infmap.addData(message)
  
  for key, value := range infmap.Data {
  	wprint("key: ",key," =>")
    for key, value := range value.Data {
      wprint("key: ",key," =>")
      for key, value := range value.Data {
	    wprint("Key:", key, "Value:", value)
      }
      
    }
  }
  
}

func (infmap Infmap) addData (message *Tmessage) {

  infmap.RLock()
  defer infmap.RUnlock()
   
  if _, ok := infmap.Data[message.Infid]; !ok { //Initialize message.Infid with message.Vmid 
	infmap.RUnlock()
	defer infmap.RLock()
	infmap.Lock()
	defer infmap.Unlock()
	infmap.Data[message.Infid] = &Infdata{ 
	  Data: map[string]*Vmdata{ 
	    message.Vmid: &Vmdata{ 
	      Data: map[string]int{},
	    },
	  },
	}
  } else { //Initialize message.Vmid in message.Infid
	 infmap.infidRLock(message.Infid)
	 if _, ok := infmap.Data[message.Infid].Data[message.Vmid]; !ok {
	 	infmap.infidRUnlock(message.Infid)
	 	infmap.infidLock(message.Infid)
	 	defer infmap.infidUnlock(message.Infid)
	  	infmap.Data[message.Infid].Data[message.Vmid] = &Vmdata{
		  Data: map[string]int{},
	  	}
     }	
  }
  
  infmap.vmidLock(message.Infid, message.Vmid)
  defer infmap.vmidUnlock(message.Infid, message.Vmid)
  infmap.Data[message.Infid].Data[message.Vmid].Data = message.Data
  
  wprint("Inf " + message.Infid + " instance " + message.Vmid, infmap.Data[message.Infid].Data[message.Vmid].Data)
	
}

func (infmap Infmap) vmidRLock(infid string, vmid string) {
	//infmap.RLock()
	infmap.Data[infid].RLock()
	infmap.Data[infid].Data[vmid].RLock()
}

func (infmap Infmap) vmidRUnlock(infid string, vmid string) {
	infmap.Data[infid].Data[vmid].RUnlock()
	infmap.Data[infid].RUnlock()
	//infmap.RUnlock()
}

func (infmap Infmap) vmidLock(infid string, vmid string) {
	//infmap.RLock()
	infmap.Data[infid].RLock()
	infmap.Data[infid].Data[vmid].Lock()
}

func (infmap Infmap) vmidUnlock(infid string, vmid string) {
	infmap.Data[infid].Data[vmid].Unlock()
	infmap.Data[infid].RUnlock()
	//infmap.RUnlock()
}

func (infmap Infmap) infidRLock (infid string) {
	//infmap.RLock()
	infmap.Data[infid].RLock()
}

func (infmap Infmap) infidRUnlock(infid string) {
	infmap.Data[infid].RUnlock()
	//infmap.RUnlock()
}

func (infmap Infmap) infidLock(infid string) {
	//infmap.RLock()
	infmap.Data[infid].Lock()
}

func (infmap Infmap) infidUnlock(infid string) {
	infmap.Data[infid].Unlock()
	//infmap.RUnlock()
}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}