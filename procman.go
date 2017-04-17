package main

import (
	"bufio"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"sync"
)

type App struct {
	Name string
	Dir  string
	Args []string
	*os.Process
}

var managedApps = make([]*App, 5)
var started = make(chan *App, 1)

func init() {
	handleSignals()
	LoadConfig()
}

func main() {
	apps := make([]App, 0)
	if err := viper.UnmarshalKey("apps", &apps); err != nil {
		panic(err)
	}
	log.Println(apps)

	var wg sync.WaitGroup
	for _, app := range apps {
		log.Printf("Starting %s\n", app.Name)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			RunApp(&app)
		}(&wg)
	}

	for i := 0; i< len(apps); i++ {
		select {
			case a := <-started:
				managedApps = append(managedApps, a)
		}
	}


	log.Println(managedApps)

	wg.Wait()

	log.Println("Procman exited!")
}

func RunApp(app *App) {
	cmdDir := app.Dir
	cmdName := cmdDir + app.Name

	cmd := exec.Command(cmdName)
	cmd.Dir = cmdDir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			//TODO
			log.Println(scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}


	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// command fails to run or doesn't complete successfully
			procState := exitError.ProcessState
			if procState.Exited() {
				log.Printf("Process %d exited with error, details: %s\n", procState.Pid(), procState.String())
				//os.Exit(1)
			} else {
				// Not exited?
				log.Println("Process %d failed but NOT exited!", procState.Pid())
				os.Exit(1)

			}

		} else {
			// Other types of error
			log.Fatal(err)
		}
		return 
	}

	app.Process = cmd.Process
	started <- app
	return 
}



