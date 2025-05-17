package impl

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	utils "study-session/utils/go"
)

// 文字列からinterface{}への変換（int → float64 → stringの順で試行）
func parseValue(str string) interface{} {
	str = strings.TrimSpace(str)

	// int試行
	if i, err := strconv.Atoi(str); err == nil {
		return i
	}
	// float64試行
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	}
	// fallback to string（ダブルクォートがあれば外す）
	return strings.Trim(str, `"`)
}

// loadSortTestData は入力ファイルと期待値ファイルを読み込む
func loadSortTestData(fileDir string) ([]interface{}, []interface{}, error) {
	// 入力ファイル読み込み
	inputData, err := ioutil.ReadFile(strings.Join([]string{fileDir, "input.txt"}, "/"))
	if err != nil {
		return nil, nil, fmt.Errorf("入力ファイルの読み込みに失敗しました: %v", err)
	}

	var inputArray []interface{}
	inputStr := strings.TrimSpace(string(inputData))
	inputStr = strings.Trim(inputStr, "[]")

	if inputStr != "" {
		for _, s := range strings.Split(inputStr, ",") {
			inputArray = append(inputArray, parseValue(s))
		}
	}

	// 期待値ファイル読み込み
	expectedData, err := ioutil.ReadFile(strings.Join([]string{fileDir, "expected.txt"}, "/"))
	if err != nil {
		return inputArray, nil, fmt.Errorf("期待値ファイルの読み込みに失敗しました: %v", err)
	}

	var expectedArray []interface{}
	expectedStr := strings.TrimSpace(string(expectedData))
	expectedStr = strings.Trim(expectedStr, "[]")

	if expectedStr != "" {
		for _, s := range strings.Split(expectedStr, ",") {
			expectedArray = append(expectedArray, parseValue(s))
		}
	}

	return inputArray, expectedArray, nil
}

// MeasureSortPerformance はSortの性能と正当性を計測する
func MeasureSortPerformance(fileDir string, iterations int) map[string]interface{} {
	array, expectedOutput, err := loadSortTestData(fileDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sorter := &SortImplementation{}

	fmt.Printf("Sort実装のパフォーマンス計測と正当性検証:\n")
	fmt.Printf("配列サイズ: %d\n", len(array))
	fmt.Printf("繰り返し回数: %d\n", iterations)

	var sorted []interface{}

	results := utils.MeasurePerformance("Sort", func() {
		for i := 0; i < iterations; i++ {
			arrayCopy := make([]interface{}, len(array))
			copy(arrayCopy, array)

			sorted = sorter.Sort(arrayCopy)
			if iterations == 1 {
				fmt.Printf("ソート前の先頭5要素: ")
				for j := 0; j < 5 && j < len(array); j++ {
					fmt.Printf("%v ", array[j])
				}
				fmt.Println()

				fmt.Printf("ソート後の先頭5要素: ")
				for j := 0; j < 5 && j < len(sorted); j++ {
					fmt.Printf("%v ", sorted[j])
				}
				fmt.Println()
			}
		}
	})

	valid := utils.VerifyResult("Sort", sorted, expectedOutput)
	results["valid"] = valid

	return results
}
