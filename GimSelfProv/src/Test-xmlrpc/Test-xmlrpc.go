package main

import (
    "fmt"
    "log"
    "bytes"
    //"strings"
    //"io/ioutil"
    "net/http"
    "github.com/divan/gorilla-xmlrpc/xml"
)

type Response struct {
	
}


//func XmlRpcCall(method string, args struct { IMCred struct { Id string`xml:"id"`; Type string`xml:"type"`; Username string`xml:"username"`; Password string`xml:"password"` } }) (reply struct{Message string}, err error) {
func XmlRpcCall(method string, args struct { Credential []interface{} }) (reply struct { Response []map[string][]interface{} `xml:"params"`; }, err error) {
//func XmlRpcCall(method string, args [1]interface{}) (reply struct{Message string}, err error) {	
    
    type IMResponse struct { 
    	IMR []interface{} 
    }
    
    var vIMR []interface{} 
    vIMR = append(vIMR, true)
    vIMR = append(vIMR, []string{ "e4454b32-d9c2-11e7-ba70-42010af08f38", "d547b12a-d9d5-11e7-ba70-42010af08f38" })
    vIMResponse := struct {
    	Message []interface{}
    }{
    	Message: vIMR,
    }
//	vIMResponse := &IMResponse {
//		IMR: vIMR,
//	}
    //reply = vIMResponse
    buf2, _ := xml.EncodeClientRequest(method, &vIMResponse)
    fmt.Println("buf2: " + string(buf2))
    
    buf, _ := xml.EncodeClientRequest(method, &args)
    
    fmt.Println(string(buf))
//    fmt.Printf("%+v\n", args.Credential[0])
//    fmt.Println(args.Credential)

    resp, err := http.Post("http://http://35.205.104.163:8080", "text/xml", bytes.NewBuffer(buf))
    if err != nil {
        return
    }
    //defer resp.Body.Close()
    
//    bodyBytes, _ := ioutil.ReadAll(resp.Body)
//    bodyString := string(bodyBytes)
//    
//    bodyString = strings.TrimSuffix(bodyString, "\n")
//    
//    fmt.Println(bodyString)
    
//    vIMResponse := &IMResponse{
//    	IMR: struct {
//    		
//    	} {
//    		
//    	}
//    }
	
//	var creply struct { Message []interface{} }
//	creply.Message = make([]interface{}, 2)
//	err = xml.DecodeClientResponse(resp.Body, &creply)
//    reply = creply

    //err = xml.DecodeClientResponse(resp.Body, &reply)
    resp.Body.Close()
    return
}

func main() {

	var MyCred []interface{}

	MyCred = append(MyCred, struct {
					Id string`xml:"id"`; 
					Type string`xml:"type"`; 
					Username string`xml:"username"`; 
					Password string`xml:"password"`
				}{
					Id: "INDRAIM",
				    Type: "InfrastructureManager",
				    Username: "admin",
				    Password: "Oem1234-;.",
				},
			)	

    reply, err := XmlRpcCall("GetInfrastructureList", struct {Credential []interface{}}{Credential: MyCred})
    //reply, err := XmlRpcCall("GetInfrastructureList", MyCred)
    if err != nil {
        log.Println("XmlRpcCall error: ",err)
    }

    log.Printf("Response: %v\n", reply.Response)
}

