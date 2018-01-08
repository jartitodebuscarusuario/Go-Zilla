package main
import (
	"flag"
	"fmt"
	"time"
    "net"
    "os"
    "log"
    "path/filepath"
    "io/ioutil"
    "sync"
    "sync/atomic"
    "strconv"
    "math/rand"
    "github.com/satori/go.uuid"
    //"bytes"
    //"encoding/gob"
    "encoding/json"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/cpu"
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
var rerrorcount int64
var serrorcount int64
var confInfid string
var confVmid string

func main() {
	
	var ninfids int
    var nvmids int
    var iter int
    var nmessages int
    var encoder string
    var tpause int
    
    var CONN_HOST string
    var CONN_TYPE string
    var CONN_PORT string
    var cpuPeriod interface{}
    //var fnoprint string
    
    flag.IntVar(&ninfids, "ninfids", 1, "Number of infids (default: 1)")
    flag.IntVar(&nvmids, "nvmids", 1, "Number of vmids (default: 1)")
    flag.IntVar(&iter, "iter", 1, "Number of iterations (default: 1)")
    flag.IntVar(&nmessages, "nmessages", 1, "Number of iterations (default: 1)")
    flag.IntVar(&tpause, "tpause", 0, "Number of seconds to pause between messages (default: 0)")
    flag.StringVar(&CONN_HOST, "host", "::1", "Hostname to connect to (default: ::1)")
    flag.StringVar(&CONN_PORT, "port", "8888", "Port to connect to (default: 8888)")
    flag.StringVar(&CONN_TYPE, "type", "tcp6", "Connection type (default: tcp)")
    flag.BoolVar(&noprint, "noprint", false, "Supress output (default: false)")
    flag.StringVar(&encoder, "encoder", "json", "Message encoding type (default: json)")
    
    flag.Parse()
    
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }
    
    defconf := readConf(dir + "/clientConfig.json")
    
    //Conf file takes priority
    if rhost, ok := defconf["host"]; ok {
    	CONN_HOST = string(rhost.(string))
    }
    
    if rport, ok := defconf["port"]; ok {
    	CONN_PORT = string(rport.(string))
    }
    
    if valCpu, ok := defconf["cpuPeriod"]; ok {
    	cpuPeriod = int(valCpu.(float64))
    } else {
    	cpuPeriod = 15
    }
    
    if valInfid, ok := defconf["infid"]; ok {
    	confInfid = valInfid.(string)
    }
    
    if valVmid, ok := defconf["vmid"]; ok {
    	confVmid = valVmid.(string)
    }
    
    rerrorcount = 0
    serrorcount = 0
	
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
    fmt.Println("HOST, PORT:", CONN_HOST, CONN_PORT)
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
			        serrorcount = atomic.AddInt64(&serrorcount, 1)
			        wg.Done()
			        return
			        //os.Exit(1)
		    }
	
			for k := 0; k < nmessages; k++ {
	
			infid := rand.Intn(ninfids)
			vmid := rand.Intn(nvmids)
			
		    //encoder := gob.NewEncoder(conn)
		    encoder := json.NewEncoder(conn)
		    v, _ := mem.VirtualMemory()
		    c, _ := cpu.Percent(time.Duration(cpuPeriod.(int)) * time.Second, false)
		    var strInfid string
		    var strVmid string
		    
		    if confInfid != "" {
		    	strInfid = confInfid
		    } else {
		    	strInfid = infidarray[infid]
		    }
		    
		    if confVmid != "" {
		    	strVmid = confVmid
		    } else {
		    	strVmid = vmidarray[vmid]
		    }
		    
		    cmessage := &Tmessage{
		    	//Infid:infidarray[infid],
		    	//Vmid:vmidarray[vmid],
		    	Infid:strInfid,
		    	Vmid:strVmid,
		    	Data:map[string]int{
		    		//"cpu":rand.Intn(100),
		    		//"mem":rand.Intn(100),
		    		"cpu":int(c[0]),
		    		"mem":int(v.UsedPercent),
		        },
		    }
		    encoder.Encode(cmessage)
		    wprint("write to server = ", *cmessage)
	
		    reply := make([]byte, 1024)
	
		    _, err = conn.Read(reply)
		    if err != nil {
	        fmt.Println("Write to server failed (error reading reply):", err.Error())
	        //os.Exit(1)
	        rerrorcount = atomic.AddInt64(&rerrorcount, 1)
	        wg.Done()
	        return
		    }
			//conn.Close()
		    wprint("reply from server=", string(reply))
		    time.Sleep(time.Duration(tpause) * time.Second)
		    }
	    conn.Close()
	    wg.Done()
		}(&wg)
	    time.Sleep(time.Duration(tpause) * time.Second)
	}
	wg.Wait()
	t := time.Now()
	fmt.Println("Done in",t.Sub(start)," Errors read: ",rerrorcount,"Errors connect: ",serrorcount)
	//time.Sleep(10 * time.Second)
}

//Read config file in json format
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
