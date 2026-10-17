# サードパーティに関する通知

このプロジェクトには、サードパーティ製のソフトウェアおよび素材が含まれています。プロジェクト全体に適用されるMIT Licenseは、それらの素材に適用されるライセンスに代わるものではありません。リポジトリに取り込まれたファイルに含まれる著作権表示、ライセンス表示、および帰属表示は保持する必要があります。

## フロントエンドソフトウェア

| コンポーネント | バージョン | ライセンス | 著作権／アップストリーム | 配布方法 |
| --- | --- | --- | --- | --- |
| Alpine.js | 3.17.4 | MIT | Copyright © 2019–2025 Caleb Porzio and contributors; <https://github.com/alpinejs/alpine/blob/v3.17.4/LICENSE.md> | バージョン固定されたjsDelivr CDN URLからSRI付きで配信 |
| Bootstrap | 5.3.1 | MIT | Copyright 2011–2023 The Bootstrap Authors; <https://github.com/twbs/bootstrap/blob/v5.3.1/LICENSE> | バージョン固定されたjsDelivr CDN URLからSRI付きで配信 |
| Cash | 8.1.5 | MIT | Copyright (c) 2014–2020 Ken Wheeler, 2020–present Fabio Spampinato; <https://github.com/fabiospampinato/cash/blob/8.1.5/license> | `backend/web/js/external/cash.min.js` |
| jquery.qrcode | 1.0.0 | MIT | Copyright (c) 2011 Jerome Etienne; <https://github.com/jeromeetienne/jquery-qrcode/blob/master/MIT-LICENSE.txt> | `backend/web/js/external/qrcode.min.js`（Cashを使用するよう変更） |

フロントエンドに関する追加の簡略な通知は、`backend/web/js/external/LICENSES.md`に保持されています。圧縮済みファイルに含まれるライセンスおよび著作権のヘッダーは、必要な通知の一部であり、削除してはなりません。

BootstrapとAlpine.jsは実行時に取得され、公開されるツリーには取り込まれていません。バージョン固定されたjsDelivrのファイルとSRI値は、2026-10-10にダウンロードした正確なバイト列と照合済みです。Git管理対象外のローカル開発用コピーは、アプリケーションから参照されません。

## Font Awesome Free 5.15.4

Font Awesome Freeの実行時CSSとWOFF2 Webフォントは、`backend/web/fontawesome/`以下に取り込まれています。開発用ソース、メタデータ、SVGコレクション、JavaScriptビルド、および旧形式のフォントは配布されません。Font Awesomeでは、ファイルの種類ごとに異なるライセンスが適用されます。

- SVGおよびJavaScript形式のアイコン：Creative Commons Attribution 4.0 International（CC BY 4.0）
- Webフォントおよびデスクトップフォント：SIL Open Font License 1.1
- コード、およびフォントやアイコン以外のファイル：MIT License

Copyright Fonticons, Inc. Font Awesomeのブランドアイコンは各所有者の商標であり、対応する企業、製品、またはサービスを表す目的でのみ使用する必要があります。

同梱されているライセンスと帰属表示の要件は、`backend/web/fontawesome/LICENSE.txt`に記録されています。<https://fontawesome.com/license/free>も参照してください。

## Goモジュール

`backend/go.mod`、`backend/go.sum`、`backend/cmd/db/go.mod`、および`backend/cmd/db/go.sum`で宣言されているGoモジュールは、それぞれのアップストリームプロジェクトからダウンロードされ、このリポジトリにはソースコードとして取り込まれていません。

コンパイル済みバイナリまたはコンテナイメージを配布する前に、そのビルドで実際に使用する依存関係についてライセンスレポートを再生成し、リンクされるモジュールが要求するすべての通知を含めてください。依存関係のライセンスはバージョン間で変更される可能性があります。依存関係一覧の正式な情報源は、この文書ではなくGoモジュールファイルです。

## 削除した素材とプロジェクトで作成した代替素材

