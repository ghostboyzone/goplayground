package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmdName := flag.String("cmd", "", "-cmd \"ls -lh\"")
	loopNum := flag.Int("loop", 1, "-loop 20")
	loopDelay := flag.String("delay", "1s", "-delay 1s")
	flag.Parse()
	log.Println(*cmdName, *loopNum, *loopDelay)
	if *cmdName == "" {
		flag.Usage()
		os.Exit(1)
	}

	pd, _ := time.ParseDuration(*loopDelay)

	for {
		if *loopNum == 0 {
			break
		}
		ret, err := exec.Command("sh", "-c", *cmdName).Output()
		log.Println(ret, err)

		if *loopNum > 0 {
			*loopNum--
		}
		time.Sleep(pd)
	}

}
