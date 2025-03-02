# p-blog-api

curl -X POST http://localhost:8080/api/categories -H "Content-Type: application/json" -d '{
  "name": "Tecnologia"
}'

curl -X GET http://localhost:8080/api/categories


curl -X POST http://localhost:8080/api/categories -H "Content-Type: application/json" -d '{
  "name": "Tecnologia"
}'

curl -X POST http://localhost:8080/api/posts -H "Content-Type: application/json" -d '{
  "title": "Meu primeiro post",
  "category_id": "67c3b2fe36728a583c1fefb8",
  "content": "Este é o conteúdo do post em Markdown."
}'

curl -X GET http://localhost:8080/api/posts