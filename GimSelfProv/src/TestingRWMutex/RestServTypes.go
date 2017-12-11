package main

//Type to receive http json commands
type Tcommand struct {
	Command string
	Value string
	Data map[string]interface{}
}

type Tresponse struct {
	Result string
	Infid string
	Data map[string]interface{}
}

type Configuration struct {
	Data map[string]interface{}
}