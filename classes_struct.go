package main

import "fmt"

// Person is a person :)
type Person struct {
	sex string
}

func CreatePersonOfSex(sex string) *Person {
	return &Person{sex}
}

// Recognize
func (p *Person) Recognize() string {
	switch p.sex {
	case "Male":
		return "Male here"
	case "Female":
		return "Female here"
	default:
		return "I would rather not say"
	}
}

type Reproducable interface {
	Recognize() string
}

func RecognizeMe(r Reproducable) {
	fmt.Printf("Hey I am %s\n", r.Recognize())
}

func main() {
	aPerson := CreatePersonOfSex("")
	RecognizeMe(aPerson)
}
