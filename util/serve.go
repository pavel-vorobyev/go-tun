package util

import "sync"

func Serve() {
	//termSignal := make(chan os.Signal, 1)
	//signal.Notify(termSignal, os.Interrupt, syscall.SIGTERM)
	//<-termSignal
	//fmt.Println("\nShutting down...")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
