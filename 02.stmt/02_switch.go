package main

import "fmt"

func main() {

	// ---- 01. Simple switch
	{
		day := 3

		switch day {
		case 1:
			fmt.Println("> Monday")
		case 2:
			fmt.Println("> Tuesday")
		default:
			fmt.Printf("> Unknown day: %v\n", day)
		}
	}

	// ---- 02. Tagless
	{
		// age := 20
		switch age := 20; {
		case age < 18:
			fmt.Println("> Minor")
		default:
			fmt.Println("> Adult")
		}
	}

	// ---- 03.  case list
	{
		c := 'a'
		fmt.Printf("> %T, %c\n", c, c) // int32(rune), 97

		switch c {
		case 'a', 'b':
			fmt.Println("> alphabet:", c)
		default:
			fmt.Println("> unknown:", c)
		}
	}

	// ---- 04. Switch type
	{
		var i interface{} = 'r'
		switch i.(type) {
		case float64:
			fmt.Println("> float64", c)
		case int32:
			fmt.Printf("> rune, %U, %U, %U, %U\n", 'a', 'z', 'A', 'Z')

			ri := i.(rune)
			if (ri >= 'a' && ri <= 'z') || (ri >= 'A' && ri <= 'Z') {
				fmt.Printf("> rune, alphabet: %c, %#U\n", ri, ri)
				break
			}

			fmt.Printf("> rune, non alphabet: %c\n", i)
		default:
			fmt.Printf("> rune, unknown %T, %c\n", c, c)
		}
	}

	// fallthrough

	// break
}
