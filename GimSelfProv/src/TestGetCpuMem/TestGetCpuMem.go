package main

import (
    "fmt"
    "time"
    //"github.com/shirou/gopsutil/mem"
    //"github.com/shirou/gopsutil/load"
    "github.com/shirou/gopsutil/cpu"
)

func main() {
    //v, _ := mem.VirtualMemory()

    // almost every return value is a struct
    //fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

    // convert to JSON. String() is also implemented
    //fmt.Println(v)
    
    //v, _ := load.Avg()
    
    // almost every return value is a struct
    //fmt.Printf("Load1: %v, Load5:%v, Load15:%v\n", v.Load1, v.Load5, v.Load15)

    // convert to JSON. String() is also implemented
    //fmt.Println(v)
    
    v, _ := cpu.Percent(15 * time.Second, false)
    
    fmt.Printf("Percent(15, false): %v\n", v)
    
}
