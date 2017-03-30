package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigConf struct {
	Server ServerInfo
	Client ClientInfo
}

type ServerInfo struct {
	Addr            string `json:"addr"`
	ReadBufferSize  int    `json:"read_buffer_size"`
	WriteBufferSize int    `json:"write_buffer_size"`
}

type ClientInfo struct {
	Addr         string   `json:"addr"`
	Paths        []Path   `json:"paths"`
	SendChannels int      `json:"send_channels"`
	ExcludeFiles []string `json:"exclude_files"`
}

type Path struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type HashMap map[string]string

func NewConf(filePath string) ConfigConf {
	log.Println("Read Conf: ", filePath)
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	var myConf ConfigConf
	decoder.Decode(&myConf)
	log.Println("Conf: ", myConf)
	return myConf
}

func (c *ConfigConf) GetPaths() HashMap {
	rMap := make(HashMap)
	for _, v := range c.Client.Paths {
		rMap[v.Local] = v.Remote
	}
	return rMap
}
