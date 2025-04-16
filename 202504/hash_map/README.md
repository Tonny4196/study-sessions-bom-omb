# HashMap 実装ガイド

## 📂 ファイル構造

```
hashmap/
├── go/
│   ├── main.go
│   └── impl/
│       ├── hashmap_implementation.go     # 実装ファイル
│       └── hashmap_performance_measurement.go
│
├── javascript/
│   ├── hashmap_implementation.js         # 実装ファイル
│   └── hashmap_performance_measurement.js
│
├── ruby/
│   ├── hashmap_implementation.rb         # 実装ファイル
│   └── hashmap_performance_measurement.rb
│
└── test_cases/                           # テストケース
    ├── case1/
    ├── case2/
    └── ...
```

## 🔍 ファイルの役割

| ファイル名 | 説明 |
|------------|------|
| **hashmap_implementation.xx** | HashMapのメイン実装ファイル。ここに実装を記述します |
| **hashmap_performance_measurement.xx** | 実装したHashMapの検証とパフォーマンス計測を行うためのファイル |

## ⚙️ 実装方法

実装は各言語の **hashmap_implementation** ファイルに記述してください。

- ファイル分割は自由に行っていただいて構いません
- ただし、インターフェイスは **hashmap_implementation** で定義されているものに合わせてください
  - インターフェイスが一致しないと、正確なパフォーマンス測定および正当性検証ができません

## 🔑 必須インターフェイス

実装するHashMapには、以下の4つのメソッドを必ず含めてください：

1. **put(key, value)**: キーと値のペアを格納
   - 既存のキーの場合は値を更新
   - 新規のキーの場合は追加

2. **get(key)**: キーに対応する値を取得
   - キーが存在しない場合はnull/nil/undefinedを返す

3. **remove(key)**: キーに対応するエントリを削除
   - 削除に成功した場合はtrue/削除した値、失敗した場合はfalse/nilを返す

4. **all_entries()**: 全てのキーと値のペアを取得（テスト用）
   - 戻り値はハッシュ/オブジェクト/マップ形式で全てのエントリを含む

内部実装では衝突解決やリサイズなど必要な機能を含めてください。

## 🧪 テストと計測方法

### Go

```bash
go run <your_path>/hashmap/go/main.go "<your_path>/hashmap/test_cases/<テストケース>"
```

**例:**
```bash
go run hashmap/go/main.go "./hashmap/test_cases/case1"
```

### JavaScript (Node.js)

```bash
node <your_path>/hashmap/javascript/hashmap_performance_measurement.js "<your_path>/hashmap/test_cases/<テストケース>"
```

**例:**
```bash
node hashmap/javascript/hashmap_performance_measurement.js "./hashmap/test_cases/case1"
```

### Ruby

```bash
ruby <your_path>/hashmap/ruby/hashmap_performance_measurement.rb "<your_path>/hashmap/test_cases/<テストケース>"
```

**例:**
```bash
ruby hashmap/ruby/hashmap_performance_measurement.rb "./hashmap/test_cases/case1"
```

## 📝 実装のポイント

- **固定サイズのバケット配列**: 初期サイズのバケット配列を使用してハッシュマップを実装
- **効率的なハッシュ関数**: 衝突を最小限に抑えるハッシュ関数の実装
- **衝突解決機構**: チェイニングまたはオープンアドレス法などの衝突解決手法の実装
- **O(1)時間複雑度**: 効率的な検索、挿入、削除操作の実現
- **動的リサイズ**: 負荷係数に基づいた動的なリサイズ機能の実装
- **任意のデータ型対応**: 様々なキーと値のデータ型に対応できる設計

## 💡 テストケースの構成

各テストケースディレクトリには以下のファイルが含まれています:

- `operations.txt`: 実行する操作のリスト（JSON形式）
  - 例: `[{"action": "put", "key": "key1", "value": "value1"}, {"action": "get", "key": "key1"}, ...]`
- `expected.txt`: 操作後の期待されるハッシュマップの状態（JSON形式のキー・値オブジェクト）

これらのファイルを使用して、実装したHashMapの正確性とパフォーマンスが検証されます。

## 🔄 要求するインターフェイス
実装の際には、必ず、以下の4つの関数を実装してください。
- put: HashMapにキーバリューをセットするための関数
- get: HashMapから特定キーに対応するバリューを取得するための関数
- remove: HashMapから特定キーを削除するための関数
- all_entries: 各言語で標準実装されたHash/Map/Object形式でキーバリューペアの組み合わせを出力するための関数

以下に例を示します。

### Ruby
```ruby
class HashMapImplementation
  def put(key, value)
    # キーと値のペアを格納
  end
  
  def get(key)
    # キーに対応する値を取得
  end
  
  def remove(key)
    # キーに対応するエントリを削除
  end
  
  def all_entries
    # 全てのキーと値のペアを取得
  end
end
```

### Go
```go
type HashMapImplementation struct {
  // 内部フィールド
}

func (h *HashMapImplementation) Put(key, value string) {
  // キーと値のペアを格納
}

func (h *HashMapImplementation) Get(key string) (string, bool) {
  // キーに対応する値を取得（第2戻り値は存在フラグ）
}

func (h *HashMapImplementation) Remove(key string) bool {
  // キーに対応するエントリを削除
}

func (h *HashMapImplementation) GetAllEntries() map[string]string {
  // 全てのキーと値のペアを取得
}
```

### JavaScript
```javascript
class HashMapImplementation {
  put(key, value) {
    // キーと値のペアを格納
  }
  
  get(key) {
    // キーに対応する値を取得
  }
  
  remove(key) {
    // キーに対応するエントリを削除
  }
  
  getAllEntries() {
    // 全てのキーと値のペアを取得
  }
}
```