package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	// Parameters
	var filenameFlag = flag.String("filename", "timesheet.txt", "the path to the timesheet file")
	var perHourFlag = flag.Int64("per-hour", 300000, "the amount you of each hour work")
	var currencyFlag = flag.String("currency", "IRR", "the currency unit")
	flag.Parse()

	// Open file
	file, err := os.Open(*filenameFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var sumDiff int64
	var perHour int64 = *perHourFlag

	var eachMinutes int64
	eachMinutes = perHour / 60

	// Read each line
	scanner := bufio.NewScanner(file)
	fmt.Println("+---------------------------------------------------------------+")
	fmt.Println("|------- W O R K   H O U R S   C O M P U T A T I O N -----------|")
	fmt.Println("+---------------------------------------------------------------+")
	for scanner.Scan() {
		line := scanner.Text()

		period := strings.Split(line, " ~ ")

		// Try to parse open date
		openDate, err := time.Parse("2006/01/02 15:04", period[0])
		if err != nil {
			fmt.Printf("Open date '%s' is not valid\n", period[0])
			os.Exit(1)
		}

		// Try to parse close date
		closeDate, err := time.Parse("2006/01/02 15:04", period[1])
		if err != nil {
			fmt.Printf("Close date '%s' is not valid\n", period[1])
			os.Exit(1)
		}

		// Get diff
		diff := (closeDate.Unix() - openDate.Unix()) / 60
		sumDiff += diff

		HH, mm := minutesToHHMM(diff)
		fmt.Printf("| %s ~ %s |  %02d:%02d\n", openDate, closeDate, HH, mm)
	}

	p := message.NewPrinter(language.English)

	fmt.Printf("+---------------------------------------------------------------+\n")
	fmt.Printf(" Full Timesheet:	%d:%02d (%d Minutes)\n", sumDiff/60, sumDiff%60, sumDiff)
	p.Printf(" Salary Per Hour:	%d %s (Each Minutes %d %s)\n", perHour, *currencyFlag, eachMinutes, *currencyFlag)
	p.Printf(" Final Salary:		%d %s\n", sumDiff*eachMinutes, *currencyFlag)
	fmt.Printf("+---------------------------------------------------------------+\n")
	fmt.Printf("| (c) 2020 @mrcyna, All Rights Reserved                         |\n")
	fmt.Printf("| Source available on https://github.com/mrcyna/go-timesheet    |\n")
	fmt.Printf("+---------------------------------------------------------------+\n")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func minutesToHHMM(m int64) (int64, int64) {
	return m / 60, m % 60
}
