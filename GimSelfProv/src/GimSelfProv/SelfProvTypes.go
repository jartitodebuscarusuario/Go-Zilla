package main

import (
    "sync"
)

//type Tresponse struct {
//	Result string
//	Infid string
//	Data map[string]int
//}

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
	Conf map[string]interface{}
	Alarm map[string][]bool
}

//Type to store remote infraestructures
type Infmap struct {
	sync.RWMutex
	Data map[string]*Infdata
}

