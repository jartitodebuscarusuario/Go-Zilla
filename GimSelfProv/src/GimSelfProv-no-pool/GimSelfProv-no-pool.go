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

func main() {

    var CONN_HOST string
    var CONN_TYPE string
    var CONN_PORT string
    
    var nvm int
    //var fnoprint string
    
    flag.StringVar(&CONN_HOST, "host", "agracia7", "Hostname to bind (default: agracia7)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Hostname to bind (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", false, "Supress output (default: false)")
    flag.IntVar(&nvm, "nvm", 1, "Machines per inf (default: 1)")
    
    flag.Parse()
    
    // Open file to write received data
	//f, err = os.Create("D:\\agracia\\Documents\\g\\perl\\GoReport_server.txt")
	
	//Channel to send inf data to hashwriter
	chaninf := make(chan *Tmessage)
	
	//Init hashes
	hashinf = make(map[string]map[string]map[string]int)
	hashconf = make(map[string]map[string]int)
	hashalarm = make(map[string]map[string][]bool)
	//hashtimer := make(map[string]time.Timer)
	
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
		"nvm": nvm,
	}

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
	mstart := time.Now()
	for i := 0; i < nTmessage; i++ {
		chaninf <- prueba[i]
	}
	mt := time.Now()
	fmt.Println("Done in",mt.Sub(mstart))
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

  // Make a buffer to hold incoming data.
  dec := gob.NewDecoder(conn)
  message := &Tmessage{}
  dec.Decode(message)

  // Write received data to disk
  //f.Write(request[:read_len])

  wprint(message.Infid,message.Vmid,message.Data,message.Data["cpu"],message.Data["mem"])
  
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received " + message.Infid + " " + message.Vmid))
  
  // Close the connection when you're done with it.
  conn.Close()
  
  wprint("Before writing to infchan ...")

  chaninf <- message
  wprint("After writing to infchan ...")

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
		
		var totalcpu, totalmem, upcpu, downcpu, upmem, downmem int
		
		if len(hashinf[tmessage.Infid]) == hashconf[tmessage.Infid]["nvm"] {
			
			totalcpu, totalmem = evalselfprov(tmessage.Infid, hashinf[tmessage.Infid])
			hashinf[tmessage.Infid] = make(map[string]map[string]int)
			if totalcpu > hashconf[tmessage.Infid]["upcpu"] {
				upcpu = 1
			}
			if totalcpu < hashconf[tmessage.Infid]["downcpu"] {
				downcpu = 1
			}
			if totalcpu > hashconf[tmessage.Infid]["upmem"] {
				upmem = 1 
			}
			if totalcpu < hashconf[tmessage.Infid]["downmem"] {
				downmem = 1 
			}
		}
		
		//Check why this is not printing
		wprint(totalcpu, totalmem, upcpu, downcpu, upmem, downmem)

		i++
		if i >= nTmessage - 1 {
			t := time.Now()
			fmt.Println("Done read channel in (",nTmessage," messages)",t.Sub(start))
			start = time.Now()
			i = 0
		}
	}
}

func evalselfprov (infid string, infidhash map[string]map[string]int) (totalcpu int, totalmem int) {
	
//	var totalcpu int
//	var totalmem int
	
	vtotal := make(map[string][]int)
	
	for _, value := range infidhash {
	    vtotal["cpu"] = append(vtotal["cpu"], value["cpu"])
	    vtotal["mem"] = append(vtotal["mem"], value["mem"])
	}
	
	for _, value := range vtotal["cpu"] {
		totalcpu += value
	}
	
	for _, value := range vtotal["mem"] {
		totalmem += value
	}
	
	totalcpu = totalcpu/len(vtotal["cpu"])
	totalmem = totalmem/len(vtotal["mem"])
	
//	fmt.Println(totalcpu/len(vtotal["cpu"]))
//	fmt.Println(totalmem/len(vtotal["mem"]))
	return totalcpu, totalmem

}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}