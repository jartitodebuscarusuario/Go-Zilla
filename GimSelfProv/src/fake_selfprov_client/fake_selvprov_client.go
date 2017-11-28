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

func main() {
	
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
    flag.StringVar(&CONN_HOST, "host", "::1", "Hostname to connect to (default: agracia7)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Port to connect to (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp6", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", false, "Supress output (default: false)")
    
    flag.Parse()
	
	var infidarray []string
	var vmidarray []string
	
	for i := 0; i < ninfids; i++ {
		infid := uuid.NewV1()
		infidarray = append(infidarray, fmt.Sprintf("%s", infid))
	}
	
	for i := 0; i < nvmids; i++ {
		vmidarray = append(vmidarray, strconv.Itoa(i))
	}
	
    servAddr := "[" + CONN_HOST + "]:" + CONN_PORT
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
			conn, err := net.DialTCP(CONN_TYPE, nil, tcpAddr)
			if err != nil {
			        fmt.Println("Dial failed:", err.Error())
			        os.Exit(1)
		    }
	
			infid := rand.Intn(ninfids)
			vmid := rand.Intn(nvmids)
		    encoder := gob.NewEncoder(conn)
		    cmessage := &Tmessage{
		    	Infid:infidarray[infid],
		    	Vmid:vmidarray[vmid],
		    	Data:map[string]int{
		    		"cpu":rand.Intn(100),
		    		"mem":rand.Intn(100),
		        },
		    }
		    encoder.Encode(cmessage)

		    wprint("write to server = ", *cmessage)
	
		    reply := make([]byte, 1024)
	
		    _, err = conn.Read(reply)
			    if err != nil {
		        fmt.Println("Write to server failed:", err.Error())
		        os.Exit(1)
		    }
			conn.Close()
		    wprint("reply from server=", string(reply))
		    //time.Sleep(100 * time.Millisecond)
		    wg.Done()
		}(&wg)
    }
	wg.Wait()
	t := time.Now()
	fmt.Println("Done in",t.Sub(start))
}
