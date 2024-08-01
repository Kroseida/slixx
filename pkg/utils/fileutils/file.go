package fileutils

import (
	"os"
	"sort"
	"strings"
)

type FileInfo struct {
	Name          string
	FullDirectory string
	RelativePath  string
	CreatedAt     int64
	Directory     bool
	Size          uint64
}

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FixedPathName(path string) string {
	fixedPath := path

	// Remove double slashes
	for strings.Contains(fixedPath, "//") {
		fixedPath = strings.ReplaceAll(fixedPath, "//", "/")
	}

	return fixedPath
}

func ParentDirectory(path string) string {
	var pos = strings.LastIndex(path, "/")
	if pos == -1 {
		return path
	}
	return path[0:pos]
}

func SplitArrayBySize(array []FileInfo, n int) ([][]FileInfo, []uint64) {
	// Sort the array in descending order based on the "Size" field
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Size > array[j].Size
	})

	result := make([][]FileInfo, n)
	totalSizes := make([]uint64, n)

	for i := range array {
		if len(totalSizes) == 0 {
			break
		}

		// Find the index of the chunk with the smallest total size so far
		minIndex := 0
		minSize := totalSizes[0]
		for j := 1; j < n; j++ {
			if totalSizes[j] < minSize {
				minIndex = j
				minSize = totalSizes[j]
			}
		}

		// Append the current element to the chunk with the smallest total size
		result[minIndex] = append(result[minIndex], array[i])
		// Update the total size of the chunk
		totalSizes[minIndex] += array[i].Size
	}

	return result, totalSizes
}

func SortByLength(array []FileInfo) []FileInfo {
	sort.SliceStable(array, func(i, j int) bool {
		return len(array[i].RelativePath) > len(array[j].RelativePath)
	})
	return array
}
