#!/bin/bash


# 対象のコマンドをインストールするような処理
if !(type "sqlite3" > /dev/null 2>&1); then
	sudo apt-get install -y sqlite3
fi

# if !(type "yq" > /dev/null 2>&1); then
# 	sudo snap install yq 
# fi


# STORE_MASTER ---------------------------------
# . ../../../.env
# STORE_ID=`cat ../../../internal/env/${ENVIRONMENT}.yaml | yq '.STORE_ID'`
echo "--1 start ----"
STORE_MASTER="master.sqlite3"
STORE_TARGET=${STORE_ID}.sqlite3

echo "--2 start ----"
sqlite3 ${STORE_MASTER} < drop.sql
sqlite3 ${STORE_MASTER} < create.sql
sqlite3 ${STORE_TARGET} < drop.sql
sqlite3 ${STORE_TARGET} < create.sql

echo "--product.tsv start ----"
cat product.tsv        | sqlite3 ${STORE_MASTER} ".mode tabs" ".import /dev/stdin product"

echo "--tag_primary.tsv start ----"
cat tag_primary.tsv    | sqlite3 ${STORE_MASTER} ".mode tabs" ".import /dev/stdin tag_primary"

echo "--tag_secondary.tsv start ----"
cat tag_secondary.tsv  | sqlite3 ${STORE_MASTER} ".mode tabs" ".import /dev/stdin tag_secondary"

