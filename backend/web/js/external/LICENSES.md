# フロントエンドライブラリのライセンス

次のブラウザー向けライブラリは、このディレクトリに取り込まれ、アプリケーションから読み込まれます。ソースファイル、バージョン、ライセンス、およびローカルでの変更内容は、リンク先のアップストリームリポジトリと照合済みです。

## Cash 8.1.5

- ファイル：`cash.min.js`
- ライセンス：MIT
- Copyright (c) 2014-2020 Ken Wheeler, 2020-present Fabio Spampinato
- アップストリーム：<https://github.com/fabiospampinato/cash/tree/8.1.5>
- SHA-256: `9a044188efdb625c5e04d1220698c099927ff16bfb434c37cd7f04dd5ee1ae1f`
- ローカルでの変更：なし。ファイルはアップストリームの8.1.5配布物と一致します。

## jquery.qrcode 1.0.0

- ファイル：`qrcode.min.js`
- ライセンス：MIT
- Copyright (c) 2011 Jerome Etienne, <http://jetienne.com>
- アップストリーム：<https://github.com/jeromeetienne/jquery-qrcode>
- SHA-256: `41ea01fd8d9a05bb9c45e73b0b3b3e3f62dc1e44924abdff9fd39f9362f3379a`
- ローカルでの変更：取り込まれたCash APIを使用するため、アップストリームの圧縮済みファイル末尾にあるグローバル依存を`jQuery`から`cash`へ変更しています。

## MIT License本文（英語原文）

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Alpine.jsとBootstrapは、バージョン固定されたCDN URLからSRI検証付きで読み込まれます。詳細はリポジトリルートの`THIRD_PARTY_NOTICES.md`に記載しています。
