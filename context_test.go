package belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "c")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextC, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))

	fmt.Println(contextA.Value("b"))
}

func CreateCounter() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)

		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()
	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goruoutine", runtime.NumGoroutine())

	destination := CreateCounter()

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func CreateCounterWithCancel(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)

		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) //simulate slow process
			}
		}
	}()
	return destination
}

func TestContextWithTimeOut(t *testing.T) {
	fmt.Println("Total Goruoutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounterWithCancel(ctx)

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}
	fmt.Println("Total Goruoutine", runtime.NumGoroutine())

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goruoutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5 *time.Second))
	defer cancel()

	destination := CreateCounterWithCancel(ctx)

	fmt.Println("Total Goruoutine", runtime.NumGoroutine())
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}