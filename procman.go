package main

import (
	"bufio"
	"github.com/spf13/viper"
	"log"
	"os/exec"
	"sync"
)

type App struct {
	Name string
	Dir  string
	Args []string
}

func init() {
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
			log.Println("Program output | %s", scanner.Text())
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
