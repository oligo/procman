package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type App struct {
	Executable string
	Root       string
	Args       []string
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

//type AppOutput struct {
//    // process id
//    Pid int
//    Msg []byte
//}
//
//func (o *AppOutput) push() {
//    o.msgchan <- o.Msg
//}
//
//func FlushOutput(pid int, msg []byte) {
//    out := AppOutput{Pid: pid, Msg: msg}
//
//}

func RunAll(apps []*App) {
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

func KillAll(apps []*App) {
	for _, app := range apps {
		if app.Process != nil {
			if exitStatus := app.Quit(); exitStatus != nil {
				log.Printf("Killing process %s with exit message: %s", app.Executable, exitStatus.Error())
			}
		}
	}
	return
}

func (app *App) Run() error {
	if len(app.Executable) == 0 {
		return errors.New("Empty app config, omitting.")
	}
	cmd := exec.Command(app.Executable, app.Args...)
	if len(app.Root) > 0 {
		cmd.Dir = app.Root
	}

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
			return AppError{app.Executable, procState.Pid(), procState.String()}

		} else {
			// Other types of error
			return AppError{app.Executable, app.Process.Pid, err.Error()}

		}
	}

	return nil
}

func (app *App) Quit() error {
	if app.Process == nil || app.Process.Pid == 0 {
		return errors.New("os: process not initialized")
	}

	if app.Process.Pid == -1 {
		return errors.New("os: process already released")
	}

	if exitStatus := app.Process.Signal(syscall.SIGTERM); exitStatus != nil {
		fmt.Printf("pid: %d, error: %s\n", app.Process.Pid, exitStatus)
		app.Process.Kill()
		return nil
	}
	fmt.Printf("Pid: %d exited successfully!\n", app.Process.Pid)

	return nil
}
