package query

type FetchConfig struct {
	Limit       int
	Offset      int
	WithTrashed bool     // 削除されたものを含むかどうか　default: false
	Preload     []string // list of メンバ変数名の文字列(ネスト可/クエリ不可)
	Joins       []string // list of メンバ変数名の文字列(ネスト不可/クエリ可)
	Order       []string // list of {+|-}メンバ変数名の文字列
	// list of {
	//    メンバ変数名の文字列 <op> <value> |
	//    メンバ変数名の文字列 {in | not in } <value1>,<value2>,... |
	//    メンバ変数名の文字列 is { NULL | NOTNULL }
	// }
	// op is one of  =, !=, <, <=, >, >=
	Query []string
}
