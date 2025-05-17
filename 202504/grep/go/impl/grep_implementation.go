package impl

import (
	"bufio"
	"log"
	"os"
	"io"
	"bytes"
)

// GrepImplementation はGrepの基本実装を提供する
type GrepImplementation struct{}

// Search はファイルから特定のパターンを検索する
func (g *GrepImplementation) Search(filePath, pattern string) []string {
	var result []string

	// ファイルを開く
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open file %s: %v", filePath, err)
		return result
	}
	defer f.Close()

	info, _ := f.Stat()

	// 行単位でスキャン
	reader := bufio.NewReaderSize(f, optimalBufSize(info.Size()))
	keyword := []byte(pattern)

	for {
		line, err := reader.ReadSlice('\n')
		if err != nil && err != io.EOF {
			log.Printf("failed to read file %s: %v", filePath, err)
			return nil
		}
		// パターンが含まれているか判定
		if bytes.Contains(line, keyword) {
			// 行の末尾の改行を削除
			if n := len(line); n > 0 && line[n-1] == '\n' {
				line = line[:n-1]
				if n > 1 && line[n-2] == '\r' {
						line = line[:n-2]
				}
			}
      result = append(result, string(line))
    }
		if err == io.EOF {
      break
    }
	}

	return result
}

// optimalBufSize はファイルサイズに基づいて最適なバッファサイズを決定する
func optimalBufSize(fileSize int64) int {
    switch {
    case fileSize < 64*1024:
        return 4 * 1024
    case fileSize < 5*1024*1024:
        return 64 * 1024
    default:
        return 256 * 1024
    }
}
