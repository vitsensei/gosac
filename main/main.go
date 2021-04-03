package main

import (
	"fmt"
	"github.com/vitsensei/gosac/examples"
	"github.com/vitsensei/gosac/ransac"
)

func main() {
	cloud := examples.CreateRandomCloud(100)
	r, _ := ransac.NewRansac(1000, 0.5, 10, 2)
	r.SetModel(&examples.Line{})
	r.SetPopulation(cloud)

	bestLine, loss := r.Run()
	fmt.Println("Best model:", bestLine)
	fmt.Println("with loss: ", loss)

	cloud.ToFile("cloud.csv")
	line, ok := bestLine.(*examples.Line)
	if ok {
		line.ToFile("line.csv")
	}
}
