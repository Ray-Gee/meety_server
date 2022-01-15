1. pgadmin設定
1. client 一覧
1. client create画面

./init.sh

client >> react
server >> go

brew services start postgresql
psql postgres 

docker compose up -d
cd server
go run main.go
<!-- オートリロード -->
gin -p 3001 -i run main.go

cd client
yarn start

psql -p 5432 -h localhost -d book_keeper -U postgres 
<!-- 環境変数でパスワード対応済み
password = password -->