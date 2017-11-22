package main
import (
	"flag"
	"fmt"
	"time"
    "net"
    "os"
    "sync"
    "strconv"
    "math/rand"
    "github.com/satori/go.uuid"
    //"bytes"
    "encoding/gob"
)

type Tmessage struct {
	Infid string
	Vmid string
	Data map[string]int
}

func wprint(param ...interface{}) {
	if !noprint {
		fmt.Println(param)
	}
}

var noprint bool
/*var ninfids int
var nvmids int
var iter int
var err error*/

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
	}
	
	if len(args) > 2 {
		iter, err = strconv.Atoi(args[2])
		if err != nil {
	        iter = 1
	    }
	} else {
		iter = 1
	}
	
	if len(args) > 3 {
		ninfids, err = strconv.Atoi(args[3])
		if err != nil {
	        ninfids = 1
	    }
	} else {
		ninfids = 1
	}
	
	if len(args) > 4 {
		nvmids, err = strconv.Atoi(args[4])
		if err != nil {
	        nvmids = 1
	    }
	} else {
		nvmids = 1
	}
	
	fmt.Println(noprint,iter,ninfids,nvmids)*/
	
	var ninfids int
    var nvmids int
    var iter int
    var CONN_HOST string
    var CONN_TYPE string
    var CONN_PORT string
    //var fnoprint string
    
    flag.IntVar(&ninfids, "ninfids", 1, "Number of infids (default: 1)")
    flag.IntVar(&nvmids, "nvmids", 1, "Number of vmids (default: 1)")
    flag.IntVar(&iter, "iter", 1, "Number of iterations (default: 1)")
    flag.StringVar(&CONN_HOST, "host", "agracia7", "Hostname to connect to (default: agracia7)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Port to connect to (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", false, "Supress output (default: false)")
    
    flag.Parse()
	
	var infidarray []string
	var vmidarray []string
	
	for i := 0; i < ninfids; i++ {
		infid := uuid.NewV1()
		//infidarray[i] = fmt.Sprintf("%s", infid)
		infidarray = append(infidarray, fmt.Sprintf("%s", infid))
	}
	
	for i := 0; i < nvmids; i++ {
		//vmidarray[i] = strconv.Itoa(i)
		vmidarray = append(vmidarray, strconv.Itoa(i))
	}
	
    servAddr := CONN_HOST + ":" + CONN_PORT
    tcpAddr, err := net.ResolveTCPAddr(CONN_TYPE, servAddr)
    if err != nil {
        fmt.Println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }
    
    var wg sync.WaitGroup
    wg.Add(iter)

	start := time.Now()

	for j := 0; j < iter; j++ {
		//wg.Add(1)
		go func(wg *sync.WaitGroup) {
			//fmt.Println("Go in goroutine ...")
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
			        fmt.Println("Dial failed:", err.Error())
			        os.Exit(1)
		    }
	
			infid := rand.Intn(ninfids)
			vmid := rand.Intn(nvmids)
			//cpu := strconv.Itoa(rand.Intn(100))
			//mem := strconv.Itoa(rand.Intn(100))
			//cmessage := new(bytes.Buffer)
		    //gob.NewEncoder(cmessage).Encode(Tmessage{Infid:infidarray[infid],Vmid:vmidarray[vmid],Data:map[string]int{"cpu":rand.Intn(100),"mem":rand.Intn(100)}})
			//wprint("message:",Tmessage{Infid:infidarray[infid],Vmid:vmidarray[vmid],Data:map[string]int{"cpu":rand.Intn(100),"mem":rand.Intn(100)}})
			//strEcho := "{\"infid\":\"" + infidarray[infid] + "\",\"vmid\":\"" + vmidarray[vmid] + "\",\"data\":{ \"cpu\": " + cpu + ",\"mem\": " + mem + " }}\n"
			//wprint(strEcho)
		    //_, err = conn.Write([]byte(strEcho))
		    encoder := gob.NewEncoder(conn)
		    cmessage := &Tmessage{Infid:infidarray[infid],Vmid:vmidarray[vmid],Data:map[string]int{"cpu":rand.Intn(100),"mem":rand.Intn(100)}}
		    encoder.Encode(cmessage)
		    /*_, err = conn.Write(cmessage.Bytes())
			    if err != nil {
		        fmt.Println("Write to server failed:", err.Error())
		        os.Exit(1)
		    }*/
	
		    //wprint("write to server = ", strEcho)
		    wprint("write to server = ", *cmessage)
	
		    reply := make([]byte, 1024)
	
		    _, err = conn.Read(reply)
			    if err != nil {
		        fmt.Println("Write to server failed:", err.Error())
		        os.Exit(1)
		    }
	
		    wprint("reply from server=", string(reply))
	
		    conn.Close()
		    wg.Done()
		}(&wg)
    }
	wg.Wait()
	t := time.Now()
	fmt.Println("Done in",t.Sub(start))
}
