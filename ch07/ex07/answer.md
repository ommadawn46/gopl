## 練習問題 7.7
20.0のデフォルト値は°を含んでいないのに、ヘルプメッセージが°を含んでいる理由を説明しなさい。

## 回答
下記のコード中の`PrintDefaults()`がヘルプメッセージのデフォルト値をアウトプットしているメソッドである。
コードを見ると`s += fmt.Sprintf(" (default %q)", flag.DefValue)`でデフォルト値を表す文字列を作成している。この部分から、デフォルト値としてアウトプットされているのは、`Flag`構造体のメンバ`DefValue`の値であることが読み取れる。
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

`Flag`構造体の定義は以下の通りである。`DefValue`の型は`string`となっている。
```go
type Flag struct {
	Name     string // name as it appears on command line
	Usage    string // help message
	Value    Value  // value as set
	DefValue string // default value (as text); for usage message
}
```

`Flag`のインスタンスは、下記のコード中のメソッド`Var()`で作成される。
コードから、`Flag.DefValue`には`value.String()`の戻り値がセットされていることがわかる。
問題のプログラムtempflag.goの場合、`value`の中身には`celsiusFlag`がセットされた状態でこのメソッドが呼ばれることになる。
```go
func (f *FlagSet) Var(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{name, usage, value, value.String()}
	... 省略 ...
}
```

`celsiusFlag`は`Celsius`を含む構造体であるため、`celsiusFlag.String()`を呼び出した際には暗黙的に`Celsius.String()`が呼び出される。
以下のコードの通り、`Celsius.String()`は°を含んだ文字列を返すため、最終的なヘルプメッセージにも°が含まれることになる。

```go
func (c Celsius) String() string    { return fmt.Sprintf("%g°C", c) }
```
