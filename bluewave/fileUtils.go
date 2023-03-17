package bluewave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/liyue201/gostl/ds/set"
	"github.com/liyue201/gostl/utils/comparator"
)

type CachedData struct {
	Version string
	Data    any
}

func PageSkipConditions(pageText string) *set.Set[bool] {
	conds := set.New[bool](comparator.BoolComparator, set.WithGoroutineSafe())

	conds.Insert(strings.Contains(pageText, "FORM FDA "))
	conds.Insert(strings.Contains(pageText, "Form FDA "))
	conds.Insert(strings.Contains(pageText, "PAPERWORK REDUCTION ACT"))
	conds.Insert(strings.Contains(pageText, "PAYMENT IDENTIFICATION NUMBER"))
	conds.Insert(strings.Contains(pageText, "For more assistance with Adobe Reader"))
	conds.Insert(strings.Contains(pageText, "latest version of Adobe Reader"))
	conds.Insert(strings.Contains(pageText, ".................................................................."))
	conds.Insert(strings.Contains(pageText, "Safety Data Sheet"))
	conds.Insert(strings.Contains(pageText, "SAFETY DATA SHEET"))
	conds.Insert(strings.Contains(pageText, "Contains Nonbinding Recommendations"))

	return conds
}

func BlockSkipConditions(blockText string) *set.Set[bool] {
	conds := set.New[bool](comparator.BoolComparator, set.WithGoroutineSafe())

	conds.Insert(strings.Contains(blockText, "510(k)"))
	conds.Insert(strings.Contains(blockText, "New Hampshire Avenue"))
	conds.Insert(strings.Contains(blockText, "ISO "))
	conds.Insert(strings.Contains(blockText, "IEC "))
	conds.Insert(strings.Contains(blockText, ".............."))
	conds.Insert(strings.Contains(blockText, "Tel.:"))
	conds.Insert(strings.Contains(blockText, "TEL:"))
	conds.Insert(strings.Contains(blockText, "FAX:"))
	conds.Insert(strings.Contains(blockText, "Fax:"))
	conds.Insert(strings.Contains(blockText, "+86"))
	conds.Insert(strings.Contains(blockText, "86-519"))

	return conds
}

func GetPageBlockAndHashes(fileName string, pageNum int) (any, any) {
	var imageHashes []any
	var textBlocks []any

	// open file and process
	doc, err := fitz.New(fileName)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	// cumLenText := 0
	// cumLenDigits := 0
	// page := doc.Page(pageNum)

	pageSkipFlag := false
	pageText, err := doc.Text(pageNum)
	if err != nil {
		fmt.Println("Error in reading page")
		return imageHashes, textBlocks
	}

	pageSkipConditons := PageSkipConditions(pageText)
	for iter := pageSkipConditons.Begin(); iter.IsValid(); iter.Next() {
		if iter.Value() {
			pageSkipFlag = true
			break
		}
	}

	if !pageSkipFlag {
		// blockNum := 0
		pageImage, err := doc.Image(pageNum)
		if err != nil {
			fmt.Println("Error getting image")
		}
		fmt.Println(pageImage)
	}

	return imageHashes, textBlocks
}

func IsCompatible(vCurrent any, vCache any) bool {
	return true
}

func GetPageCount(fileName string) int {
	doc, err := fitz.New(fileName)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	nPages := doc.NumPage()
	return nPages
}

func ReadBlocksAndHashes(fileName string) (any, any, any, error) {
	fileExt := filepath.Ext(fileName)
	fmt.Println(fileExt)
	if fileExt != ".pdf" {
		return nil, nil, nil, fmt.Errorf("fitz cannot read non-PDF file %s", fileName)
	}

	var textBlocks []any
	var imageHashes []any

	nPages := GetPageCount(fileName)
	// runtime.GOMAXPROCS()
	for i := 0; i < nPages; i++ {
		imgHashes, txtBlocks := GetPageBlockAndHashes(fileName, i)
		imageHashes = append(imageHashes, imgHashes)
		textBlocks = append(textBlocks, txtBlocks)
	}

	return textBlocks, imageHashes, nPages, nil
}

func GetFileData(
	fileName string,
	index int,
	regenCache bool,
	version string,
	awsClient AWSClient,
) any {
	var blocks any
	var imageHashes any
	var nPages any

	// Check if cached exists next to PDF
	// if exists, check if version is compatible
	cachedFilename := fileName + ".jsoncached"
	if (awsClient != AWSClient{}) {
		downloadCacheFromS3(cachedFilename)
	}
	_, err := os.Stat(cachedFilename)
	if err == nil && !regenCache {
		content, err := ioutil.ReadFile(cachedFilename)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
		var cached CachedData
		err = json.Unmarshal(content, &cached)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}
		if IsCompatible(version, cached.Version) {
			// blocks, image_hashes, n_pages = cached["data"]
		}
	}

	if blocks == nil && imageHashes == nil && nPages == nil {
		blocks, imageHashes, nPages, err := ReadBlocksAndHashes(fileName)
		if err != nil {
			fmt.Println("Error while reading blocks and hashes")
			fmt.Println(err)
			return ""
		}
		fmt.Println(blocks)
		fmt.Println(imageHashes)
		fmt.Println(nPages)
	}

	return ""
}
