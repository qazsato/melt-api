# Melt API

[Melt](https://github.com/qazsato/melt) で利用するための REST API 。
API Gateway + Lambda (Serverless Framework) の構成で、ランタイムは Node.js。

## ディレクトリ構成

```
.
├── config # 各環境の設定ファイル
│   ├── local.yml
│   ├── development.yml
│   └── production.yml
├── functions # Lambda関数
│   ├── post_images
│   │   ├── function.yml # serverless.ymlの分割ファイル
│   │   └── handler.js # Lambda関数の実行ファイル
│   └── ...
├── package.json
├── package-lock.json
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
