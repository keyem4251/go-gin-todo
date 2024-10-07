# ベースイメージを指定（Go 1.21以上）
FROM golang:1.23.2-alpine

# 作業ディレクトリを設定
WORKDIR /workspace

# Goモジュールファイルをコピーして依存関係をダウンロード
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# アプリケーションをビルド
RUN go build -o /todo-app

# アプリケーションを実行
CMD ["/todo-app"]