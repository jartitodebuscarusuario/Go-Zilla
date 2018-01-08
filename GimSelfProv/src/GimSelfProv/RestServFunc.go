package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "net/http"
    "html/template"
    "encoding/json"
    "github.com/gorilla/websocket"   
)

//var readLogTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html>
//<html>
//<head>
//<meta charset="utf-8">
//<script>  
//window.addEventListener("load", function(evt) {
//    var output = document.getElementById("output");
//    var input = document.getElementById("input");
//    var ws;
//    var print = function(message) {
//        var d = document.createElement("div");
//        d.innerHTML = message;
//        output.appendChild(d);
//    };
//    document.getElementById("open").onclick = function(evt) {
//        if (ws) {
//            return false;
//        }
//        ws = new WebSocket("{{.}}");
//        ws.onopen = function(evt) {
//            print("OPEN");
//        }
//        ws.onclose = function(evt) {
//            print("CLOSE");
//            ws = null;
//        }
//        ws.onmessage = function(evt) {
//            print("RESPONSE: " + evt.data);
//        }
//        ws.onerror = function(evt) {
//            print("ERROR: " + evt.data);
//        }
//        return false;
//    };
//    document.getElementById("send").onclick = function(evt) {
//        if (!ws) {
//            return false;
//        }
//        print("SEND: " + input.value);
//        ws.send(input.value);
//        return false;
//    };
//    document.getElementById("close").onclick = function(evt) {
//        if (!ws) {
//            return false;
//        }
//        ws.close();
//        return false;
//    };
//});
//</script>
//</head>
//<body>
//<table>
//<tr><td valign="top" width="50%">
//<p>Click "Open" to create a connection to the server, 
//"Send" to send a message to the server and "Close" to close the connection. 
//You can change the message and send multiple times.
//<p>
//<form>
//<button id="open">Open</button>
//<button id="close">Close</button>
//<p><input id="input" type="text" value="Hello world!">
//<button id="send">Send</button>
//</form>
//</td><td valign="top" width="50%">
//<div id="output"></div>
//</td></tr></table>
//</body>
//</html>`))

var readLogTemplate = template.Must(template.New("").Parse(readTemplate("index.html")))

var upgrader = websocket.Upgrader{}

func readTemplate (templ string) string {
	dir, derr := filepath.Abs(filepath.Dir(os.Args[0]))
    if derr != nil {
	  wprint(derr)
    }
	dat, _ := ioutil.ReadFile(dir + "/" + templ)
	return string(dat[:])
}

func loadconf(conf map[string]interface{}, dconf map[string]interface{}) (result map[string]interface{}, intresult map[string]interface{}) {
	result = make(map[string]interface{})
	intresult = make(map[string]interface{})
	if _, ok := conf["defaultconf"]; !ok {
		for key, _ := range dconf {
			switch param := key; param {
			  	case "evalsp":
				  	if _, ok := conf[key]; !ok {
				  		intresult[key] = string(dconf[key].(string))
				  	} else {
				  		intresult[key] = string(conf[key].(string))
				  	}
				  	//intresult[key] = result[key]
				  	/*var paramsp int
					for key, val := range evalsp {
					  	if key == string(result["evalsp"]) {
					  	paramsp = val
					  	}
					}
					result[key] = paramsp */
				  	result[key] = intresult[key]
			  	default:
				  	if _, ok := conf[key]; !ok {
				  		result[key] = int(dconf[key].(float64))
				  	} else {
				  		result[key] = int(conf[key].(float64))
				  	}
				  	intresult[key] = result[key]
		  	}
		}
	} else {
		for key, _ := range dconf {
			switch param := key; param {
			  	case "evalsp":
				  	/* var paramsp int
					for key, val := range evalsp {
					  	if key == string(dconf["evalsp"].(string)) {
					  	paramsp = val
					  	}
					} */
					result[key] = string(dconf["evalsp"].(string))
					intresult[key] = dconf[key]
			  	default:
				  	result[key] = int(dconf[key].(float64))
				  	intresult[key] = result[key]	  	
		  	}
		}
	}
	return result, intresult
}

func HttpService(w http.ResponseWriter, r *http.Request) {
    //w.Header().Set("Content-Type", "application/json")
    fmt.Println(r.RequestURI)
    //decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()
    //var t Tcommand

    switch myURL := r.RequestURI; myURL {
    	case "/reginfid":
	    	regInfid(w, r)
	    case "/readlog":
		    readLog(w, r)	
	    case "/wslog":
		    wsLog(w, r, &LogChan)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side   
    }
}

func readLog (w http.ResponseWriter, r *http.Request) {
	readLogTemplate.Execute(w, "ws://"+r.Host+"/wslog")
}

func wsLog (w http.ResponseWriter, r *http.Request, channel *chan string) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wprint("upgrade:", err)
		fmt.Println("upgrade:", err)
		return
	}
	defer c.Close()
	mt, message, err := c.ReadMessage()
	if err != nil {
		wprint("read:", err)
		return
	} else {
		wprint("Websocket client said", string(message))
		fmt.Println("Websocket client said", string(message))
	}
	
	for {
		smessage := <-(*channel)
		err = c.WriteMessage(mt, []byte(smessage))
		fmt.Println("Send message: " + string(smessage))
		if err != nil {
			wprint("write:", err)
			fmt.Println("write:", err)
			break
		}
	}
}

func regInfid (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var t Tcommand
	err := decoder.Decode(&t)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
		return
    }
    wprint("Command/Value: ", t.Command, t.Value)
    switch method := r.Method; method {	
	    //if t.Command == "regInfid" {
	    case "PUT":
	        switch command := t.Command; command {
	        	case "regInfid":	
			    	infmap.RLock()
			    	if _, ok := infmap.Data[t.Value]; !ok {
				    	infmap.RUnlock()
				    	newconf, sendconf := loadconf(t.Data, defconf)
				    	var emptyAlarm []bool
						for i := 0; i < newconf["numsamples"].(int); i++ {
						  emptyAlarm = append(emptyAlarm, false)	
						}
						
				    	infmap.Lock()
						infmap.Data[t.Value] = &Infdata{ 
						  Data: map[string]*Vmdata{},
						  Conf: newconf,
						  Alarm: map[string][]bool {
						  	"upcpu": emptyAlarm,
						  	"downcpu": emptyAlarm,
						  	"upmem": emptyAlarm,
						  	"downmem": emptyAlarm,
						  },
						}
				    	
						infmap.Unlock()
						infmap.RLock()
						rmessage := &Tresponse{
					      Result: "OK",
					      Infid:t.Value,
					      Data:sendconf,
						}
						w.WriteHeader(http.StatusCreated)
						json.NewEncoder(w).Encode(rmessage)
						wprint("Registered inf: " + t.Value)
			    	} else {
			    		w.WriteHeader(http.StatusConflict)
				    	fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
			    	}
			    	infmap.RUnlock()
			    default:
				    w.WriteHeader(http.StatusConflict)
			    	fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
	        }		
	    default:
	    	w.WriteHeader(http.StatusMethodNotAllowed)
	    	fmt.Fprintf(w, "{ \"result\" : \"NOK\" }") // send data to client side
	}
}