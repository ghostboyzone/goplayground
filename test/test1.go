package main

import (
	// "encoding/binary"
	// "encoding/gob"
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	"log"
	// "os"
	"os/exec"
	"strings"
	// "reflect"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	// "bytes"
	// proto123 "github.com/emicklei/proto"
)

func main() {
	tt, err := ParseFile("/home/map/odp_commodity/php/phplib/protobuf/proto/google/protobuf/descriptor.proto", "/home/map/odp_commodity/php/phplib/protobuf/proto/")
	if err != nil {
		panic(err)
	}
	log.Println(tt)
}

type errCmd struct {
	output []byte
	err    error
}

func (this *errCmd) Error() string {
	return this.err.Error() + ":" + string(this.output)
}

func ParseFile(filename string, paths ...string) (*descriptor.FileDescriptorSet, error) {
	return parseFile(filename, false, true, paths...)
}

func parseFile(filename string, includeSourceInfo bool, includeImports bool, paths ...string) (*descriptor.FileDescriptorSet, error) {
	args := []string{"--proto_path=" + strings.Join(paths, ":")}
	if includeSourceInfo {
		args = append(args, "--include_source_info")
	}
	if includeImports {
		args = append(args, "--include_imports")
	}
	args = append(args, "--descriptor_set_out=/dev/stdout")
	args = append(args, filename)
	cmd := exec.Command("protoc", args...)
	cmd.Env = []string{}
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, &errCmd{data, err}
	}
	fileDesc := &descriptor.FileDescriptorSet{}
	if err := proto.Unmarshal(data, fileDesc); err != nil {
		return nil, err
	}
	return fileDesc, nil
}
