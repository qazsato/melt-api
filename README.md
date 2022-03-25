# Melt API

[Melt](https://github.com/qazsato/melt) で利用するための REST API 。

API Gateway + Lambda (Serverless Framework) の構成で、ランタイムは Go。

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
├── package.json
└── serverless.yml
```

## セットアップ

```
npm install
```

ローカル起動
```
npm run offline
```
