package main

import (
	"math/rand"
	"strconv"
	"bufio"
	"fmt"
	"os"
)

func openFile(filename string) (slice []string) {
	var file, _ = os.Open(filename)
	var scan *bufio.Scanner = bufio.NewScanner(file)
	for scan.Scan() {
		slice = append(slice, scan.Text())
	}
	file.Close()
	return slice
}

func createFile(filename string, slice []string) {
	var file, _ = os.Create(filename)
	for i := range slice {
		fmt.Fprintln(file, slice[i])
	}
	file.Close()
}

func appendFile(filename, content string) {
	var file, _ = os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	fmt.Fprintln(file, content)
	file.Close()
}

func chunkSlice(slice []string, size int) (chunks [][]string) {
	for size < len(slice) {
		slice, chunks = slice[size:], append(chunks, slice[0:size:size])
	}
	return append(chunks, slice)
}

func removeString(slice []string, str string) []string {
	for i := range slice {
		if (slice[i] == str) {
			return append(slice[:i], slice[i + 1:]...)
		}
	}
	return slice
}

func randomIntString(length int) string {
	var bytes []byte = make([]byte, rand.Intn(length) + 1)
	for i := range bytes {
		bytes[i] = "0123456789"[rand.Intn(10)]
	}
	return string(bytes)
}

func formatNumber(number int64) string {
    var in string = strconv.FormatInt(number, 10)
	var out []byte = make([]byte, len(in) + (len(in) - 1) / 3)
	for i, j, k := len(in) - 1, len(out) - 1, 0; ; i, j = i - 1, j - 1 {
		out[j] = in[i]
		if (i == 0 ){
			return string(out)
		}
        if k++; (k == 3) {
			j, k = j - 1, 0
			out[j] = ','
        }
    }
}
