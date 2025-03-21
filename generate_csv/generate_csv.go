package main

//generates normally distributed test data for your csv crunching pleasure
//takes aroun 17 mins to generate a 1b row (13ish gb) csv
//would not recommend trying to open said csv in your text editor/ide
//prrrobably could multithread it, but I'm pretty confident io speeds are the primary bottleneck
//adjust file lenght, and if you change the file name from large_data.csv please for the love of god add it to your gitignore
import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
)

const file_length = 1000000000
const mean = 50
const sigma = 10

func main() {
	file, err := os.OpenFile("large_data.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("error creating csv", err)
	}
	defer file.Close()
	header := []string{"value", "name"}
	writer := csv.NewWriter(file)
	writer.Write(header)
	defer writer.Flush()
	for i := 0; i <= file_length; i++ {
		row := []string{"test_name", strconv.FormatFloat(rand.NormFloat64()*sigma+mean, 'G', 2, 64)}
		err := writer.Write(row)
		if err != nil {
			panic(err)
		}
	}
	// blow in header
	writer.Write(header)
	defer writer.Flush()
}
