package main

import (
	"fmt"
)

// Entity ... An entity is an entity
type Entity struct {
	firstName string
	weight    float32
}

func updateByVal(val Entity) {
	val.weight *= 1.1
}

func updateByRef(val *Entity) {
	val.weight *= 1.1
}

func main() {
	anEntity := Entity{firstName: "Issa", weight: 3.330}
	fmt.Printf("Hello, my name is %s and I have %.2f \n", anEntity.firstName, anEntity.weight)
	updateByVal(anEntity)
	fmt.Printf("Post update by val: %f\n", anEntity.weight)
	updateByRef(&anEntity)
	fmt.Printf("Post update by ref: %f\n", anEntity.weight)

	fmt.Println(anEntity)

	fmt.Println(new(Entity))
	ptrEntity := new(Entity)

	array1 := [5]Entity{anEntity, *ptrEntity, Entity{"Ayee", 45.056}}
	fmt.Println("array1 ", array1)
	array2 := new([5]Entity)
	fmt.Println("array2 ", array2)

	slice1 := array1[1:3]
	fmt.Println("slice1 ", slice1)
}
