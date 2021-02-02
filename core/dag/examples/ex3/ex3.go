package main

import "infini-framework/core/dag"

func main() {

	d := dag.New()
	d.Spawns(f1, f2, f3).Join().Pipeline(f4)
	d.Run()
}

func f1() {
	println("f1")
}
func f2() {
	println("f2")
}
func f3() {
	println("f3")
}
func f4() {
	println("f4")
}
