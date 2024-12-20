package grizzly

import (
	"math"
	"runtime"
	"strconv"
	"sync"
)

func arrayFloatBase(initValue float64, data []float64, operation func(info float64, result float64) float64) chan float64 {
	length := len(data)
	if length == 0 {
		// Handle empty data case by returning a closed channel immediately
		emptyChan := make(chan float64)
		close(emptyChan)
		return emptyChan
	}
	numGoroutines := runtime.NumCPU()
	if numGoroutines > length {
		numGoroutines = length // Avoid creating more goroutines than necessary
	}
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	resultChan := make(chan float64, numGoroutines)

	// Function to calculate the sum of a chunk
	worker := func(start, end int) {
		defer wg.Done()
		result := initValue
		// Always starts from second value to calculate Mean Correctly
		for i := start; i < end; i++ {
			if !math.IsNaN(data[i]) {
				result = operation(data[i], result)
			}
		}
		resultChan <- result
	}

	// Launch goroutines to process chunks
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func arrayStringBase(initValue float64, data []string, operation func(info string, result float64) float64) chan float64 {
	length := len(data)
	if length == 0 {
		// Handle empty data case by returning a closed channel immediately
		emptyChan := make(chan float64)
		close(emptyChan)
		return emptyChan
	}
	numGoroutines := runtime.NumCPU()
	if numGoroutines > length {
		numGoroutines = length // Avoid creating more goroutines than necessary
	}
	chunkSize := (length + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup
	resultChan := make(chan float64, numGoroutines)

	// Function to calculate the sum of a chunk
	worker := func(start, end int) {
		defer wg.Done()
		result := initValue
		// Always starts from second value to calculate Mean Correctly
		for i := start; i < end; i++ {
			result = operation(data[i], result)
		}
		resultChan <- result
	}

	// Launch goroutines to process chunks
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go worker(start, end)
	}

	// Wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func arrayStringCountWord(data []string, word string) float64 {
	CPUNumbers := runtime.NumCPU()

	chunkSize := len(data) / CPUNumbers
	if len(data)%CPUNumbers != 0 {
		chunkSize++ // Handle cases where data is not evenly divisible
	}

	results := make(chan float64, CPUNumbers)
	var wg sync.WaitGroup

	// Worker function to count occurrences in a chunk
	worker := func(start, end int) {
		defer wg.Done()
		var localCount float64
		for i := start; i < end && i < len(data); i++ {
			if data[i] == word {
				localCount++
			}
		}
		results <- localCount
	}

	// Divide the work into chunks and spawn goroutines
	for i := 0; i < CPUNumbers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		wg.Add(1)
		go worker(start, end)
	}

	// Close the results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results
	var total float64
	for count := range results {
		total += count
	}

	return total
}

func arrayFloatCountValue(data []float64, value float64) float64 {
	CPUNumbers := runtime.NumCPU()

	chunkSize := len(data) / CPUNumbers
	if len(data)%CPUNumbers != 0 {
		chunkSize++ // Handle cases where data is not evenly divisible
	}

	results := make(chan float64, CPUNumbers)
	var wg sync.WaitGroup

	// Worker function to count occurrences in a chunk
	worker := func(start, end int) {
		defer wg.Done()
		var localCount float64
		for i := start; i < end && i < len(data); i++ {
			if data[i] == value {
				localCount++
			}
		}
		results <- localCount
	}

	// Divide the work into chunks and spawn goroutines
	for i := 0; i < CPUNumbers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		wg.Add(1)
		go worker(start, end)
	}

	// Close the results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results
	var total float64
	for count := range results {
		total += count
	}

	return total
}

func arrayFloatCountNaNValue(data []float64) float64 {
	CPUNumbers := runtime.NumCPU()

	chunkSize := len(data) / CPUNumbers
	if len(data)%CPUNumbers != 0 {
		chunkSize++ // Handle cases where data is not evenly divisible
	}

	results := make(chan float64, CPUNumbers)
	var wg sync.WaitGroup

	// Worker function to count occurrences in a chunk
	worker := func(start, end int) {
		defer wg.Done()
		var localCount float64
		for i := start; i < end && i < len(data); i++ {
			if math.IsNaN(data[i]) {
				localCount++
			}
		}
		results <- localCount
	}

	// Divide the work into chunks and spawn goroutines
	for i := 0; i < CPUNumbers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		wg.Add(1)
		go worker(start, end)
	}

	// Close the results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results
	var total float64
	for count := range results {
		total += count
	}

	return total
}

func arrayCountFloatDuplicates(elements []float64) map[float64]int {
	numCPU := runtime.NumCPU()
	length := len(elements)
	if length == 0 {
		return nil // Handle empty input
	}

	// Calculate chunk size for splitting work
	chunkSize := (length + numCPU - 1) / numCPU
	resultChan := make(chan map[float64]int, numCPU)
	var wg sync.WaitGroup

	// Worker function to count duplicates in a subset of elements
	countDuplicates := func(subset []float64) {
		defer wg.Done()
		localCounts := make(map[float64]int)
		for _, element := range subset {
			localCounts[element]++
		}
		resultChan <- localCounts // Send local result to the channel
	}

	// Start goroutines to process chunks of the array
	for i := 0; i < numCPU && i*chunkSize < length; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > length {
			end = length
		}
		wg.Add(1)
		go countDuplicates(elements[start:end])
	}

	// Close the channel once all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// mergeFloat results from all goroutines
	combinedCounts := make(map[float64]int)
	for localCounts := range resultChan {
		for key, count := range localCounts {
			combinedCounts[key] += count
		}
	}

	// Filter results to retain only elements with counts greater than 1
	finalResult := make(map[float64]int)
	for key, count := range combinedCounts {
		if count > 1 {
			finalResult[key] = count
		}
	}

	return finalResult
}

