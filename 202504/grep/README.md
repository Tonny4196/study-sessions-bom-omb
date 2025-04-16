# Grep 実装ガイド

## 📂 ファイル構造

```
grep/
├── go/
│   ├── main.go
│   └── impl/
│       ├── grep_implementation.go     # 実装ファイル
│       └── grep_performance_measurement.go
│
├── javascript/
│   ├── grep_implementation.js         # 実装ファイル
│   └── grep_performance_measurement.js
│
├── ruby/
│   ├── grep_implementation.rb         # 実装ファイル
│   └── grep_performance_measurement.rb
│
└── test_cases/                        # テストケース
    ├── case1/
    ├── case2/
    └── ...
```

## 🔍 ファイルの役割

| ファイル名 | 説明 |
|------------|------|
| **grep_implementation.xx** | grep処理を行うためのメイン実装ファイル。ここに実装を記述します |
| **grep_performance_measurement.xx** | 実装したgrepの検証とパフォーマンス計測を行うためのファイル |

## ⚙️ 実装方法

実装は各言語の **grep_implementation** ファイルに記述してください。

- ファイル分割は自由に行っていただいて構いません
- ただし、インターフェイスは **grep_implementation** で定義されているものに合わせてください
  - インターフェイスが一致しないと、正確なパフォーマンス測定および正当性検証ができません

## 🧪 テストと計測方法

### Go

```bash
go run <your_path>/grep/go/main.go "<your_path>/grep/test_cases/<テストケース>"
```

**例:**
```bash
go run grep/go/main.go "./grep/test_cases/case1"
```

### JavaScript (Node.js)

```bash
node <your_path>/grep/javascript/grep_performance_measurement.js "<your_path>/grep/test_cases/<テストケース>"
```

**例:**
```bash
node grep/javascript/grep_performance_measurement.js "./grep/test_cases/case1"
```

### Ruby

```bash
ruby <your_path>/grep/ruby/grep_performance_measurement.rb "<your_path>/grep/test_cases/<テストケース>"
```

**例:**
```bash
ruby grep/ruby/grep_performance_measurement.rb "./grep/test_cases/case1"
```

## 📝 実装のポイント

- 正確なパターンマッチング: テキストファイル内の行から指定されたパターンを検索
- 効率的な検索アルゴリズム: 大きなファイルでもパフォーマンスが低下しないように実装
- メモリ管理: 大きなファイルを効率的に処理できるように留意

## 💡 テストケースの構成

各テストケースディレクトリには以下のファイルが含まれています:

- `input.txt`: 検索対象のテキストファイル
- `pattern.txt`: 検索するパターン
- `expected.txt`: 期待される出力結果

これらのファイルを使用して、実装したGrepの正確性が検証されます。
