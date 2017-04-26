package main

import (
	//"bufio"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"sync"
    //"io"
    "fmt"
)

type App struct {
	Executable string
	Root  string
	Args []string
	*os.Process
}

type AppError struct {
    // App executable name
    Name string
    // App PID
    Pid int
    // Error message
    Message string
}

func (err AppError) Error() string {
    return fmt.Sprintf("Process %s(PID %d) exited with error, details: %s\n", err.Name, err.Pid, err.Message)
}

type AppOutput struct {
    // process id
    Pid int
    Msg []byte
}

func (o *AppOutput) push() {
    o.msgchan <- o.Msg
}

func FlushOutput(pid int, msg []byte) {
    out := AppOutput{Pid: pid, Msg: msg}
    
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
			if err := app.Run(); err != nil {
                log.Printf("App quited with error: %s", err)
            }
		}(&wg, app)
	}
    
    log.Println("Procman Started!")
    wg.Wait()
}


func (app *App) Run() error {
	if len(app.Executable) == 0 {
		panic("App Executable is empty!")
	}
	
	cmd := exec.Command(app.Executable)
	if len(app.Root) > 0 {
		cmd.Dir = app.Root
	}
	
	////var stderr io.ReadCloser
	//var stdout io.ReadCloser
	//var err error
	////stderr, err = cmd.StderrPipe()
	//stdout, err = cmd.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//scanner := bufio.NewScanner(stdout)
	//go func() {
	//	for scanner.Scan() {
	//		//TODO
	//		fmt.Println(scanner.Text())
	//	}
	//
	//	if err := scanner.Err(); err != nil {
	//		fmt.Fprintln(os.Stderr, "Error encountered while reading input:", err)
	//	}
	//}()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	
	app.Process = cmd.Process
	log.Printf("Process %s started with PID %d\n", app.Executable, app.Process.Pid)
	
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// command fails to run or doesn't complete successfully
			procState := exitError.ProcessState
			if !procState.Exited() {
                app.Process.Kill()
            }
            return AppError{ app.Executable, procState.Pid(), procState.String()}
			
		} else {
			// Other types of error
            return AppError{ app.Executable, app.Process.Pid, err.Error()}
            
        }
	}
    
	return nil
}