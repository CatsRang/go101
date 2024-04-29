package main

import "fmt"

type E interface {
	int64 | string
}

func Backward(s []string) func(func(int, string) bool) {
	return func(yield func(int, string) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
		return
	}
}

func main() {
	// ---- if/else
	{
		if 7%2 == 0 {
			fmt.Println("> 7 is even")
		} else {
			fmt.Println("> 7 is odd")
		}

		// 선언된 변수는 if/else 블록 내에서만 유효하다.
		if num := 8; num%5 == 0 {
			fmt.Println("> 8 is divisible by 5")
		} else if num2 := 9 % 4; num2 == 0 {
			fmt.Printf("> 9%4 == %v\n", num2)
		} else {
			fmt.Println(">", num, num2)
		}
	}

	// ---- for
	{
		// multiple initializers
		for sum, i := 100, 1; i < 10; i++ {
			sum = sum + i
			fmt.Println("> Sum", i, sum)
		}

		for i := range 10 {
			fmt.Println(">", i)
		}

	}

	// ---- for break, continue, goto
	{
	Label3:
		for {
			for {
				fmt.Println("> Infinite loop")
				break Label3
			}
		}
		goto Label4

	Label4:
		fmt.Println("> Label4")
	}

	// ---- for range
	// https://yourbasic.org/golang/for-loop-range-array-slice-map-channel/

	{
		// 01. basic for-each loop(slice or array)
		aA := []string{"a", "b", "c", "d", "e"}
		for i, s := range aA {
			fmt.Println(">", i, s)
		}

		// 02. String iteration: runes or bytes
		for i, ch := range "안녕하세요" {
			fmt.Printf("> %v, %c, %#U\n", i, ch, ch)
		}

		// 03. Map iteration: keys and values
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		for k, v := range m {
			fmt.Printf("> %v, %v\n", k, v)
		}

		// 04. Channel iteration
		ch := make(chan int)
		go func() {
			ch <- 1
			ch <- 2
			ch <- 3
			close(ch)
		}()

		for v := range ch {
			fmt.Println("> ch received", v)
		}
	}

	{
		sl := []string{"hello", "world", "golang"}
		Backward(sl)(func(i int, s string) bool {
			fmt.Printf("%d : %s\n", i, s)
			return true
		})

		// go 1.21 way ???
		/*
			for i, s := range Backward(sl) {
				fmt.Printf("%d : %s\n", i, s)
			}
		*/
	}

}
