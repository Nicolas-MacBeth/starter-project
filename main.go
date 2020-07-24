package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

const jmxPort int = 8686

var jmxFlags = []string{"-Dcom.sun.management.jmxremote", fmt.Sprintf("-Dcom.sun.management.jmxremote.port=%v", jmxPort), "Dcom.sun.management.jmxremote.local.only=true", "-Dcom.sun.management.jmxremote.authenticate=false", "-Dcom.sun.management.jmxremote.ssl=false"}

func main() {
	log.Println("Welcome to the suite of ingredient finding services")

	// waitGroup for parent main to stay alive while subprocesses are still alive
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(4)

	// Each child is a *process* but is still spawned from different goroutines, for better performance and waitGroup management
	go startChildProcess("foodfinder", waitGroup)
	go startChildProcess("foodsupplier", waitGroup)
	go startChildProcess("foodvendor", waitGroup)
	go startChildProcessJava("foodcounter", "ServerCounter", waitGroup)

	waitGroup.Wait()
}

func startChildProcessJava(directory string, file string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Compile the class
	build := exec.Command("javac", fmt.Sprintf("%v/%v.java", directory, file))
	attachChildOutputToParent(build)
	errBuild := build.Run()
	if errBuild != nil {
		log.Printf("%v %v", directory, errBuild)
		return // Exit from this specific goroutine, since build failed,
	}

	// Run the server
	childProcess := exec.Command("java", jmxFlags[0], jmxFlags[1], jmxFlags[2], jmxFlags[3], jmxFlags[4], fmt.Sprintf("%v.%v", directory, file))
	attachChildOutputToParent(childProcess)

	errProcess := childProcess.Run() // This line will block execution and deferred Done() will not run until the server crashes, keeping parent main alive
	if errProcess != nil {
		log.Printf("%v %v", directory, errProcess)
	}
}

// Start every child process individually
func startChildProcess(directory string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Build the executable
	build := exec.Command("go", "build", "server.go")
	build.Dir = directory
	attachChildOutputToParent(build)
	errBuild := build.Run()
	if errBuild != nil {
		log.Printf("%v %v", directory, errBuild)
		return // Exit from this specific goroutine, since build failed,
	}

	// Run the executable
	childProcess := exec.Command("./server")
	childProcess.Dir = directory
	attachChildOutputToParent(childProcess)

	errProcess := childProcess.Run() // This line will block execution and deferred Done() will not run until the server crashes, keeping parent main alive
	if errProcess != nil {
		log.Printf("%v %v", directory, errProcess)
	}
}

// Swap the child processes' Stdout and Stderr to something more useful (parent processe's Stdout for now)
func attachChildOutputToParent(childProcess *exec.Cmd) {
	// Note: Stdout could also be routed to something else, say a log file, in the future if logs become too cumbersome
	// For now, each server has its logging prefix, which is organized enough
	childProcess.Stdout = os.Stdout
	childProcess.Stderr = os.Stderr
}
