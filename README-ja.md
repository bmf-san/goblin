[English](https://github.com/bmf-san/goblin) [日本語](https://github.com/bmf-san/goblin/blob/master/README-ja.md)

# goblin
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GitHub release](https://img.shields.io/github/release/bmf-san/goblin.svg)](https://github.com/bmf-san/goblin/releases)
[![CircleCI](https://circleci.com/gh/bmf-san/goblin/tree/master.svg?style=svg)](https://circleci.com/gh/bmf-san/goblin/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmf-san/goblin)](https://goreportcard.com/report/github.com/bmf-san/goblin)
[![codecov](https://codecov.io/gh/bmf-san/goblin/branch/master/graph/badge.svg?token=ZLOLQKUD39)](https://codecov.io/gh/bmf-san/goblin)
[![GitHub license](https://img.shields.io/github/license/bmf-san/goblin)](https://github.com/bmf-san/goblin/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/bmf-san/goblin.svg)](https://pkg.go.dev/github.com/bmf-san/goblin)
[![Sourcegraph](https://sourcegraph.com/github.com/bmf-san/goblin/-/badge.svg)](https://sourcegraph.com/github.com/bmf-san/goblin?badge)

トライ木をベースにしたGo製のHTTP Routerです。

<img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/d654ddf2b81c2b4123684f93071af0cf559eb0b5.png" alt="goblin" title="goblin" width="250px">

このロゴは[gopherize.me](https://gopherize.me/gopher/d654ddf2b81c2b4123684f93071af0cf559eb0b5)で作成しました。

# 目次
- [goblin](#goblin)
- [目次](#目次)
- [特徴](#特徴)
- [インストール](#インストール)
- [例](#例)
- [使い方](#使い方)
  - [メソッドベースのルーティング](#メソッドベースのルーティング)
  - [名前付きパラメータのルーティング](#名前付きパラメータのルーティング)
  - [正規表現を使ったルーティング](#正規表現を使ったルーティング)
  - [ミドルウェア](#ミドルウェア)
  - [カスタム可能なエラーハンドラー](#カスタム可能なエラーハンドラー)
  - [デフォルトOPTIONSハンドラー](#デフォルトoptionsハンドラー)
- [ベンチマークテスト](#ベンチマークテスト)
- [設計](#設計)
- [コントリビューション](#コントリビューション)
- [スポンサー](#スポンサー)
- [ライセンス](#ライセンス)
  - [作者](#作者)

# 特徴
- Go1.20 >= 1.16
- トライ木をベースとしたシンプルなデータ構造
- 軽量
  - Lines of codes:2428
  - Package size: 140K
- 標準パッケージ以外の依存性なし
- net/httpとの互換性
- net/httpの[Servemux](https://pkg.go.dev/net/http#ServeMux)よりも高機能
  - メソッドベースのルーティング
  - 名前付きパラメータのルーティング
  - 正規表現を使ったルーティング
  - ミドルウェア
  - カスタム可能なエラーハンドラー
  - デフォルトOPTIONSハンドラー
- 0allocs
  - 静的なルーティングにおいて0allocsを達成
  - 名前付きルーティングについては3allocs程度
    - パラメータのslice生成やパラメータをcontextに格納する部分でヒープ割当が発生

# インストール
```sh
go get -u github.com/bmf-san/goblin
```

# 例
サンプルの実装を用意しています。

[_examples](https://github.com/bmf-san/goblin/blob/master/_examples)をご参照ください。

# 使い方
## メソッドベースのルーティング
任意のHTTPメソッドに基づいてルーティングを定義することができます。

以下のHTTPメソッドをサポートしています。
`GET/POST/PUT/PATCH/DELETE/OPTIONS`

```go
r := goblin.NewRouter()

r.Methods(http.MethodGet).Handler(`/`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/")
}))

r.Methods(http.MethodGet, http.MethodPost).Handler(`/methods`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        fmt.Fprintf(w, "GET")
    }
    if r.Method == http.MethodPost {
        fmt.Fprintf(w, "POST")
    }
}))

http.ListenAndServe(":9999", r)
```

## 名前付きパラメータのルーティング
名前付きパラメータ(`:paramName`)を使ったルーティングを定義することができます。

```go
r := goblin.NewRouter()

r.Methods(http.MethodGet).Handler(`/foo/:id`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))

r.Methods(http.MethodGet).Handler(`/foo/:name`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    name := goblin.GetParam(r.Context(), "name")
    fmt.Fprintf(w, "/foo/%v", name)
}))

http.ListenAndServe(":9999", r)
```

## 正規表現を使ったルーティング
名前付きパラメータに正規表現を使うこと(`:paramName[pattern]`)で正規表現を使ったルーティングを定義することができます。

```go
r.Methods(http.MethodGet).Handler(`/foo/:id[^\d+$]`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    id := goblin.GetParam(r.Context(), "id")
    fmt.Fprintf(w, "/foo/%v", id)
}))
```

## ミドルウェア
リクエストの前処理、レスポンスの後処理に役立つミドルウェアをサポートしています。

任意のルーティングに対してミドルウェアを定義することができます。

グローバルにミドルウェアを設定することもできます。グローバルにミドルウェアを設定すると、すべてのルーティングにミドルウェアが適用されるようになります。

ミドルウェアは1つ以上設定することができます。

ミドルウェアはhttp.Handlerを返す関数として定義する必要があります。

```go
// http.Handlerを返す関数としてミドルウェアを実装
func global(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "global: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "global: after\n")
	})
}

func first(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "first: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "first: after\n")
	})
}

func second(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "second: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "second: after\n")
	})
}

func third(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "third: before\n")
		next.ServeHTTP(w, r)
		fmt.Fprintf(w, "third: after\n")
	})
}

r := goblin.NewRouter()

// グローバルにミドルウェアを設定
r.UseGlobal(global)
r.Methods(http.MethodGet).Handler(`/globalmiddleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "/globalmiddleware\n")
}))

// Useメソッドを使用することでミドルウェアを適用できます
r.Methods(http.MethodGet).Use(first).Handler(`/middleware`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middleware\n")
}))

// ミドルウェアは複数設定することができます
r.Methods(http.MethodGet).Use(second, third).Handler(`/middlewares`, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "middlewares\n")
}))

http.ListenAndServe(":9999", r)
```

`/globalmiddleware`にリクエストすると、次のような結果が得られます。

```
global: before
/globalmiddleware
global: after
```

`/middleware`にリクエストすると、次のような結果が得られます。

```
global: before
first: before
middleware
first: after
global: after
```

`/middlewares`にリクエストすると、次のような結果が得られます。

```
global: before
second: before
third: before
middlewares
third: after
second: after
global: after
```

## カスタム可能なエラーハンドラー
独自のエラーハンドラーを定義することができます。

定義可能なエラーハンドラは以下の2種類です。

- NotFoundHandler
  - ルーティングにマッチする結果が得られなかったときに実行されるハンドラです
- MethodNotAllowedHandler
  - マッチするメソッドがなかった場合に実行されるハンドラです

```go
func customMethodNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "customMethodNotFound")
	})
}

func customMethodAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "customMethodNotAllowed")
	})
}

r := goblin.NewRouter()
r.NotFoundHandler = customMethodNotFound()
r.MethodNotAllowedHandler = customMethodAllowed()

http.ListenAndServe(":9999", r)
```

## デフォルトOPTIONSハンドラー
OPTIONSメソッドでのリクエストの際に実行されるデフォルトのハンドラを定義することができます。

```go
func DefaultOPTIONSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNoContent)
	})
}

r := goblin.NewRouter()
r.DefaultOPTIONSHandler = DefaultOPTIONSHandler()

http.ListenAndServe(":9999", r)
```

デフォルトOPTIONSハンドラーは例えば、CORSのOPTIONSリクエスト（preflight request）の対応などに役立ちます。

# ベンチマークテスト
goblinのベンチマークテストを実行するコマンドを用意しています。

[Makefile](https://github.com/bmf-san/goblin/blob/master/Makefile)をご参照ください。

他のHTTP Routerとのベンチマーク比較結果が気になりますか？

こちらをご覧ください！
[bmf-san/go-router-benchmark](https://github.com/bmf-san/go-router-benchmark)

# 設計
goblinの内部的なデータ構造について解説します。

パフォーマンスが最適化されたHTTP Routerにおいては、[基数木](https://ja.wikipedia.org/wiki/%E5%9F%BA%E6%95%B0%E6%9C%A8)が採用されていることが多いですが、goblinは[トライ木](https://ja.wikipedia.org/wiki/%E3%83%88%E3%83%A9%E3%82%A4_(%E3%83%87%E3%83%BC%E3%82%BF%E6%A7%8B%E9%80%A0))をベースとしたデータ構造を採用しています。

基数木と比較すると、トライ木はメモリ使用量に劣る為、パフォーマンス面では不利です。しかしアルゴリズムの単純さ、理解しやすさは圧倒的にトライ木に軍配が上がるでしょう。

HTTP Routerは一見単純な仕様を持つアプリケーションに思えるかもしれませんが、意外と複雑です。これはテストケースを見て頂ければわかるかと思います。
（もっと良い感じのテストケースの実装アイデアがあればぜひ教えてください。）

単純なアルゴリズムを採用していることのメリットとしては、コードのメンテナビリティに貢献するという点です。（基数木の実装の難しさに対する言い訳とも聞こえるかもしれません・・実際のところ基数木をベースにしたHTTP Routerの実装の難しさには一度挫折しました・・）

[_examples](https://github.com/bmf-san/goblin/blob/master/_examples)のソースコードを例に、goblinの内部的なデータ構造について説明します。

ルーティングの定義を表で表すと、次のようになります。

| Method | Path | Handler | Middleware |
| -- | -- | -- | -- |
| GET | / | RootHandler | N/A |
| GET | /foo | FooHandler | CORS |
| POST | /foo | FooHandler | CORS |
| GET | /foo/bar | FooBarHandler | N/A |
| GET | /foo/bar/:name | FooBarNameHandler | N/A |
| POST | /foo/:name | FooNameHandler | N/A|
| GET | /baz | BazHandler | CORS |

gobinではこのようなルーティングは次のような木構造として表現されます。

```
凡例：<HTTP Method>,[Node]

<GET>
    ├── [/]
    |
    ├── [/foo]
    |        |
    |        └── [/bar]
    |                 |
    |                 └── [/:name]
    |
    └── [/baz]

<POST>
    └── [/foo]
             |
             └── [/:name]
```

HTTPメソッドごとに木を構築するようになっています。

各ノードはハンドラーやミドルウェアの定義をデータとして持っています。

ここでは説明を簡素にするため、名前付きルーティングのデータや、グローバルミドルウェアのデータなどを省略しています。

内部で構築される木には他にも色々なデータが保持されます。

詳しく知りたい場合はデバッカーを使って内部構造を覗いてみてください。

改善のアイデアがあればぜひ教えてください！

# コントリビューション
IssueやPull Requestはいつでもお待ちしています。

気軽にコントリビュートしてもらえると嬉しいです。

コントリビュートする際は、以下の資料を事前にご確認ください。

[CODE_OF_CONDUCT](https://github.com/bmf-san/goblin/blob/master/.github/CODE_OF_CONDUCT.md)
[CONTRIBUTING](https://github.com/bmf-san/goblin/blob/master/.github/CONTRIBUTING.md)

# スポンサー
もし気に入って頂けたのならスポンサーしてもらえると嬉しいです！
[GitHub Sponsors - bmf-san](https://github.com/sponsors/bmf-san)

あるいはstarを貰えると嬉しいです！

継続的にメンテナンスしていく上でのモチベーションになります :D

# ライセンス
MITライセンスに基づいています。

[LICENSE](https://github.com/bmf-san/goblin/blob/master/LICENSE)

## 作者
[bmf-san](https://github.com/bmf-san)

- Email
  - bmf.infomation@gmail.com
- Blog
  - [bmf-tech.com](http://bmf-tech.com)
- Twitter
  - [bmf-san](https://twitter.com/bmf-san)


