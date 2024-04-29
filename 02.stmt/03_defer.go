package main

func main() {

	// ---- if/else
	{
		if 7%2 == 0 {
			println("> 7 is even")
		} else {
			println("> 7 is odd")
		}

		if num := 8 % 4; num == 0 {
			println("> 8 is divisible by 4")
		} else if num%3 == 0 {
			println("> 8 is divisible by 3")
		}
	}

}
