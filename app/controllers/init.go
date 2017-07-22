package controllers

import (
	"github.com/revel/revel"
)

// var (
// 	intSig  chan os.Signal = make(chan os.Signal, 1)
// 	termSig chan os.Signal = make(chan os.Signal, 1)
// 	stopSig chan os.Signal = make(chan os.Signal, 1)
// 	quitSig chan os.Signal = make(chan os.Signal, 1)
// )

// func sigHandler(s os.Signal) {
// 	cleanUp()
// 	log.Println("Quitting service...")
// 	os.Exit(-1)
// }

// func installSignalHandler(cb func(os.Signal)) {
// 	signal.Notify(termSig, syscall.SIGINT)
// 	signal.Notify(termSig, syscall.SIGTERM)
// 	signal.Notify(stopSig, syscall.SIGTSTP)
// 	signal.Notify(quitSig, syscall.SIGQUIT)
// 	select {
// 	case <-intSig:
// 		cb(syscall.SIGINT)
// 	case <-termSig:
// 		cb(syscall.SIGTERM)
// 	case <-stopSig:
// 		cb(syscall.SIGTSTP)
// 	case <-quitSig:
// 		cb(syscall.SIGQUIT)
// 	}
// }

// func cleanUp() {
// 	redis.Pool.Close()
// }

func init() {

	// go installSignalHandler(sigHandler)
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