監査により、出所または再配布条件を十分に確認できない素材が見つかりました。それらは2026-10-14に配布対象のツリーから削除されました。古い非公開コピーから誤ってファイルを復元しないよう、以下の調査結果を履歴として残します。

代替ファイルの`backend/web/image/app-icon.svg`、`backend/web/image/allergy.svg`、および`backend/web/image/product.svg`は、2026-10-14にこのリポジトリ向けに作成されたプロジェクト独自の素材であり、プロジェクト全体のMIT Licenseが適用されます。

### face-api.jsのモデル重み

以前の`backend/web/js/external/weights/`ディレクトリには、`backend/web/js/verify_twenty.js`で使用される5種類のモデル群が含まれていました。公開されているすべてのマニフェストとシャードについて、2026-10-10にアップストリームのファイルとバイト単位で比較しました。

| モデル | 検証元 | 結果 | 再配布状況 |
| --- | --- | --- | --- |
| SSD MobileNet V1 | `justadudewhohacks/face-api.js`のcommit `a86f011d72124e5fb93e59d5c4ab98f699dd5c9c` | 3ファイルすべてが一致 | リポジトリはMIT Licenseですが、READMEには重みの取得元が`yeephycho/tensorflow-face-detection`であり、WIDER FACEを用いて学習したと記載されています。そのアップストリームプロジェクトはコードとモデルにApache 2.0を適用する一方、データセットの権利についてはWIDER FACEを参照しています。両方の通知を保持し、予定する配布と利用がデータセットの条件に適合することを確認する必要があります。 |
| 68点顔ランドマーク | `justadudewhohacks/face-api.js`のcommit `a86f011d72124e5fb93e59d5c4ab98f699dd5c9c` | 2ファイルとも一致 | リポジトリはMIT Licenseで、モデルファイルを公開素材へコピーするよう利用者に案内しています。READMEには約35,000枚のラベル付き顔画像で学習したと記載されていますが、データセットやそのライセンスは特定されていません。そのため、データセットに由来する権利は未確認のままです。 |
| 顔認識 | `justadudewhohacks/face-api.js`のcommit `a86f011d72124e5fb93e59d5c4ab98f699dd5c9c` | 3ファイルすべてが一致 | リポジトリはMIT Licenseですが、READMEでは学習者をDavis Kingとし、dlibの顔認識モデルと同等であると記載しています。ローカルリポジトリには、これらの重みの再配布権を付与する個別の通知がありません。 |
| Tiny Face Detector | `justadudewhohacks/face-api.js`のcommit `a86f011d72124e5fb93e59d5c4ab98f699dd5c9c` | 2ファイルとも一致 | リポジトリはMIT Licenseで、モデルファイルを公開素材へコピーするよう利用者に案内しています。READMEには約14,000枚の画像からなる独自データセットについて記載されていますが、その出所やライセンスは特定されていません。そのため、データセットに由来する権利は未確認のままです。 |
| Tiny Yolo v2 | `justadudewhohacks/face-api.js-models`のcommit `8401d0fde3e60fef4dd374cc3d36c5bcc4c1103a` | 5ファイルすべてが一致 | 配布元リポジトリにはライセンスファイルも、検出可能なライセンスのメタデータもありません。再配布の許可は確認できませんでした。 |

face-api.jsのREADMEは、モデルのマニフェストとシャードを公開素材のディレクトリへコピーするよう利用者に案内しており、face-api.jsのメインリポジトリにはMIT Licenseが適用されています。これは、そのリポジトリでホストされるファイルの再配布を裏付けるものですが、アップストリームモデルに不足している通知や、学習データセットに関連する権利および利用制限の問題を解決するものではありません。特に、公開GitHubリポジトリであること自体はライセンスになりません。

これらのモデルファイルは、**この監査では公開再配布を承認していません**。そのため、すべてのマニフェストと重みシャードを削除しました。ブラウザーはこれらのファイルおよび`face-api.js`を読み込まなくなりました。写真撮影と既存のサーバー側検証フローは、従来のクライアント側の顔数確認の補助表示がなくても引き続き利用できます。

