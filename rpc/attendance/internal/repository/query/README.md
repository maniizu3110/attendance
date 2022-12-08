# Query Builder 仕様書

すべての実装は[Gorm](https://gorm.io/ja_JP/docs/)に従う．

# 使用方法

## 注意点

- **preload と join は同時に使用できない**
- ネストした型（例：User.Parent）を使用できるかどうか注意（下記図参照）

|             | ネストした型のサポート | 多重ネストした型のサポート | ID 指定での取得のサポート    | 配列型のサポート |
| ----------- | ---------------------- | -------------------------- | ---------------------------- | ---------------- |
| query       | ○                      | ×                          | × 　　　　　　　　　　　　　      | ×                |
| preload     | ○                      | ○                          | ○ 　　　　　　　　　　　　       | ○                |
| joins       | ○                      | ×                          | ×                            | ×                |
| order       | ×                      | ×                          | ×                            | ×                |
| withTrashed | ×                      | ×                          | ×                            | ×                |

## query

```go

query: ["A = 10", "B != 10", "C > 0", "D is null","A includes some","B in notnull", "A in a,b,c"]
// ネストしないパターン


query: ["Store.A = 10", "Store.B != 10", "Store.C > 0", "Store.D is null","Store.A includes some","Store.B in notnull", "Store.A in a,b,c"]
// join でネストした場合は、ネストしたテーブル名をつける

```

## preload

```go

preload: ["Store","Company.Setting","SplitList.SplitShareList"]
// Example.Store の中身も一緒に取得

```

## order

```go

order:["-A", "B"]
// Example.A の大きい順, Example.B の小さい順 に並べ替えた結果を取得
// **順番は保証される**

```

## limit

```go
limit: 10
// 先頭 10 件を取得
// offsetと併用することで任意の位置からのデータが取得可能

```

## offset

```go
offset: 15
// 15 件目から取得
```

## withTrashed

```go
withTrased: true
// 論理削除されたものも取得
```

使用例

```go
params: {
    query:["A = 10", "B != 10", "C > 0", "D is null","A includes some","B in notnull", "A in a,b,c",""],
    limit: 10,
    offset: 15,
    order:["-A", "B"],
    preload: ["Store","Company.Setting"],
    withTrased: true,
    joins: ["Store"],
}
```

## 詳細

より詳しい使用方法は query_test.go を参照してください

# 1. Preload

## 1.1. 概要

指定された関連データを事前読み込みする為のクエリを作成．

## 1.2. 使い方

`ApplyPreloads`に以下の 3 つ（場合によれば 4 つ）のデータを渡すと，中で適切に Preload クエリが作成され，処理する Gorm の DB に，作成された Preload クエリが適用される．

> 引数：`model T, preloadFields []string, query *gorm.DB`
>
> `model`:モデルデータ（非 nil なものであれば中身が空でもよい）
>
> `preloadFields`:事前読み込みしたいメンバーを格納した配列．`FetchConfig`の`Preload`配列を入れる．
>
> `query`:処理する対象の Gorm DB．

ただし，埋め込まれた構造体 ( embedded struct ) のデータをフィールドに持つものは，その構造体の中身のフィールドは Preload の対象とならない点に注意．
配列に指定するものはフィールド名．大文字小文字も正確に区別して指定する．

## 1.3. 使用例

次の様な構造体に対して`Preload`を適用する場合のデザインパターンとアンチパターンを，以下に列挙する．

```go
type MainStruct {
    ID        int
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
    CreatorID *uint
    EmbededStruct // MainStructに埋め込まれた構造体（中身のFieldはPreload対象外）
    Data NotEmbededStruct // MainStructに埋め込まれていない構造体（中身のFieldはPreloadの対象）
}

type EmbededStruct {
    Name string
    Age  uint32
}

type NotEmbededStruct {
    Address   string
    WorkPlace string
}
```

> ### 1.3.1. デザインパターン:ok_hand:
>
> 渡す`Preload`配列の例のみを示す．
>
> ```go
> members := FetchConfig{
>     Preload: []string{"Data"}
> }
> ```
>
> ### 1.3.2. アンチパターン:x:
>
> 以下の様な記述はエラーが出るので避ける．
>
> ```go
> // パターン１：フィールド名に従っていない．大文字小文字を正確に分ける必要がある．
> members := FetchConfig{
>     Preload: []string{"data"} // DataはDataと記述しなければならない．
> }
>
> // パターン２：埋め込まれた構造体に対してPreloadしよと試みる．
> members := FetchConfig{
>     Preload: []string{"Name"} // 埋め込まれていない構造体（またはスライス）に対してのみPreloadが行える．
> }
> ```

# 2. Order

## 2.1. 概要

指定されたフィールドでソートする為のクエリを作成．

## 2.2. 使い方

`ApplyOrdersToDB`に以下の 3 つのデータを渡すと，中で適切に Order クエリが作成され，処理する Gorm の DB に，作成された Order クエリが適用される．

> 引数：`model T, members []string, query *gorm.DB`
>
> `model`:モデルデータ（非 nil なものであれば中身が空でもよい）
>
> `orderFields`:ソートしたいメンバーを格納した配列．`FetchConfig`の`Order`配列を入れる．
>
> `query`:処理する対象の Gorm DB．

ただし，埋め込まれていない構造体 ( not embedded struct ) のデータをフィールドに持つものは，その構造体の中身のフィールドは Order の対象とならない点に注意．

配列に指定するものはフィールド名．大文字小文字も正確に区別して指定する．

フィールド名だけを指定する場合，デフォルトで昇順ソートになる．降順にしたい場合はフィールド名の接頭辞として`-`をつける（ex. `Score`フィールドを降順にしたいなら`-Score`）．

## 2.3. 使用例

次の様な構造体に対して`Order`を適用する場合のデザインパターンとアンチパターンを，以下に列挙する．

```go
type MainStruct {
    ID        int
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
    CreatorID *uint
    EmbededStruct // MainStructに埋め込まれた構造体（中身のFieldはOrder対象）
    Data NotEmbededStruct // MainStructに埋め込まれていない構造体（中身のFieldはOrder対象外）
}

type EmbededStruct {
    Name string
    Age  uint32
}

type NotEmbededStruct {
    Address   string
    WorkPlace string
}
```

> ### 2.3.1. デザインパターン:ok_hand:
>
> 渡す`Order`配列の例のみを示す．例えば，Name は昇順，Age は降順でソートしたい場合
>
> ```go
> members := FetchConfig{
>     Order: []string{"Name", "-Age"}
> }
> ```
>
> ### 2.3.2. アンチパターン:x:
>
> 以下の様な記述はエラーが出るので避ける．
>
> ```go
> // パターン１：フィールド名に従っていない．大文字小文字を正確に分ける必要がある．
> members := FetchConfig{
>     Order: []string{"name", "-age"} // NameはName，AgeはAgeと記述しなければならない．
> }
>
> // パターン２：埋め込まれていない構造体のフィールドでOrderしようとする．
> members := FetchConfig{
>     Order: []string{"Address", "-Age"}
> }
>
> // パターン３：prefixとフィールド名が離れている．
> members := FetchConfig{
>     Order: []string{"- Age"} // 実はこれでもノーエラーで適切に処理されるが，見栄えが悪いのでアンチパターンで...
> }
> ```

# 3. Query

## 3.1. 概要

指定されたフィールドを指定条件で取得する為のクエリを作成．

## 3.2. Query の種類とそれに対応する使い方・使用例

`ApplyQueriesToDB`に以下の 3 つのデータを渡すと，中で適切に Query クエリが作成され，処理する Gorm の DB に，作成された Query クエリが適用される．

> 引数：`model T, conditions []string, query *gorm.DB`
>
> `model`:モデルデータ（非 nil なものであれば中身が空でもよい）
>
> `conditions`:指定フィールドに対する指定条件を格納した配列．`FetchConfig`の`Query`配列を入れる．
>
> `query`:処理する対象の Gorm DB．

また，以下全ての Query に対する使用例は，次の構造体を用いること想定する．

```go
type MainStruct {
    ID        int
　　CreatedAt time.Time
　　UpdatedAt time.Time
　　DeletedAt gorm.DeletedAt
　　CreatorID *uint
    EmbededStruct // MainStructに埋め込まれた構造体（中身のFieldはQueryの対象）
    Data NotEmbededStruct // MainStructに埋め込まれていない構造体（中身のFieldはQuery対象外）
}

type EmbededStruct {
    Name string
    Age  uint32
}

type NotEmbededStruct {
    Address   string
    WorkPlace string
}
```

### 3.2.1. Comparition（ `<=` , `>=` , `=` , `!=` , `<` , `>` ）

比較演算であるが，改めて各定義を記述しておく．
| 演算子 | 使用例 | 機能 |
| ---- | ---- | ---- |
| `<=` | A <= B | A は B 以下 |
| `>=` | A >= B | A は B 以上 |
| `=` | A = B | A と B は等しい |
| `!=` | A != B | A と B は等しくない |
| `>` | A > B | A は B より大きい |
| `<` | A < B | A は B より小さい |

> ### 3.2.1.1. デザインパターン:ok_hand:
>
> ```go
> members := FetchConfig {
>     Query: []string{"Name = hoge", "Age >= 15"}
> }
> ```
>
> ### 3.2.1.2. アンチパターン:x:
>
> ```go
> // パターン１：埋め込まれていない構造体に対してのQueryを試みる
> members := FetchConfig {
>     Query: []string{"Address != huge"}
> }
>
> // パターン２：不正なクエリ
> members := FetchConfig {
>     Query: []string{"Name ! = gopher", "Age > = 15", "CreatedAt >> 2022-07-23T15:04:05Z07:00"} // <,>,!と=は切りはなしてはいけない．また，存在しない演算子を使用するのも勿論不正
> }
> ```

### 3.2.2. `IN`

| 演算子 | 使用例              | 機能                         |
| ------ | ------------------- | ---------------------------- |
| `IN`   | A in (x, y, z, ...) | A は括弧内のいずれかと等しい |

> ### 3.2.2.1. デザインパターン:ok_hand:
>
> ```go
> members := FetchConfig {
>     Query: []string{"Name in hoge,fuga,foo,bar"}
> }
> ```
>
> ### 3.2.2.2. アンチパターン:x:
>
> ```go
> // パターン１：埋め込まれていない構造体に対してのQueryを試みる
> members := FetchConfig {
>     Query: []string{"Address in hoge,fuga"}
> }
>
> // パターン２：不正なクエリ
> members := FetchConfig {
>     Query: []string{"Name in gopher in golang"} // in hoge,fuga,...の形を用いること
> }
> ```

### 3.2.3. `NOT IN`

| 演算子   | 使用例                  | 機能                             |
| -------- | ----------------------- | -------------------------------- |
| `NOT IN` | A not in (x, y, z, ...) | A は括弧内のいずれとも等しくない |

> ### 3.2.3.1. デザインパターン:ok_hand:
>
> ```go
> members := FetchConfig {
>     Query: []string{"Name not in hoge,fuga,foo,bar"}
> }
> ```
>
> ### 3.2.3.2. アンチパターン:x:
>
> ```go
> // パターン１：埋め込まれていない構造体に対してのQueryを試みる
> members := FetchConfig {
>     Query: []string{"Address not in hoge,fuga"}
> }
>
> // パターン２：不正なクエリ
> members := FetchConfig {
>     Query: []string{"Name not in gopher,gogogo not in golang"} // not in hoge,fuga,...の形を用いること
> }
> ```

### 3.2.4. `LIKE`

| 演算子 | 使用例   | 機能                                                               |
| ------ | -------- | ------------------------------------------------------------------ |
| `LIKE` | A LIKE B | A は正規表現を使用して文字列検索した場合，B を部分文字列として含む |

> ### 3.2.4.1. デザインパターン:ok_hand:
>
> `LIKE`の場合は`includes`を持ちいる．また，指定できる部分文字列は 1 つのみである．
>
> ```go
> members := FetchConfig {
>     Query: []string{"Name includes hoge"}
> }
> ```
>
> ### 3.2.4.2. アンチパターン:x:
>
> ```go
> // パターン１：埋め込まれていない構造体に対してのQueryを試みる
> members := FetchConfig {
>     Query: []string{"Address includes hoge"}
> }
>
> // パターン２：指定部分文字列を複数指定する
> members := FetchConfig {
>     Query: []string{"Name includes gopher,gopher"} // 指定できるのは1つのみ
> }
>
> // パターン３：Typo
> members := FetchConfig {
>     Query: []string{"Name include gopher,gopher"} // includes（三人称）である点に注意
> }
> ```

### 3.2.5. `NULL` or `NOT NULL`

| 演算子        | 使用例        | 機能             |
| ------------- | ------------- | ---------------- |
| `is NULL`     | A is NULL     | A は NULL である |
| `is NOT NULL` | A is NOT NULL | A は NULL でない |

> ### 3.2.5.1. デザインパターン:ok_hand:
>
> ```go
> members := FetchConfig {
>     Query: []string{"DeletedAt is NULL", "UpdatedAt is NOTNULL"}
> }
> ```
>
> ### 3.2.5.2. アンチパターン:x:
>
> ```go
> // パターン１：埋め込まれていない構造体に対してのQueryを試みる
> members := FetchConfig {
>     Query: []string{"Address is NULL"}
> }
>
> // パターン２：小文字を交えてしまう
> members := FetchConfig {
>     Query: []string{"Name is Null"} // 全て大文字（実装で小文字を含めるようにすることも可能だが，冗長になるので辞めた）
> }
>
> // パターン３：NOT と NULLを切り離してしまう
> members := FetchConfig {
>     Query: []string{"DeletedAt is NOT NULL"} // NOTとNULLは1つの文字列にまとめる（NOTNULLとする）
> }
> ```
