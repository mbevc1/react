package main

import (
	"flag"
	"fmt"
	"github.com/mbevc1/react/ssh"
	"github.com/mbevc1/react/yaml"
	"os"
	"time"
)

var cmd = flag.String("cmd", "", "Command to run")
var bck = flag.Bool("b", false, "Backup MT device")
var rbt = flag.Bool("r", false, "Reboot MT device")
var upg = flag.Bool("u", false, "Upgrade MT device")
var host = flag.String("h", "", "Host to connect")
var grp = flag.String("g", "", "Filter group")
var mts = flag.String("m", "mts.yaml", "Configuration file")
var ver = flag.Bool("v", false, "Show version")

func init() {
	// example with short version for long flag
	flag.StringVar(cmd, "c", "", "Command to run")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("%s v%v\n", Name, Version)
		if GitCommit != "" {
			fmt.Printf("Commit hash: %v\n", GitCommit)
		}
		os.Exit(0)
	} else if *cmd == "" && !*bck && !*rbt && !*upg {
		fmt.Printf("%s v%v\n", Name, Version)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read Yaml
	config := yaml.Parse(*mts)

	for index, element := range config.Cfgs {
		// Filter by group
		if *grp != "" && *grp != element.Group {
			continue
		}
		//Filter by host
		if *host != "" && *host != element.Host {
			continue
		}
		// Connect to host
		ssh.Conn(&element.Host, &element.User, &element.Pass)
		if *cmd != "" {
			fmt.Printf("Entry %v: %v\nRunning...", index, element.Host)
			ssh.Run(cmd)
			fmt.Println("done!")
		}
		if *bck {
			fmt.Printf("Entry %v: %v\nBacking up...", index, element.Host)
			dow := string([]rune(time.Now().Weekday().String())[0:3])
			bck_dir := "backups/"
			dir := bck_dir + element.Group + "/" + element.Name
			merr := os.MkdirAll(dir, os.ModePerm)
			if merr != nil {
				panic(merr)
			}
			ssh.Download(bck_dir+dow+".rsc", dir+"/"+dow+".rsc")
			ssh.Download(bck_dir+dow+".backup", dir+"/"+dow+".backup")
			fmt.Println("done!" + dow)
		}
		if *rbt {
			cmd := "/system reboot;"
			fmt.Printf("Entry %v: %v\nRebooting...", index, element.Host)
			ssh.Run(&cmd)
			fmt.Println("done!")
		}
		if *upg {
			cmd := "/system package update check; /system package update install; /system routerboard upgrade"
			fmt.Printf("Entry %v: %v\nUpgrading...", index, element.Host)
			ssh.Run(&cmd)
			fmt.Println("done!")
		}
		// Disconnect from host
		ssh.Close()
	}
}
