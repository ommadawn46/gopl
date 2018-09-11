## 練習問題 7.7
20.0のデフォルト値は°を含んでいないのに、ヘルプメッセージが°を含んでいる理由を説明しなさい。

## 回答
以下の `PrintDefaults()` がヘルプメッセージ内のデフォルト値を表示している関数である。
デフォルト値として表示されているのは `Flag` 構造体のメンバ `flag.DefValue` であることが分かる。
```go
func (f *FlagSet) PrintDefaults() {
      ... 省略 ...
			if _, ok := flag.Value.(*stringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", flag.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", flag.DefValue)
			}
		}
		fmt.Fprint(f.Output(), s, "\n")
	})
}
```

`Flag` 構造体の定義を見ると、`DefValue` の型は `string` となっている。
```go
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}
```

`FlagSet` のメソッド `Var()` の中で `Flag` のインスタンスが作成されている。
`Flag.DefValue` には `value.String()` の戻り値がセットされる。
問題のプログラムtempflag.goの場合、 `value` の中身は `celsiusFlag` である。
```go
func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{name, usage, value, value.String()}
	... 省略 ...
}
```

`celsiusFlag` は `Celsius` を含む構造体であるため、 `String()` を暗黙的に満足している。
そのため、Flagへ値がセットされる際に `Celsius.String()` が呼び出され、`°C` を含んだ表示がされる。