func arrayContainsInteger(arr []int, target int) bool {
	for _, value := range arr {
		if value == target {
			return true // Element found
		}
	}
	return false // Element not found
}

func arrayContainsString(arr []string, target string) bool {
	for _, value := range arr {
		if value == target {
			return true // Element found
		}
	}
	return false // Element not found
}

// arrayGetNonFloatValues identifies non-convertible float values using goroutines
func arrayGetNonFloatValues(input []string) []string {
	numGoroutines := runtime.NumCPU() // Number of goroutines to use
	chunkSize := (len(input) + numGoroutines - 1) / numGoroutines

	var wg sync.WaitGroup
	var mu sync.Mutex // Protects shared access to `nonConvertible`
	var nonConvertible []string

	// Launch goroutines to process chunks
	for g := 0; g < numGoroutines; g++ {
		start := g * chunkSize
		end := start + chunkSize
		if start >= len(input) {
			break
		}
		if end > len(input) {
			end = len(input)
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localNonConvertible := []string{} // Local slice to collect results

			for _, str := range input[start:end] {
				if _, err := strconv.ParseFloat(str, 64); err != nil {
					localNonConvertible = append(localNonConvertible, str)
				}
			}

			// Append results to the shared slice
			mu.Lock()
			nonConvertible = append(nonConvertible, localNonConvertible...)
			mu.Unlock()
		}(start, end)
	}

	wg.Wait()
	return nonConvertible
}

func arrayResizeString(input []string, targetLength int, defaultValue string) []string {
	for len(input) < targetLength {
		input = append(input, defaultValue)
	}
	return input
}

func arrayResizeFloat(input []float64, targetLength int, defaultValue float64) []float64 {
	for len(input) < targetLength {
		input = append(input, defaultValue)
	}
	return input
}

func arrayUniqueValuesFloat(arr []float64) []float64 {
	if len(arr) == 0 {
		return []float64{}
	}

	numGoroutines := runtime.NumCPU()
	if numGoroutines > len(arr) {
		numGoroutines = len(arr) // Avoid spawning more goroutines than necessary
	}

	chunkSize := (len(arr) + numGoroutines - 1) / numGoroutines

	// Channel to collect results from goroutines
	results := make(chan map[float64]struct{}, numGoroutines)
	var wg sync.WaitGroup

	// Divide work among goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunk := arr[start:end]

		wg.Add(1)
		go func(data []float64) {
			defer wg.Done()
			// Calculate unique values for the chunk
			chunkUnique := make(map[float64]struct{}, len(data))
			for _, val := range data {
				chunkUnique[val] = struct{}{}
			}
			results <- chunkUnique
		}(chunk)
	}

	// Close results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// mergeFloat results from all chunks
	finalUnique := make(map[float64]struct{})
	for chunkResult := range results {
		for key := range chunkResult {
			finalUnique[key] = struct{}{}
		}
	}

	// Convert the unique values map to a slice
	uniqueValues := make([]float64, 0, len(finalUnique))
	for key := range finalUnique {
		uniqueValues = append(uniqueValues, key)
	}

	return uniqueValues
}

func arrayUniqueValuesString(arr []string) []string {
	if len(arr) == 0 {
		return []string{}
	}

	numGoroutines := runtime.NumCPU()
	if numGoroutines > len(arr) {
		numGoroutines = len(arr) // Avoid spawning more goroutines than necessary
	}

	chunkSize := (len(arr) + numGoroutines - 1) / numGoroutines

	// Channel to collect results from goroutines
	results := make(chan map[string]struct{}, numGoroutines)
	var wg sync.WaitGroup

	// Divide work among goroutines
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunk := arr[start:end]

		wg.Add(1)
		go func(data []string) {
			defer wg.Done()
			// Calculate unique values for the chunk
			chunkUnique := make(map[string]struct{}, len(data))
			for _, val := range data {
				chunkUnique[val] = struct{}{}
			}
			results <- chunkUnique
		}(chunk)
	}

	// Close results channel once all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// mergeFloat results from all chunks
	finalUnique := make(map[string]struct{})
	for chunkResult := range results {
		for key := range chunkResult {
			finalUnique[key] = struct{}{}
		}
	}

	// Convert the unique values map to a slice
	uniqueValues := make([]string, 0, len(finalUnique))
	for key := range finalUnique {
		uniqueValues = append(uniqueValues, key)
	}

	return uniqueValues
}

func oneHotEncode(data []string) ([]float64, map[float64]string) {
	arrayUniqueValuesString(data)
	// Create a map to hold indices and categories
	indexToCategory := make(map[float64]string)
	categoryIndex := make(map[string]float64)

	index := 0.0

	// Build the categoryIndex and reverse indexToCategory maps
	for _, category := range data {
		if _, exists := categoryIndex[category]; !exists {
			categoryIndex[category] = index
			indexToCategory[index] = category
			index++
		}
	}

	// Determine the total size of the one-hot encoded data
	numCategories := int(index)
	encoded := make([]float64, len(data)*numCategories)

	// Populate the flattened one-hot encoding
	for i, category := range data {
		rowStart := i * numCategories            // Row offset in the flattened array
		colIndex := int(categoryIndex[category]) // Column index for the category
		encoded[rowStart+colIndex] = 1.0         // Set the corresponding position to 1
	}

	return encoded, indexToCategory
}
