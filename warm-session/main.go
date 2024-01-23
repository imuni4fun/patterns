package main

import (
	"context"
	"fmt"
	"time"
	"warm-session/synchronizer"
)

func main() {
	_ = initSynchronizer()

	go func() {
		pullImage("nginx:latest", "Always", time.Duration(2*time.Minute))
	}()

	killAfter(5)
}

var syncro synchronizer.ISychronizer

func initSynchronizer() synchronizer.ISychronizer {
	ctx, _ := context.WithCancel(context.Background())

	sessionChan := make(chan synchronizer.Session, 1000) // must give lots of buffer to ensure no deadlock or dropped requests
	completedChan := make(chan string, 10)               // must give some buffer to ensure no deadlock

	asyncTimeout := time.Duration(10 * time.Minute)

	syncro = synchronizer.GetSynchronizer(ctx, asyncTimeout, sessionChan, completedChan)
	synchronizer.RunPullerLoop(ctx, sessionChan, completedChan)

	return syncro
}

func pullImage(image, pullPolicy string, callerTimeout time.Duration) error {
	ses, err := syncro.StartPull(image, pullPolicy)
	if err != nil {
		err = fmt.Errorf("main: pull image: error returned from StartPull for image %s %s", image, err)
		fmt.Println(err.Error())
		return err
	}
	err = syncro.WaitForPull(ses, callerTimeout)
	if err != nil {
		err = fmt.Errorf("main: pull image: error returned from WaitForPull for image %s %s", image, err.Error())
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("main: pull image: success for image %s\n", image)
	return nil
}

func killAfter(sec int) {
	killAfter := time.Duration(sec) * time.Second
	fmt.Printf("sleeping now for %v...\n", killAfter)
	time.Sleep(killAfter)
	fmt.Printf("shutting down...\n")
}
