package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	pm "github.com/uscott/go-pm"
)

type CX struct {
	errorC chan error
	doneC  chan struct{}
	wait   time.Duration
}

func NewCX() *CX {
	return &CX{errorC: make(chan error), doneC: make(chan struct{})}
}

func (x *CX) Done() chan struct{} {
	return x.doneC
}

func (x *CX) Error() chan error {
	return x.errorC
}

func (x *CX) WaitTime() time.Duration {
	return x.wait
}

func (x *CX) OnFork() {
	fmt.Println("Forking")
}

func (x *CX) OnQuit() {
	fmt.Println("Exiting")
}

func (x *CX) Run() {
	rand.Seed(time.Now().UnixNano())
	for {
		n := rand.Int31n(20)
		fmt.Println(n)
		switch n {
		case 0:
			x.errorC <- fmt.Errorf("n == 0")
		case 19:
			x.doneC <- struct{}{}
		default:
		}
		time.Sleep(time.Second / 10)
	}
}

func main() {

	x := NewCX()
	x.wait = 15 * time.Second

	pm, err := pm.NewProcessManager(":8008")
	if err != nil {
		fmt.Printf("Exiting: %v\n", err.Error())
		os.Exit(1)
	}

	pm.Sub = x
	if err = pm.Run(); err != nil {
		fmt.Printf("Error, exiting: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Exiting")
}
