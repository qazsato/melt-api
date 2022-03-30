# Melt API

[Melt](https://github.com/qazsato/melt) で利用するための REST API 。

API Gateway + Lambda (Serverless Framework) の構成で、ランタイムは Go。

## ドキュメント

[API 仕様書](https://qazsato.github.io/melt-api)

## ディレクトリ構成

```
.
├── config # 各環境の設定ファイル
│   ├── local.yml
│   ├── development.yml
│   └── production.yml
├── functions # Lambda関数
│   ├── images
│   │   ├── function.yml # serverless.ymlの分割ファイル
│   │   └── main.go # Lambda関数の実行ファイル
│   └── ...
├── go.mod
├── Makefile
├── package.json
└── serverless.yml
```

## セットアップ

```
make init
```

ローカル起動
```
make serve
```
