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
	// myCh         chan int
	confPath     *string
	confHost     *string
	confPort     *string
	confUsername *string
	confPassword *string
	confThread   *int
	confDb       *string
	sqlSucc      int64
	sqlErr       int64
)

func main() {

	sqlSucc = 0
	sqlErr = 0
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

	myCh := make(chan int, *confThread)

	var allPaths []string

	filepath.Walk(*confPath, func(path string, info os.FileInfo, err error) error {

		if len(path)-strings.LastIndex(path, ".sql") == 4 {
			// log.Println(path)
			allPaths = append(allPaths, path)
		}

		return nil
	})

	totalLoop := 0

	reg, _ := regexp.Compile(`/\*.*\*/`)

	for _, onePath := range allPaths {
		log.Println(onePath)

		f, _ := os.Open(onePath)
		scanner := bufio.NewScanner(f)

		totalStr := ""

		// var sqlArr []string

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

			// strings.LastIndex(lineStr, ";")

			totalStr += lineStr

			if len(lineStr)-strings.LastIndex(lineStr, ";") == 1 {
				totalStr = strings.Replace(totalStr, "INSERT INTO", "INSERT IGNORE INTO", 1)
				myCh <- 1
				go inDb(totalStr, myCh)
				totalStr = ""
			}

			/*
				for _, i := range strings.Split(lineStr, "") {
					totalStr += i
					if i == ";" {
						if strings.Replace(totalStr, " ", "", -1) != ";" {

							totalStr = strings.Replace(totalStr, "INSERT INTO", "INSERT IGNORE INTO", 1)
							// sqlArr = append(sqlArr, totalStr)
							myCh <- 1
							go inDb(totalStr, myCh)

						}
						totalStr = ""
					}
				}
			*/

			// if len(sqlArr) >= 100 {
			//     myCh <- 1
			//     inDb(sqlArr, myCh)
			//     sqlArr :=
			// }

			totalLoop++
			if totalLoop >= 1000 {
				log.Println(sqlSucc, sqlErr)
				totalLoop = 0
			}
		}

		if len(totalStr) != 0 {
			// sqlArr = append(sqlArr, totalStr)
			myCh <- 1
			inDb(totalStr, myCh)
		}

		// log.Println(sqlArr)

		// myCh <- 1
		// inDb(sqlArr, myCh)
	}

}

func inDb(sql string, ch chan int) {
	conn, err := client.Connect(*confHost+":"+*confPort, *confUsername, *confPassword, *confDb)
	if err != nil {
		log.Panic(err)
	}

	_, err = conn.Execute(sql)
	if err != nil {
		sqlErr++
		log.Println("Err:", err, " Sql:", sql)
	} else {
		sqlSucc++
	}

	<-ch
}
