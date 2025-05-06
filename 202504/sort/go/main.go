package main

import (
	"fmt"
	"os"
	"study-session/sort/go/impl"
)

// ===============================================
// メイン関数
// ===============================================

func main() {
	// プログラム開始時に表示するメッセージ
	fmt.Println("==============================")
	fmt.Println("Sort性能計測と正当性検証")
	fmt.Println("==============================")

	// ソートの計測と検証を開始
	fmt.Println("Sort実装のテスト")
	fileDir := os.Args[1]                                            // コマンドライン引数からファイルディレクトリを取得
	sortResults := impl.MeasureSortPerformance(fileDir, 1, "string") // ソートのパフォーマンスを計測

	// 検証結果の要約を表示
	fmt.Println("\n==============================")
	fmt.Println("テスト結果サマリー")
	fmt.Println("==============================")

	// ソート結果の検証
	sortValid := false
	if sortResults != nil {
		// "valid"キーの値がbool型か確認し、bool型の場合に変換
		if valid, ok := sortResults["valid"].(bool); ok {
			sortValid = valid // bool型に変換できた場合、変数に結果を格納
		}
	}

	// 結果を表示
	fmt.Printf("Sort: %s\n", boolToCheckmark(sortValid))
}

// boolToCheckmark はブール値をチェックマーク文字列に変換
func boolToCheckmark(b bool) string {
	if b {
		// 成功の場合は「成功 ✓」を返す
		return "成功 ✓"
	}
	// 失敗の場合は「失敗 ✗」を返す
	return "失敗 ✗"
}
