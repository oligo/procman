package main

import (
	"bufio"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"sync"
	"fmt"
    "io"
)

type App struct {
	Executable string
	Root  string
	Args []string
	*os.Process
}

func init() {
	handleSignals()
	LoadConfig()
}

func main() {
	apps := make([]*App, 0)
	if err := viper.UnmarshalKey("apps", &apps); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, app := range apps {
		log.Printf("Starting %s\n", app.Executable)
		wg.Add(1)
		go func(wg *sync.WaitGroup, app *App) {
			defer wg.Done()
			RunApp(app)
		}(&wg, app)
	}
    
    log.Println("Procman Started!")
    wg.Wait()
}

func RunApp(app *App) {
    if len(app.Executable) == 0 {
        panic("App Executable is empty!")
    }
    
	cmd := exec.Command(app.Executable)
    if len(app.Root) > 0 {
        cmd.Dir = app.Root
    }
    
    //var stderr io.ReadCloser
    var stdout io.ReadCloser
    var err error
    //stderr, err = cmd.StderrPipe()
    stdout, err = cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			//TODO
			fmt.Println(scanner.Text())
		}
		
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error encountered while reading input:", err)
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
    
    app.Process = cmd.Process
    log.Printf("Process %s started with PID %d\n", app.Executable, app.Process.Pid)
    
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// command fails to run or doesn't complete successfully
			procState := exitError.ProcessState
			if procState.Exited() {
				log.Printf("Process %s(PID %d) exited with error, details: %s\n", app.Executable, procState.Pid(), procState.String())
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

	
	return 
}



