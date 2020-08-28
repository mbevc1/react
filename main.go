package main

import (
	"flag"
	"fmt"
	"github.com/mbevc1/react/ssh"
	"github.com/mbevc1/react/yaml"
	"log"
	"os"
	"time"
)

var cmd = flag.String("cmd", "", "Command to run")
var bck = flag.Bool("b", false, "Backup MT device")
var bck_path = flag.String("bck-path", "backups", "Backup path")
var rbt = flag.Bool("r", false, "Reboot MT device")
var upg = flag.Bool("u", false, "Upgrade MT device")
var host = flag.String("h", "", "Host to connect")
var grp = flag.String("g", "", "Filter group")
var mts = flag.String("m", "mts.yaml", "Configuration file")
var hlp = flag.Bool("help", false, "Show command parameters")
var ver = flag.Bool("v", false, "Show version")

func init() {
	// example with short version for long flag
	flag.StringVar(cmd, "c", "", "Command to run")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("%s %v\n", Name, Version)
		if GitCommit != "" {
			fmt.Printf("Commit hash: %v\n", GitCommit)
		}
		os.Exit(0)
	} else if *cmd == "" && !*bck && !*rbt && !*upg || *hlp {
		fmt.Printf("%s %v\n", Name, Version)
		flag.PrintDefaults()
		if *hlp {
			os.Exit(0)
		} else {
			fmt.Println("Missing and action!")
			os.Exit(1)
		}
	}

	//Read Yaml
	config := yaml.Parse(*mts)

	for index, element := range config.Cfgs {
		//Filter by group
		if *grp != "" && *grp != element.Group {
			continue
		}
		//Filter by host
		if *host != "" && *host != element.Host {
			continue
		}
		//Connect to host
		ssh.Conn(&element.Host, &element.User, &element.Pass)
		if *cmd != "" {
			log.Printf("| Device %v: %v\nRunning...\n", index, element.Host)
			ssh.Run(cmd)
			log.Println("| done!")
		}
		if *bck {
			log.Printf("| Device %v: %v\nBacking up...\n", index, element.Host)
			dow := string([]rune(time.Now().Weekday().String())[0:3])
			bck_dir := "backups"
			dir := *bck_path + "/" + element.Group + "/" + element.Name
			merr := os.MkdirAll(dir, os.ModePerm)
			if merr != nil {
				panic(merr)
			}
			//MT backup command
			cmd := fmt.Sprintf("/export compact file=backup ; /tool fetch address=127.0.0.1 mode=ftp user=%s password=%s src-path=backup.rsc dst-path=%s/%s.rsc ; /file remove backup.rsc ; /system backup save name=%s/%s.backup", element.User, element.Pass, bck_dir, dow, bck_dir, dow)
			ssh.Run(&cmd)
			//Download locally
			ssh.Download(bck_dir+"/"+dow+".rsc", dir+"/"+dow+".rsc")
			ssh.Download(bck_dir+"/"+dow+".backup", dir+"/"+dow+".backup")
			log.Println("| done!")
		}
		if *rbt {
			cmd := "/system reboot;"
			log.Printf("| Device %v: %v\nRebooting...\n", index, element.Host)
			ssh.Run(&cmd)
			log.Println("| done!")
		}
		if *upg {
			cmd := "/system package update check; /system package update install; /system routerboard upgrade"
			log.Printf("| Device %v: %v\nUpgrading...\n", index, element.Host)
			ssh.Run(&cmd)
			log.Println("| done!")
		}
		//Disconnect from host
		ssh.Close()
	}
}
