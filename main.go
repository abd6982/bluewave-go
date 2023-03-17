package main

import "github.com/bballboy8/bluewave-go/bluewave"

func main() {
	bluewave.ComparePdfFiles(
		[]string{"sample_file_1.pdf", "sample_file_2.pdf"},
		[]string{},
		true,
		true,
		false,
		false,
		false,
		bluewave.AWSConfig{},
	)
}
