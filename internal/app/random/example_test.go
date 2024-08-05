package random

import "fmt"

func Example() {
	newURL := NewRandomURL()

	fmt.Println(newURL)
}
