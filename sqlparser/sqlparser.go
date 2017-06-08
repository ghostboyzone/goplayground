package main

import (
	"bufio"
	// "container/list"
	"flag"
	"github.com/siddontang/go-mysql/client"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	myCh         chan int
	confPath     *string
	confHost     *string
	confPort     *string
	confUsername *string
	confPassword *string
	confThread   *int
	confDb       *string
)

func main() {
	confPath = flag.String("path", "", "-path ./")
	confHost = flag.String("h", "127.0.0.1", "-h 127.0.0.1")
	confPort = flag.String("P", "3306", "-P 3306")
	confUsername = flag.String("u", "root", "-u root")
	confPassword = flag.String("p", "", "-p")
	confThread = flag.Int("t", 20, "-t 20")
	confDb = flag.String("db", "", "-db test")

	flag.Parse()
	if *confPath == "" || *confDb == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	myCh = make(chan int, *confThread)

	var allPaths []string

	filepath.Walk(*confPath, func(path string, info os.FileInfo, err error) error {

		if len(path)-strings.LastIndex(path, ".sql") == 4 {
			// log.Println(path)
			allPaths = append(allPaths, path)
		}

		return nil
	})

	for _, onePath := range allPaths {
		log.Println(onePath)

		f, _ := os.Open(onePath)
		scanner := bufio.NewScanner(f)

		reg, _ := regexp.Compile(`/\*.*\*/`)

		totalStr := ""

		var sqlArr []string

		for scanner.Scan() {
			if len(totalStr) != 0 {
				totalStr += " "
			}

			lineStr := scanner.Text()

			lineStr = reg.ReplaceAllString(lineStr, "")

			// skip blank line
			if len(lineStr) == 0 {
				continue
			}

			// log.Println(lineStr)

			for _, i := range strings.Split(lineStr, "") {
				totalStr += i
				if i == ";" {
					if strings.Replace(totalStr, " ", "", -1) != ";" {
						totalStr = strings.Replace(totalStr, "INSERT INTO", "INSERT IGNORE INTO", 1)
						sqlArr = append(sqlArr, totalStr)
					}
					totalStr = ""
				}
			}
		}

		if len(totalStr) != 0 {
			sqlArr = append(sqlArr, totalStr)
		}

		log.Println(sqlArr)

		myCh <- 1
		inDb(sqlArr, myCh)
	}

}

func inDb(sqls []string, ch chan int) {
	conn, err := client.Connect(*confHost+":"+*confPort, *confUsername, *confPassword, *confDb)
	if err != nil {
		log.Panic(err)
	}

	for _, sql := range sqls {
		_, err := conn.Execute(sql)
		if err != nil {
			log.Println("Err:", err, " Sql:", sql)
		}
	}

	<-ch
}
