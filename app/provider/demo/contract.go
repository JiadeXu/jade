package demo

const DemoKey = "demo.md"

type IService interface {
	GetAllStudent() []Student
}

type Student struct {
	ID   int
	Name string
}
