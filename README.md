# sb2tt

## 概要

シフトボードの共有テキストを受け取ってTimeTreeに転載するLINE Botです

このREADMEの使い方の部分はかなり省略して書いているので、そのうち書き足します

## 注意

* セキュリティが不十分
  * INSTALLATION_IDのみで連携を行なう仕様となっている
  * INSTALLATION_IDは桁数が小さく数字のみであるため、総当たりで特定される可能性がある
  * 書き込みしかできないため、そこからの情報漏洩はないが、大量の書き込みを行う等の被害の可能性は0ではない

## 使い方

### 初期設定

1. LINE Developersでチャネル(Messaging API)を、TimeTree DevでCalendar Appを作成する
   * それぞれの作成方法に関しては公式ドキュメント等をご覧ください
2. 本リポジトリをダウンロードし、Herokuにホストする
   * Dockerを用いるのがおすすめ([参考](https://devcenter.heroku.com/ja/articles/build-docker-images-heroku-yml))
3. HerokuのアドオンでHeroku Postgresを追加
4. Herokuで環境変数を指定
   * CALENDAR_APP_ID: TimeTreeで作成したCalendar Appから取得
   * DATABASE_URL: アドオンを追加したタイミングで自動指定されている
   * GIN_MODE: 「release」と設定
   * LINE_ACCESS_TOKEN: LINE Developersで作成したチャネルから取得
   * LINE_ADMIN_ID: LINE Developersから取得
   * LINE_CHANNEL_SECRET: LINE Developersで作成したチャネルから取得
   * PRIVATE_KEY: TimeTreeで作成したCalendar Appから取得
5. LINE Developersで作成したチャネルのQRコード等からLINE Botを追加する
6. Calendar Appを使いたいカレンダーにインストールする
7. LINE Botに送られてきたINSTALLATION_IDを控えておく
   * 複数人で使用する場合に共有します
8. INSTALLATION_ID(数字部分)をそのまま送信し返す

### LINE Botの仕様

以下では、

* Botに送信するメッセージ
  * Botの応答
    * 詳細

のフォーマットで説明します

* シフトボード左上のアイコンから共有(転送)したテキスト
  * TimeTreeにシフトボードの予定をコピーする
    * 「バイト先を表示」をONにした場合はそのタイトルをそのままTimeTree側でも用いる
    * 「バイト先を表示」をOFFにした場合はデフォルトで設定したタイトルをTimeTree側で用いる
* 数値のみ
  * INSTALLATION_IDを設定する
* その他の文字列
  * シフトボードで「バイト先を表示」をOFFにして共有した場合の、TimeTreeの予定タイトルを送信した文字列に設定する
    * 数値のみの設定は不可
