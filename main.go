package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

const filePath = "generate_csv/large_file.csv"
const chunkSize = 1000000
const numWorkers = 8

type chunkData struct {
	chunkSize    int
	chunkAverage int
	chunkMin     int
	chunkMax     int
}

func main() {
	start := time.Now()
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	//lesson learned here: buffer size in a channel is important
	//manage chunk size and number of workers carefully as to avoid overflowing channel bufs and causing a deadlock.
	//total workers*chunksize should do it? except if chunk size is like <2 digits for a big base file?
	jobs := make(chan [][]string, numWorkers)
	results := make(chan chunkData, numWorkers*chunkSize)
	var wg sync.WaitGroup
	header, err := reader.Read()
	if err != nil {
		fmt.Println("error reading header", err)
		return
	}
	fmt.Println("header:", header)

	// Start worker goroutines, increments waitgroup delta
	for range numWorkers {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Read CSV in chunks and send to job channel
	headerSkipped := false
	for {
		chunk, err := readChunk(reader, chunkSize)
		if err == io.EOF {
			close(jobs)
			break
		} else if err != nil {
			panic(err)
		}
		if !headerSkipped {
			chunk = chunk[1:]
			headerSkipped = true
		}

		jobs <- chunk
	}
	wg.Wait()
	close(results)

	// Aggregate results
	totalProcessed := 0
	for count := range results {
		totalProcessed += count.chunkSize
		fmt.Println(count)
	}

	fmt.Println("Total processed rows:", totalProcessed)
	fmt.Println("time spent:", time.Since(start))
}

// does exactly what it says on the tin, there's probably a
func readChunk(reader *csv.Reader, chunkSize int) ([][]string, error) {
	var chunk [][]string
	for range chunkSize {
		record, err := reader.Read()
		if err == io.EOF {
			if len(chunk) > 0 {
				return chunk, nil
			}
			return nil, err
		} else if err != nil {
			fmt.Println(err)
			continue
		}
		chunk = append(chunk, record)
	}
	return chunk, nil
}

// spawns workers that take chunks from job channel, enqueues and processes them via results
func worker(jobs <-chan [][]string, results chan<- chunkData, wg *sync.WaitGroup) {
	defer wg.Done()
	for chunk := range jobs {
		count := processChunk(chunk)
		results <- count
	}
}

func processChunk(chunk [][]string) chunkData {
	sum := 0
	base_year, err := strconv.Atoi(chunk[0][1])
	if err != nil {
		fmt.Println("Error converting base year to int for current chunk")
	}
	min_year := base_year
	max_year := base_year
	fmt.Println("chunksize: ", len(chunk))
	for row := range chunk {
		year, err := strconv.Atoi(chunk[row][1])
		if err != nil {
			fmt.Println("Error converting year to int on line:", row)
			continue
		}

		if year < min_year {
			min_year = year
		}

		if year > max_year {
			max_year = year
		}
		sum += year
	}
	size := len(chunk)
	out := chunkData{chunkSize: size,
		chunkAverage: sum / size,
		chunkMin:     min_year,
		chunkMax:     max_year,
	}
	return out
}
