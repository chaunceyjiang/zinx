package utils

import (
	"encoding/json"
	"io/ioutil"
)

type GlobalObj struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Name string `json:"name"`
	Version string `json:"version"`
	MaxPacketSize uint32 `json:"max_packet_size"`
	MaxConn int `json:"max_conn"`
	WorkerPoolSize uint32 `json:"worker_pool_size"`
	MaxWorkerTaskLen uint32 `json:"max_worker_task_len"`
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data,err:=ioutil.ReadFile("conf/zinx.json")
	if err!=nil{
		return
	}
	err = json.Unmarshal(data ,GlobalObject)
	if err!=nil{
		return
	}
}
func init()  {
	GlobalObject = &GlobalObj{
		Host:          "0.0.0.0",
		Port:          7777,
		Name:          "ZinxServer",
		Version:       "0.4",
		MaxPacketSize: 4049,
		MaxConn:       12000,
		WorkerPoolSize: 100,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}