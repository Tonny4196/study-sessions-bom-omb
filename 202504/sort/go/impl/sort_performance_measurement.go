package impl

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	utils "study-session/utils/go"
)

// loadSortTestData は入力ファイルと期待値ファイルを読み込む
func loadSortTestData(fileDir string, dataType string) ([]interface{}, []interface{}, error) {
	// 入力データの読み込み
	inputData, err := ioutil.ReadFile(strings.Join([]string{fileDir, "input.txt"}, "/"))
	if err != nil {
		return nil, nil, fmt.Errorf("入力ファイルの読み込みに失敗しました: %v", err)
	}

	// 入力データのパースと配列への変換
	var array []interface{}
	inputStr := strings.TrimSpace(string(inputData))
	inputStr = strings.Trim(inputStr, "[]")
	if inputStr != "" {
		for _, dataStr := range strings.Split(inputStr, ",") {
			dataStr = strings.TrimSpace(dataStr)
			if dataType == "int" {
				num, err := strconv.Atoi(dataStr)
				if err != nil {
					return nil, nil, fmt.Errorf("数値のパースに失敗しました: %v", err)
				}
				array = append(array, num)
			} else if dataType == "string" {
				array = append(array, dataStr)
			}
		}
	}

	// 期待値の読み込み
	expectedData, err := ioutil.ReadFile(strings.Join([]string{fileDir, "expected.txt"}, "/"))
	if err != nil {
		return array, nil, fmt.Errorf("期待値ファイルの読み込みに失敗しました: %v", err)
	}

	// 期待値のパースと配列への変換
	var expectedOutput []interface{}
	expectedStr := strings.TrimSpace(string(expectedData))
	expectedStr = strings.Trim(expectedStr, "[]")
	if expectedStr != "" {
		for _, dataStr := range strings.Split(expectedStr, ",") {
			dataStr = strings.TrimSpace(dataStr)
			if dataType == "int" {
				num, err := strconv.Atoi(dataStr)
				if err != nil {
					return array, nil, fmt.Errorf("期待値のパースに失敗しました: %v", err)
				}
				expectedOutput = append(expectedOutput, num)
			} else if dataType == "string" {
				expectedOutput = append(expectedOutput, dataStr)
			}
		}
	}

	return array, expectedOutput, nil
}

// MeasureSortPerformance はSortの性能と正当性を計測する
func MeasureSortPerformance(fileDir string, iterations int, dataType string) map[string]interface{} {
	var err error
	array, expectedOutput, err := loadSortTestData(fileDir, dataType)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sorter := &SortImplementation[interface{}]{}

	fmt.Printf("Sort実装のパフォーマンス計測と正当性検証:\n")
	fmt.Printf("配列サイズ: %d\n", len(array))
	fmt.Printf("繰り返し回数: %d\n", iterations)

	var sorted []interface{}

	// 処理時間とメモリ使用量を計測
	results := utils.MeasurePerformance("Sort", func() {
		for i := 0; i < iterations; i++ {
			// 配列のコピーを作成
			arrayCopy := make([]interface{}, len(array))
			copy(arrayCopy, array)

			sorted = sorter.Sort(arrayCopy, func(a, b interface{}) bool {
				// less 関数でタイプに応じて比較
				switch v := a.(type) {
				case int:
					return v < b.(int)
				case string:
					return v < b.(string)
				default:
					return false
				}
			})
			if iterations == 1 {
				// ソート前とソート後の最初の5要素を表示
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

	// 正当性検証
	valid := utils.VerifyResult("Sort", sorted, expectedOutput)
	results["valid"] = valid

	return results
}
