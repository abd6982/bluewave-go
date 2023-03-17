package bluewave

import (
	"fmt"
	"strings"
)

func ComparePdfFiles(
	filenames []string,
	methods []string,
	prettyPrint bool,
	verbose bool,
	regenCache bool,
	sidecarOnly bool,
	noImportance bool,
	awsConfig AWSConfig,
) any {
	// t0 := time.Now()

	// define aws client here
	awsClient := AWSClient{}
	// end of aws client

	if verbose {
		fmt.Println("Reading files...")
	}

	// readPdfSecT0 := time.Now()
	fileData := []any{}
	for ind := range filenames {
		fileName := filenames[ind]
		fData := GetFileData(fileName, ind, regenCache, "1.6.4", awsClient)
		fileData = append(fileData, fData)
	}

	// readPdfSec := time.Time.Sub(time.Now(), readPdfSecT0)
	if sidecarOnly {
		return ""
	}

	if len(filenames) < 2 {
		fmt.Println("Must have at least 2 files to compare!")
		return ""
	}

	if len(methods) == 0 {
		if verbose {
			fmt.Println("Methods not specified, using default (all).")
		}
		methods = []string{"pages", "digits", "images", "text"}
	}
	if verbose {
		fmt.Println("Using methods:", strings.Join(methods, ", "))
	}

	return ""
}