関連するアップストリームの情報源：

- <https://github.com/justadudewhohacks/face-api.js>
- <https://github.com/justadudewhohacks/face-api.js-models>
- <https://github.com/yeephycho/tensorflow-face-detection>
- <http://shuoyang1213.me/WIDERFACE/>

### 画像、音声、フォント、およびアイコン

メディア一覧と埋め込みメタデータは、2026-10-10に確認しました。

| 素材 | 出所／ライセンスの調査結果 | 再配布状況 |
| --- | --- | --- |
| `backend/web/sound/chime.mp3` | ID3メタデータから、OtoLogicおよび`Doorbell-Ding_Dong02-1`であることを確認しました。OtoLogicは無料効果音をCC BY 4.0で公開し、帰属表示を保持することを条件に再配布を許可しています。配布元：<https://otologic.jp/free/se/doorbell01.html>、利用規約：<https://otologic.jp/free/license> | 下記の帰属表示を付けることで承認 |
| `backend/web/fontawesome/css/all.min.css` | Font Awesome Free 5.15.4のコード、MIT License | 承認。埋め込みヘッダーと`backend/web/fontawesome/LICENSE.txt`を保持すること |
| `backend/web/fontawesome/webfonts/*.woff2` | Font Awesome Free 5.15.4のWebフォント、SIL OFL 1.1 | 承認。`backend/web/fontawesome/LICENSE.txt`を保持すること |
| Font Awesome Webフォントから描画されるアイコン | 配布形態はOFLライセンスのWebフォントです。SVGおよびJavaScript形式でアイコンを配布する場合は、代わりにCC BY 4.0が適用されます。ブランドアイコンには、引き続き各所有者の商標権が適用されます。 | 現在配布しているWebフォント形式について承認 |
| 以前の`backend/web/image/allergy/*.jpg` | 28ファイルすべてが`https://storage.googleapis.com/row-allergy/`の画像と一致しましたが、作成者、ライセンス、または利用条件は確認できませんでした。 | 削除し、プロジェクトで作成した`allergy.svg`へ置き換え。リモート読み込みも削除 |
| 以前の`backend/web/image/180.png`、`192.png`、`512.png`、および`favicon.ico` | レストラン記号のアートワークと一致しましたが、作成者、出所、作成記録、またはライセンスは確認できませんでした。 | 削除し、プロジェクトで作成した`app-icon.svg`へ置き換え |
| 以前の`backend/web/image/product/*.jpg` | 4枚のサンプル商品画像と`base_image.jpg`について、作成者、出所、作成記録、またはライセンスは確認できませんでした。 | 削除し、プロジェクトで作成した`product.svg`へ置き換え。リモート読み込みも削除 |

保持している効果音に必要な帰属表示：

> 効果音：OtoLogic（CC BY 4.0）— `Doorbell-Ding_Dong02-1`

この効果音のSHA-256は`717b85d29622098e92e0b74d9519d4ca8e000e28446f9e55dbb68cc827b7fa10`です。効果音を編集した場合は、この通知を保持し、素材を変更したことを明示してください。

この監査で未承認と判断された画像は、配布対象のツリーには残っていません。アプリケーションのスクリプトは、プロジェクトで作成したローカルの代替SVGのみを使用します。

## ライセンス本文

上記で特定したコンポーネントには、次のアップストリームのライセンス本文が適用されます。

- MIT License：<https://opensource.org/license/mit>
- Apache License 2.0：<https://www.apache.org/licenses/LICENSE-2.0>
- SIL Open Font License 1.1：<https://openfontlicense.org/open-font-license-official-text/>
- Creative Commons Attribution 4.0 International：<https://creativecommons.org/licenses/by/4.0/legalcode>

このプロジェクトを再配布する場合は、この通知、プロジェクトの`LICENSE`、取り込まれているFont Awesomeのライセンス、ソースファイルに埋め込まれた通知、および配布形態に応じて必要となる追加のライセンス本文を保持してください。
