package main

import "fmt"

func main() {
	// 假设有三个切片
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5}
	slice3 := []int{7, 8, 9, 10, 11, 12, 13}

	// 使用append将每个切片添加到一个新的多维切片中
	var slices [][]int
	slices = append(slices, slice1)
	slices = append(slices, slice2)
	slices = append(slices, slice3)

	// 输出结果
	fmt.Println(slices) // 输出: [[1 2 3] [4 5] [7 8 9 10 11 12 13]]
	fmt.Println(len(slices))
	for _, v := range slices {
		fmt.Println(v)
	}

}
