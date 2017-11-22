package main

import (
	"flag"
    "fmt"
    "net"
    "os"
    "sync"
    "time"
    "strconv"
    "math/rand"
    //"strconv"
    //"encoding/json"
    //"bytes"
    "github.com/satori/go.uuid"
    "encoding/gob"
)

type Tmessage struct {
	Infid string
	Vmid string
	Data map[string]int
}

type Mondata struct {
	Infid string
	Vmid string
	Data map[string]int
}

//Global consts
const nTmessage = 1000000
const ninfids = 100
const nvmids = 5

//Global vars
var hashinf map[string]map[string]map[string]int
var hashconf map[string]map[string]int
var hashalarm map[string]map[string][]bool
var defaulthashconf map[string]int
var evalsp map[string]int

var noprint bool

// Create map of channels, one per machine
var chanhash = make(map[string]map[string]chan map[string]int)

// Lock when we need write to chanhash
var lockchanhash sync.RWMutex

// Var to hold file
//var f *os.File
var err error
//var f, _ = os.Create("D:\\agracia\\Documents\\g\\perl\\GoReport_server.txt")

//var hashinf map[string]map[string]map[string]int
//var infchan chan map[string]map[string]map[string]int

func main() {
    /*args := os.Args
	
	//wprint("args 0",args,len(args))
	
	if len(args) > 1 {
		noprint, err = strconv.ParseBool(args[1])
		if err != nil {
	        noprint = false
	    }
	} else {
		noprint = false
	}*/
    
    var CONN_HOST string
    var CONN_TYPE string
    var CONN_PORT string
    //var fnoprint string
    
    flag.StringVar(&CONN_HOST, "host", "agracia7", "Hostname to bind (default: agracia7)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Hostname to bind (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", false, "Supress output (default: false)")
    
    flag.Parse()
    
    // Open file to write received data
	//f, err = os.Create("D:\\agracia\\Documents\\g\\perl\\GoReport_server.txt")
	
	/*lockchanhash.Lock()
	chanhash["2CD2B0DA-C87B-11E7-82D9-DDF68153C7F6"] = make(map[string] chan map[string]int)
	chanhash["2CD2B0DA-C87B-11E7-82D9-DDF68153C7F6"]["0"] = make(chan map[string]int)
	lockchanhash.Unlock()*/
	//Init chaninf
	//chaninf := make(chan map[string]map[string]map[string]int)
	//Channel to send inf data to hashwriter
	chaninf := make(chan *Tmessage)
	
	//Init hashinf
	hashinf = make(map[string]map[string]map[string]int)
	hashconf = make(map[string]map[string]int)
	hashalarm = make(map[string]map[string][]bool)
	//hashtimer := make(map[string]time.Timer)
	//evalsp, _ := strconv.Atoi("cpu")
	evalsp = map[string]int{
		"cpu": 0,
		"mem": 1,
	}
	defaulthashconf = map[string]int{
		"upcpu": 70,
		"downcpu": 30,
		"upmem": 80,
		"downmem": 20,
		"numalert": 3,
		"numsamples": 5,
		"evalsp": evalsp["cpu"],
		"activesp": 1,
		"tactsp": 10,
		"maxvm": 4,
		"minvm": 1,
	}
	
//	hashinf = map[string]map[string]map[string]int{
//		"infpatata": map[string]map[string]int {
//			"vmpatata": map[string]int {
//				"cpu": 66,
//				"mem": 99,
//			},
//		},
//	}

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
	
	//Only one goroutine to read/write hashinf
	go hashwriter(chaninf)
	
	bcknoprint := noprint;
	noprint = true
	start := time.Now()
	for i := 0; i < nTmessage; i++ {
		chaninf <- prueba[i]
	}
	t := time.Now()
	fmt.Println("Done in",t.Sub(start))
	noprint = bcknoprint;
	
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
        go handleRequest(conn,chaninf)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn, chaninf chan *Tmessage) {
  // Vars to hold json received data and channel read
  // var message data
  //var message interface{}
  //var message Mondata
  //var vmdata map[string]int
  // Make a buffer to hold incoming data.
  dec := gob.NewDecoder(conn)
  message := &Tmessage{}
  dec.Decode(message)
  /* request := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  read_len, err := conn.Read(request)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  // Write received data to disk
  f.Write(request[:read_len])

  var dmessage Tmessage
  gob.NewDecoder(bytes.NewBuffer(request[:read_len])).Decode(&dmessage)
  wprint(dmessage)
  
  message := &dmessage*/

  //json.Unmarshal(request[:read_len], &message)
  wprint(message.Infid,message.Vmid,message.Data,message.Data["cpu"],message.Data["mem"])
  
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received " + message.Infid + " " + message.Vmid))
  
  // Close the connection when you're done with it.
  conn.Close()
  
  /*smessage := map[string]map[string]map[string]int{
	  message.Infid: map[string]map[string]int{
	  	message.Vmid: message.Data,
	  },
  }*/
  
  wprint("Before writing to infchan ...")
  //chaninf <- smessage
  chaninf <- message
  wprint("After writing to infchan ...")
//  lockchanhash.RLock()
//  if infvmidchan, ok := chanhash[message.Infid][message.Vmid]; ok {
//	  wprint("Existe el canal" + message.Infid + " " + message.Vmid)
//	  <- infvmidchan
//	  infvmidchan <- message.Data
//  } else {
//	  wprint("No existe el canal " + message.Infid + " " + message.Vmid)
//	  lockchanhash.RUnlock()
//      lockchanhash.Lock()
//      wprint("Creando el canal " + message.Infid + " " + message.Vmid + " ...")
//	  //Initialize chan with size for non-blocking writes in goroutines
//	  chanhash[message.Infid] = make(map[string] chan map[string]int, 1)
//	  chanhash[message.Infid][message.Vmid] = make(chan map[string]int, 1)
//	  wprint("Creado el canal" + message.Infid + " " + message.Vmid)
//	  lockchanhash.Unlock()
//	  wprint("Desbloqueada la escritura en el canal " + message.Infid + " " + message.Vmid)
//	  lockchanhash.RLock()
//	  wprint("Bloqueo de lectura en el canal " + message.Infid + " " + message.Vmid)
//	  chanhash[message.Infid][message.Vmid] <- message.Data
//	  wprint("Enviados datos",message.Data,"al canal " + message.Infid + " " + message.Vmid)
//  }
//  lockchanhash.RUnlock()
}

func hashwriter (chaninf chan *Tmessage) {
	i := 0;
	start := time.Now()
	for {
		tmessage := <- chaninf
		if _, ok := hashconf[tmessage.Infid]; !ok {
		    hashconf[tmessage.Infid] = defaulthashconf
		}
		
		if _, ok := hashalarm[tmessage.Infid]; !ok {
			arrsamples := make([]bool,hashconf[tmessage.Infid]["numsamples"])
			hashalarm[tmessage.Infid] = map[string][]bool{
				"cpu": arrsamples,
				"mem": arrsamples,
			} 
		}
		
		if i == 1 {
			fmt.Println(hashalarm[tmessage.Infid])
		}

		hashinf[tmessage.Infid] = map[string]map[string]int{
			tmessage.Vmid: tmessage.Data,
		}
		
		wprint(hashinf)
		i++
		if i >= nTmessage - 1 {
			t := time.Now()
			fmt.Println("Done read channel in (",nTmessage," messages)",t.Sub(start))
			start = time.Now()
			i = 0
		}
	}
}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}