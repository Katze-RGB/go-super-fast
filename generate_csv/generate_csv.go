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
	"time"
	"fmt"
	"strconv"
)

const file_length = 1000
const mean = 50
const sigma = 15

func main() {
	start := time.Now()
	file, err := os.OpenFile("large_data.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("error creating csv", err)
	}

	names := []string{"bob", "robert", "bobby", "joe", "carol", "riley", "ali"}
	defer file.Close()

	//blow in header ahead of anything else
	header := []string{"name", "value"}
	writer := csv.NewWriter(file)
	writer.Write(header)
	defer writer.Flush()

	//loop through and write rows up to file length, with random names and normally distributed values
	for i := 0; i <= file_length; i++ {
		//randint here may not be truely random due to modulo bias, 7 is not evenly divisible by 255.
		randint := rand.Int() % len(names)
		row := []string{names[randint], strconv.FormatFloat(rand.NormFloat64()*sigma+mean, 'G', 2, 64)}
		err := writer.Write(row)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("File generation complete. Generated %v rows in %s \n",file_length, time.Since(start) )
}
