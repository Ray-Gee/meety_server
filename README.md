1. meかmeじゃないかを上の階層から渡したい/Users/admin/works/meety/client/src/containers/Messages/hooks/useMessages.jsらへん見てた
1. docker server からdb繋げないpersonControllers.go L:71dial tcp 127.0.0.1:5432: connect: connection refused

./init.shと./generate_code.shはdockerに置き換わった

client repository https://github.com/Ryuichi-g/meety_client
server repository https://github.com/Ryuichi-g/meety_server
client >> react
server >> go

<!-- db -->
brew services start postgresql
psql postgres 
<!-- server -->
docker compose up -d
cd server
go run main.go
gin -p 3001 -i run main.go <!-- ホットリロード -->
<!-- client -->
cd client
yarn start

psql -p 5432 -h localhost -d book_keeper -U postgres 
<!-- 環境変数でパスワード対応済み
password = password -->