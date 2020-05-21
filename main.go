package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	log.Println("Welcome to the suite of ingredient finding services")

	// waitGroup for parent main to stay alive while subprocesses are still alive
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(3)

	// Each child is a *process* but is still spawned from different goroutines, for better performance and waitGroup management
	go startChildProcess("foodfinder", waitGroup)
	go startChildProcess("foodsupplier", waitGroup)
	go startChildProcess("foodvendor", waitGroup)

	waitGroup.Wait()
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
