# miguser

## Overview

miguserは、ニフクラのサービスであるリモートアクセスVPNゲートウェイ(RemoteAccessVpnGateway / RAVGW)のユーザー情報の移行を支援 するツールです  
miguserを使うことで既存のRAVGWからRAVGWへ簡単にユーザー情報を一括移行することができます  
RAVGWv1をRAVGWv2に置き換える際などにご利用いただけます

## Install

### バイナリ

https://github.com/nifcloud/nifcloud-ravgw-miguser/releases
※ソースコードは今後公開予定となります

## Usage

`miguser`コマンドは、`miguser export`と`miguser import`の2つのサブコマンドで構成されます

### help

`miguser`、`miguser export`、`miguser import`は、それぞれ`-h / --help`でヘルプを参照可能です

### miguser export

移行元となるRAVGWのユーザー情報を、CSVとして出力します  
※移行元のRAVGWは事前に作成しておく必要があります

#### 事前に必要な情報

- 移行元RAVGWが配置されているリージョン
  - 確認方法: https://pfs.nifcloud.com/api/endpoint.htm
    - 例) リージョン: `east-1`、エンドポイント: `https://jp-east-1.computing.api.nifcloud.com/api/`の場合
      - `east-1`ではなく`jp-east-1`を指定してください
- 移行元ユーザーのAccessKey, SecretAccessKey
  - 確認方法: https://pfs.nifcloud.com/help/acc/key.htm
- 移行元RAVGWのRemoteAccessVpnGatewayID
  - 確認方法: https://pfs.nifcloud.com/help/ra_vpngw/detail.htm

ユーザー情報のエクスポート

```
./miguser export \
    --region "YOUR_REGION_WHERE_YOUR_RAVGW_IS_LOCATED" \
    --access-key "YOUR_ACCESS_KEY_ID" \
    --secret-access-key "YOUR_SECRET_ACCESS_KEY" \
    --ravgwid "YOUR_RAVGW_ID"
```

ショートハンド

```
./miguser export \
    -r "YOUR_REGION_WHERE_YOUR_RAVGW_IS_LOCATED" \
    -a "YOUR_ACCESS_KEY_ID" \
    -s "YOUR_SECRET_ACCESS_KEY" \
    --ravgwid "YOUR_RAVGW_ID"
```

コマンドを実行したディレクトリ内に、`<"YOUR_RAVGW_ID">.csv`という名前でCSVファイルが出力されます

### エクスポートされたユーザー情報のCSVへパスワードを追加

`miguser export`で出力されたCSVは、`UserName`、`Password`、`Description`の3列で構成されています  
CSVファイルを開くと、以下のようにPasswordの列が空欄になっています  
各ユーザーのパスワードを追記してください

例)テキストとしてCSVファイルを開いた場合  
追記前

```
UserName,Password,Description
user-1,,user-1-description
user-2,,user-2-description
```

追記後

```
UserName,Password,Description
user-1,user-1-password,user-1-description
user-2,user-2-password,user-2-description
```

### miguser import

ユーザー情報のCSVファイルを移行先のRAVGWへ追加します  
※移行先のRAVGWは事前に作成しておく必要があります

#### 事前に必要な情報

- 移行先RAVGWが配置されているリージョン
  - 確認方法: https://pfs.nifcloud.com/api/endpoint.htm
    - 例) リージョン: `east-1`、エンドポイント: `https://jp-east-1.computing.api.nifcloud.com/api/`の場合
　    - `east-1`ではなく`jp-east-1`を指定してください
- 移行先ユーザーのAccessKey, SecretAccessKey
  - 確認方法: https://pfs.nifcloud.com/help/acc/key.htm
- 移行先RAVGWのRemoteAccessVpnGatewayID
  - 確認方法: https://pfs.nifcloud.com/help/ra_vpngw/detail.htm
- CSVファイルのパス
  - Password列を埋めたCSVファイルのパスを指定します

ユーザー情報のインポート

```
./miguser import \
    --region "YOUR_REGION_WHERE_YOUR_RAVGW_IS_LOCATED" \
    --access-key "YOUR_ACCESS_KEY_ID" \
    --secret-access-key "YOUR_SECRET_ACCESS_KEY" \
    --ravgwid "YOUR_RAVGW_ID" \
    --src "YOUR_CSV_FILE"
```

ショートハンド

```
./miguser export \
    -r "YOUR_REGION_WHERE_YOUR_RAVGW_IS_LOCATED" \
    -a "YOUR_ACCESS_KEY_ID" \
    -s "YOUR_SECRET_ACCESS_KEY" \
    --ravgwid "YOUR_RAVGW_ID" \
    --src "YOUR_CSV_FILE"
```
#### お問い合わせ

ツールに関するお問い合わせは[ベーシックサポート（トラブル窓口）](https://pfs.nifcloud.com/inquiry/support.htm)のサポート範囲 外となります  
ツールに関するお問い合わせは[Issue](https://github.com/nifcloud/nifcloud-ravgw-miguser/issues)を起票してください  
コミュニティベースでのサポートとなります  
