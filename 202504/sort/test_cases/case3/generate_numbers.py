import random

# 0〜9999の範囲でランダムな整数を10,000個生成
numbers = [random.randint(0, 9999) for _ in range(10000)]

# 元のリストを書き出す
with open('input.txt', 'w', encoding='utf-8') as f:
    f.write(str(numbers))

# ソートして別リストに格納
sorted_numbers = sorted(numbers)

# ソート済みリストを書き出す
with open('expected.txt', 'w', encoding='utf-8') as f:
    f.write(str(sorted_numbers))

print("生成リストを input.txt に、ソート結果を expected.txt に書き出しました。")
