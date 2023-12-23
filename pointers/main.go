package main

import "fmt"

type Creature struct {
	Species string
}

func main() {
	var age int = 2
	var getAge *int = &age

	fmt.Println("age=", age)
	fmt.Println("getAge=", getAge)
	fmt.Println("*getAge=", *getAge)
	*getAge = 45
	fmt.Println("*getAge=", *getAge)
	fmt.Println("getAge=", getAge)
	fmt.Println("age=", age)

	//a function that is passing in an argument by value:

	var creature Creature = Creature{Species: "shark"}
	fmt.Printf("1) %+v\n", creature)
	changeCreature(creature)
	fmt.Printf("3) %+v\n", creature)
	changeCreatureWithReference(&creature)
	fmt.Printf("5) %+v\n", creature)

}

func changeCreatureWithReference(creature *Creature) {

	creature.Species = "jellyfish"
	fmt.Printf("4) %+v\n", creature)
}

func changeCreature(creature Creature) {
	creature.Species = "jellyfish"
	fmt.Printf("2) %+v\n", creature)
}
