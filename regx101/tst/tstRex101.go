package main

import (
	"fmt"
	"regexp"
)

func main()  {
	r, err := regexp.Compile(`Hello`)
	if err != nil {
		fmt.Printf("There is a problem with your regexp.\n")
		return
	}
	// Will print 'Match'
	if r.MatchString("Hello Regular Expression.") == true {
		fmt.Printf("Match ")
	} else {
		fmt.Printf("No match ")
	}
}
