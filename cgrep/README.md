# 課題 2:簡易 内容一致ファイル名検索コマンドを実装する

## コマンド仕様

### 説明

- ファイルの内容を対象に一致したファイル名を出力する
- 出力されるファイル名はカレントディレクトリからの相対パスで表記
- 引数を正規表現として解釈し、ファイル内の各行で一致を検証する
- ディレクトリ内捜索を goroutine にてマルチスレッド化して非同期で行う
- 検索範囲はデフォルトでカレントディレクトリ配下
- フラグにて検索ルートとするディレクトリを指定可能
- フラグにて一致した行も一緒に表示することが可能
- ノイズになるので `.git` ディレクトリ配下は検索対象から除外

### コマンド引数・フラグ

- 引数は検索用の正規表現
- フラグ（`-a`または`--aaa`という形式で設定するコマンドオプション）は以下の通り
  - `-d`(`--dir`): 検索ルートを指定（デフォルトは `./`）
  - `-c`(`--with-content`): 一致した行を合わせて標示させる
- 以下はコマンドのヘルプ表示

```bash

```

### 使用例

#### デフォルト動作

```bash
$ go run main.go hoge
dir1/filename1.md
filename2.txt
```

#### 検索ディレクトリを指定

```bash
$ go run main.go -d ./dir1 hoge
dir1/filename1.md
```

#### 一致した行も表示

```bash
$ go run main.go hoge
dir1/filename1.md
24: My name is **hoge**.
128: no hoge no life

filename2.txt
35: What is hoge?
```

## 実装課題

- コマンド引数・フラグを受け取る部分は実装済み
- コマンド引数の個数とパスのレンダリングに関するバリデーションは実装済み
- 検索用構造体のファクトリ関数、相対パスへの変換関数は実装済み
- 検索結果への保存・取得用関数、ファイル一覧出力用メソッドは実装済み
- エラー保存・取得用関数は実装済み
- 未実装となっている以下の処理を実装する
  - [1 週目] 配下のディレクトリ・ファイル検索機能の実装
  - [2 週目] 検索結果のレンダリング & コマンド実行時のメイン処理の実装

### 1 週目：配下のディレクトリ・ファイル検索機能の実装

- 対応ファイル：`cgrep/search/dir.go`
- 実装内容：`func (d *dir) Search()`で配下の各ディレクトリの検索を非同期で実行し、ファイル郡を`func (d *dir) GrepFiles() error`を使って内容一致検索できるように検索ロジックを実装する
- 実装対象メソッド・実装条件
  - `func (d *dir) Search()`
    - この関数は非同期で呼び出される前提の構造になっているので、引数で受け取る `d.wg` を使用してメソッドの開始と終了を知らせられるように実装する
      - `sync.WaitGroup` の使い方を事前に確認してから実装を行う
    - `d.subDirs` がレシーバと同じ `dir` 構造体のポインターのスライスになっているので、それらの `func (d *dir) Search()` を **非同期で** 実行する
    - 配下のファイル郡の検索もこの関数内で実行する
      - ファイルの検索ロジック自体は次の `func (d *dir) GrepFiles() error` で実装する
  - `func (d *dir) GrepFiles() error`
    - `d.fileFullPaths` が配下にあるファイルへのフルパスになっているので
      - ファイルを開いて
      - 内容を読み取って
      - 検索用の正規表現と一致するかを検証して
      - 一致したらファイル名と一致した行番号と行の文字列を記録して
      - ファイルを閉じる
    - ファイル名はカレントディレクトリからの相対パスとともに出力する
      - e.g.) `testdata/dir/text.txt` `../curl/client/client.go`
    - 内容の保存は `result.Set(fileName, txt string, no int)` へ渡すと保存される
    - エラーが発生したらそのタイミングでリターンする

### 2 週目：検索結果のレンダリング & コマンド実行時のメイン処理の実装

- 対応ファイル：`cgrep/result/result.go` と `cgrep/cmd/root.go`
- 実装内容：`func RenderFiles(w io.Writer)` と `func RenderWithContent(w io.Writer)`、`func ExecSearch(fullPath, regexpWord string) error` と `func Render(w io.Writer)`
- 実装対象メソッド・実装条件
  - `func RenderFiles(w io.Writer)`
    - 一致したファイル名を標準出力にすべて出力するための関数
    - 標準出力は `w io.Writer` として渡される想定
    - ファイルのみを出力するメソッドは `func (r *Result) Files() []string` として実装済み
    - 出力結果はソートされている想定（`func (r *Result) Files() []string` を使えばソート済み）
    - 単純にファイル名を改行区切りで出力すれば OK
    - 改行コードは `\n` を使用する
    - 出力フォーマットは後述
  - `func RenderWithContent(w io.Writer)`
    - `func RenderFiles(w io.Writer)` に一致した行番号と行の内容も合わせて出力する関数
    - 標準出力は `w io.Writer` として渡される想定
    - 出力結果はソートされている想定（`func (r *Result) Files() []string` を使えばソート済み）
    - 改行コードは `\n` を使用する
    - 出力フォーマットは後述
  - `func ExecSearch(fullPath, regexpWord string) error`
    - `fullPath` は検索ルートへの絶対パスが渡される
      - デフォルトではカレントディレクトリになっている
    - `regexpWord` は文字列で渡されるので正規表現オブジェクトに変換する必要がある
    - `search.New(wg *sync.WaitGroup, fullPath string, re *regexp.Regexp)` で `search.Dir` の生成と `search.Dir.Search()` の実行が必要
    - `search.Dir.Search()` は **非同期で実行** する必要がある
    - 実行は非同期だが、結果のレンダリングはすべての非同期関数が終了してからにする必要がある
    - エラーが発生したらすぐに返す
  - `func Render(w io.Writer)`
    - 結果を標準出力に出力するための関数
    - 標準出力は `w io.Writer` として渡される想定
    - グローバル変数 `withContent` が `true` の場合は一致した内容も表示、 `false` の場合はファイル名のみを表示できるようにする

#### RenderFiles の出力フォーマット

```bash
<ファイル名>
<ファイル名>
```

#### RenderWithContent の出力フォーマット

```bash
<ファイル名>
<一致した行番号>:<半角スペース><一致した行の内容>
<一致した行番号>:<半角スペース><一致した行の内容>
<空行>
<ファイル名>
<一致した行番号>:<半角スペース><一致した行の内容>
<一致した行番号>:<半角スペース><一致した行の内容>
```
