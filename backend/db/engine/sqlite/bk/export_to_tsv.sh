#!/bin/bash

# 現在の日時を取得し、フォーマットする
current_time=$(date +"%Y%m%d-%H%M%S")
STORE_TARGET=${STORE_ID}.sqlite3

# エクスポートするファイル名を設定
output_file="bk/${current_time}_${STORE_ID}_product.tsv"

# SQLiteデータベースからデータをエクスポート
sqlite3 ${STORE_TARGET} <<EOF
.headers off
.mode tabs
.output $output_file
SELECT * FROM product;
EOF

echo "データが${output_file}にエクスポートされました。"
