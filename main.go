//Task: Concurrent Matrix Addition
//Description:
//Write a program that adds two matrices concurrently using goroutines and channels.
//
//One goroutine should generate the matrices and send their elements to separate channels.
//Multiple worker goroutines should receive elements from these channels, compute the sum of corresponding elements, and send the results to a results channel.
//The main goroutine should collect and print the final matrix from the results channel.

package main

import (
	"fmt"
	"sync"
)

// Function to create a matrix
func createMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	return matrix
}

// Function to fill a matrix with random values
func fillMatrix(matrix [][]int, matrixChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	value := 1
	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = value
			matrixChan <- matrix[i][j]
			value++
		}
	}
	close(matrixChan)
}

// Function to add two matrices
func addMatrices(matrixAChan, matrixBChan, resultMatrixChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range matrixAChan {
		j := <-matrixBChan
		resultMatrixChan <- i + j
	}
}

func main() {
	rows, cols := 3, 3
	matrixA := createMatrix(rows, cols)
	matrixB := createMatrix(rows, cols)

	matrixAChan := make(chan int, rows*cols)
	matrixBChan := make(chan int, rows*cols)
	resultMatrixChan := make(chan int, rows*cols)

	var wg sync.WaitGroup
	wg.Add(2)
	go fillMatrix(matrixA, matrixAChan, &wg)
	go fillMatrix(matrixB, matrixBChan, &wg)

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go addMatrices(matrixAChan, matrixBChan, resultMatrixChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultMatrixChan)
	}()

	resultMatrix := createMatrix(rows, cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			resultMatrix[i][j] = <-resultMatrixChan
		}
	}

	fmt.Println("Matrix A:")
	for _, row := range matrixA {
		fmt.Println(row)
	}

	fmt.Println("Matrix B:")
	for _, row := range matrixB {
		fmt.Println(row)
	}

	fmt.Println("Result Matrix:")
	for _, row := range resultMatrix {
		fmt.Println(row)
	}
}
