package main

import (
	"fmt"
)

func getTypeName(x interface{}) string {
	switch x := x.(type) {
	case int:
		return "integer"

	case string:
		return "string"
	case float64:
		return "float"
	default:
		return fmt.Sprintf("Unknown type %T", x)
	}
}

func typeAssertionWithoutSwitch(y interface{}) int {

	if value, ok := y.(int); ok {
		fmt.Printf("y is an integer: %d\n", value)
		return value
	} else {
		fmt.Printf("y is not an integer\n")
		return 0
	}

}
func main() {

	fmt.Println(getTypeName("hello"))
	fmt.Println(getTypeName(45))
	fmt.Println(getTypeName(true))
	fmt.Println(typeAssertionWithoutSwitch("42"))
}
