package main

import "fmt"

func main()  {
	total := 0.0
	for i := 0; i < 10; i++ {
		duration, _ := MosaicConcurrent()
		total += duration.Seconds()
	}
	fmt.Printf("4-partition: %v sec\n", total / 10.0)

	total = 0.0
	for i := 0; i < 10; i++ {
		duration, _ := Mosaic()
		total += duration.Seconds()
	}
	fmt.Printf("no-partition: %v sec\n", total / 10.0)
}