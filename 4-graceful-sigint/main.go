// //////////////////////////////////////////////////////////////////////
// //
// // Given is a mock process which runs indefinitely and blocks the
// // program. Right now the only way to stop the program is to send a
// // SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// // want to try to gracefully stop the process first.
// //
// // Change the program to do the following:
// //   1. On SIGINT try to gracefully stop the process using
// //          `proc.Stop()`
// //   2. If SIGINT is called again, just kill the program (last resort)
// //

// package main

// import (
// 	"os"
// 	"os/signal"
// 	"syscall"
// )

// func main() {
// 	sigs := make(chan os.Signal, 1)

// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	done := make(chan bool)

// 	// Create a process
// 	proc := MockProcess{}

// 	// Run the process (blocking)
// 	// made it non blocking now and once the process is done running we will signal that its done processing using the done channel
// 	go func(){
// 		proc.Run()
// 		done <- true
// 	}()

// 	sigCount := 0

// 	select {
// 	case <- done:
// 		return
// 	case sig := <- sigs:
// 		if(sig == os.Interrupt){
// 			sigCount++
// 			if(sigCount == 1){
// 				// shutdown gracefully
// 				proc.Stop()
// 			} else if(sigCount == 2){
// 				// forced shutdown
// 				os.Exit(1)
// 			}
// 		}
// 	}

// }

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}

	// Channel to listen for OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Channel to signal when Run() completes
	done := make(chan bool)

	// Start the mock process in a goroutine
	go func() {
		proc.Run()
		done <- true
	}()

	// Wait for signals or completion
	sigCount := 0
	for {
		select {
		case <-done:
			// Process finished normally
			return
		case sig := <-sigs:
			if sig == os.Interrupt {
				sigCount++
				if sigCount == 1 {
					// First SIGINT: try graceful shutdown
					// stop is a blocking function
					go func(){
						fmt.Println("Gracefully shutting down")
						proc.Stop()
					}()
				} else if sigCount >= 2 {
					// Second SIGINT: force exit
					fmt.Println("Sayonara")
					os.Exit(1)
				}
			}
		}
	}
}