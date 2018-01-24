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
	Addr           string   `json:"addr"`
	Paths          []Path   `json:"paths"`
	SendChannels   int      `json:"send_channels"`
	BatchSendFiles int      `json:"batch_send_files"`
	ExcludeFiles   []string `json:"exclude_files"`
}

type Path struct {
	Local  string   `json:"local"`
	Remote []string `json:"remote"`
}

type Local2RemoteMap map[string][]string

func NewConf(filePath string) ConfigConf {
	log.Println("Read Conf: ", filePath)
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)
	defer file.Close()

	var myConf ConfigConf
	decoder.Decode(&myConf)
	log.Println("Conf: ", myConf)
	return myConf
}

func (c *ConfigConf) GetPaths() Local2RemoteMap {
	local2RemoteMap := make(Local2RemoteMap)
	for _, v := range c.Client.Paths {
		local2RemoteMap[v.Local] = v.Remote
	}
	return local2RemoteMap
}
