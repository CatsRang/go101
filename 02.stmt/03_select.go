package main

import (
	"fmt"
	"time"
)

func main() {
	/*
		The select statement in Go is used to handle multiple channel operations.
		It's similar to a switch statement, but it operates on channels,
		allowing you to wait on multiple channel operations and handle whichever one is ready first.
		If multiple channels are ready, one is chosen at random.
	*/

	{
		ch1 := make(chan string)
		ch2 := make(chan string)

		go func() {
			time.Sleep(2 * time.Second)
			ch1 <- "Message from ch1"
		}()

		go func() {
			time.Sleep(1 * time.Second)
			ch2 <- "Message from ch2"
		}()

		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		}
	}

	/*
		select with a default case
	*/
	{
		ch := make(chan string)

		go func() {
			time.Sleep(2 * time.Second)
			ch <- "Message"
		}()

		select {
		case msg := <-ch:
			fmt.Println("Received:", msg)
		default:
			fmt.Println("No message received")
		}
	}

	/*
		NonBlocking Select
	*/
	{
		ch1 := make(chan string, 1)
		ch2 := make(chan string, 1)

		//ch1 <- "Hello from ch1"

		select {
		case msg1 := <-ch1:
			fmt.Println("> Ch1: ", msg1)
		case msg2 := <-ch2:
			fmt.Println("> Ch1: ", msg2)
		default:
			fmt.Println("No message received")
		}

	}

}
