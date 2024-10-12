# ベースイメージを指定（Go 1.21以上）
FROM golang:1.23.2-alpine

RUN go install github.com/air-verse/air@latest

RUN apk add --no-cache git openssh curl

# 作業ディレクトリを設定
WORKDIR /workspace

# Goモジュールファイルをコピーして依存関係をダウンロード
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# アプリケーションを実行
CMD ["air"]