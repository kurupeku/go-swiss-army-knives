# go-swiss-army-knives

## Get Started

テストの実行用に以下のタスクランナーを使用する想定です。

- [Task](https://taskfile.dev/)

利用しなくともテストの実行は可能ですが、すぐにインストールできるためよろしければ導入ください。
また、タスクを実行する際はプロジェクトルートにて実行してください。

### Via Go Module

Go をローカルにインストール済みの方は OS を問わず以下のコマンドで導入可能です。

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### For Mac

```bash
brew install go-task/tap/go-task
```

### For Linux

```bash
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
```

### For Windows

```bash
choco install go-task
```

## Subject

作成する各ツールとその課題詳細は以下の通りです。

| Tool          | Desc                              |
| :------------ | :-------------------------------- |
| `curl`        | [README](./curl/README.md)        |
| `logtransfer` | [README](./logtransfer/README.md) |

## Test

各パッケージごとにテスト実行用のタスクを定義しています。
今回は同一ファイル内で実施日が異なるケースが多々存在するため、関数名を指定して実行するためのタスクを用意しています。
タスクに Args を渡す場合、 `<CMD> -- <Args>` のようにハイフン 2 つを挟むようにしてください。

### `curl` パッケージ

| Task                | Args             | Desc                                                                                     |
| :------------------ | :--------------- | :--------------------------------------------------------------------------------------- |
| `task test_curl`    | -                | `curl` パッケージ配下のテストをすべて実行します                                          |
| `task test_curl_fn` | 関数・メソッド名 | `curl` パッケージ配下のテストの内、テスト名が `<Args>` 一致するもののみ実行します        |
| `task test_lt`      | -                | `logtransfer` パッケージ配下のテストをすべて実行します                                   |
| `task test_lt_fn`   | 関数・メソッド名 | `logtransfer` パッケージ配下のテストの内、テスト名が `<Args>` 一致するもののみ実行します |

## Build

各ツール配下をすべて実装したら実際にバイナリ化して使用することができます。
出力したバイナリに Path が通っていれば、他の CLI と同じように使用することができます。

各ツールのビルド用タスクおよびバイナリ名は以下のとおりです。

| Tool          | Task              | binary  |
| :------------ | :---------------- | :------ |
| `curl`        | `task build_curl` | `scurl` |
| `logtransfer` | `task build_lt`   | `lt`    |

バイナリは各ツールのルートディレクトリ（go.mod があるところ）に出力されます。

注）出力されるバイナリは実行環境に合わせてコンパイルされます
