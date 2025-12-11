/// beginner

//////////////////////////////////////////////////////////////////////
// //
// // Your video processing service has a freemium model. Everyone has 10
// // sec of free processing time on your service. After that, the
// // service will kill your process, unless you are a paid premium user.
// //
// // Beginner Level: 10s max per request
// // Advanced Level: 10s max per user (accumulated)
// //

// package main

// import "time"

// // User defines the UserModel. Use this to check whether a User is a
// // Premium user or not
// type User struct {
// 	ID        int
// 	IsPremium bool
// 	TimeUsed  int64 // in seconds
// }

// // HandleRequest runs the processes requested by users. Returns false
// // if process had to be killed
// func HandleRequest(process func(), u *User) bool {
// 	doneChan := make(chan bool)

// 	go func() {
// 		process()
// 		doneChan <- true
// 	}()

// 	timeout := time.After(10 * time.Second)

// 	select {
// 		case d := <- doneChan:
// 			return d
// 		case <- timeout:
// 			return false
// 	}
// }

// func main() {
// 	RunMockServer()
// }

// original
// package main

// // User defines the UserModel. Use this to check whether a User is a
// // Premium user or not
// type User struct {
// 	ID        int
// 	IsPremium bool
// 	TimeUsed  int64 // in seconds
// }

// // HandleRequest runs the processes requested by users. Returns false
// // if process had to be killed
// func HandleRequest(process func(), u *User) bool {
// 	process()
// 	return true
// }

// func main() {
// 	RunMockServer()
// }

// advanced

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu        sync.Mutex // shared resource across threads
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
    u.mu.Lock()

    timeLimitSeconds := 10 - u.TimeUsed
    if timeLimitSeconds <= 0 {
        u.mu.Unlock()
        return false
    }

    timeLimit := time.Duration(timeLimitSeconds) * time.Second

    startTime := time.Now()
    doneChan := make(chan bool)

    go func() {
        process()
        doneChan <- true
    }()

    select {
    case <-time.After(timeLimit):
        u.mu.Unlock()
        return false

    case <-doneChan:
        duration := time.Since(startTime)

        u.TimeUsed = u.TimeUsed + int64(duration.Seconds())

        u.mu.Unlock()
        return true
    }
}

func main() {
	RunMockServer()
}
