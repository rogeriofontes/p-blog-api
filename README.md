# p-blog-api

curl -X POST http://localhost:8080/api/categories -H "Content-Type: application/json" -d '{
  "name": "Tecnologia"
}'

curl -X GET http://localhost:8080/api/categories

curl -X PUT http://localhost:8080/api/categories/65d4a6b4c79f4b245e7d1e0f -H "Content-Type: application/json" -d '{
  "name": "Categoria Atualizada"
}'

curl -X DELETE http://localhost:8080/api/categories/65d4a6b4c79f4b245e7d1e0f

curl -X GET http://localhost:8080/api/categories/65d4a6b4c79f4b245e7d1e0f

curl -X POST http://localhost:8080/api/categories -H "Content-Type: application/json" -d '{
  "name": "Tecnologia"
}'
----- POST

curl -X POST http://localhost:8080/api/posts -H "Content-Type: application/json" -d '{
  "title": "Meu primeiro post",
  "category_id": "67c4661d0a4028e7a38392a9",
  "content": "Este é o conteúdo do post em Markdown.",
  "tags": ["67c46876496710e2d5f1bee9","67c46942496710e2d5f1beea"]
}'

curl -X GET http://localhost:8080/api/posts -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDEwNTEyNDAsInVzZXJfaWQiOiI2N2M0NzM4MDU3ZDZjODNjNzdlMmE4MjgifQ.eFK7BpxTiAZrXGsvDe7SXPx5M_TvUIj1fudSN2mLRCQ"

curl -X PUT http://localhost:8080/api/posts/65d4b32f4c79f4b245e7d2f9a -H "Content-Type: application/json" -d '{
  "title": "Título atualizado",
  "category_id": "67c3b2fe36728a583c1fefb8",
  "content": "Conteúdo atualizado do post.",
  "tags": ["67c3b30036728a583c1fefc9", "67c3b30136728a583c1fefca"]
}'


curl -X DELETE http://localhost:8080/api/posts/65d4b32f4c79f4b245e7d2f9a

curl -X GET http://localhost:8080/api/posts/65d4b32f4c79f4b245e7d2f9a

-- TAGS -- 
curl -X POST http://localhost:8080/api/tags -H "Content-Type: application/json" -d '{
  "name": "Golang"
}'

curl -X POST http://localhost:8080/api/tags -H "Content-Type: application/json" -d '{
  "name": "Java"
}'

curl -X GET http://localhost:8080/api/tags

curl -X PUT http://localhost:8080/api/tags/65d4a6b4c79f4b245e7d1e0f -H "Content-Type: application/json" -d '{
  "name": "Nova Tag Atualizada"
}'

curl -X DELETE http://localhost:8080/api/tags/65d4a6b4c79f4b245e7d1e0f
curl -X GET http://localhost:8080/api/users/65d4b32f4c79f4b245e7d2f9a


---- Comentário ---
curl -X POST http://localhost:8080/api/comments -H "Content-Type: application/json" -d '{
  "post_id": "67c4695f496710e2d5f1beeb",
  "user_id": "12345",
  "content": "Ótimo post!"
}'

curl -X GET http://localhost:8080/api/comments/post/65d4a6b4c79f4b245e7d1e0f
curl -X GET http://localhost:8080/api/comments/67c4695f496710e2d5f1beeb

curl -X PUT http://localhost:8080/api/comments/65d4a6b4c79f4b245e7d1e0f -H "Content-Type: application/json" -d '{
  "content": "Comentário atualizado!"
}'

curl -X DELETE http://localhost:8080/api/comments/65d4a6b4c79f4b245e7d1e0f

curl -X GET http://localhost:8080/api/comments/65d4b32f4c79f4b245e7d2f9a

--- Like
curl -X POST http://localhost:8080/api/reactions -H "Content-Type: application/json" -d '{
  "post_id": "67c4695f496710e2d5f1beeb",
  "user_id": "12345",
  "reaction": true
}'


curl -X GET http://localhost:8080/api/reactions/likes/67c4695f496710e2d5f1beeb
curl -X GET http://localhost:8080/api/reactions/dislikes/67c4695f496710e2d5f1beeb
curl -X GET http://localhost:8080/api/reactions/65d4b32f4c79f4b245e7d2f9a

--- User
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{
  "username": "rogerio",
  "email": "rogerio@example.com",
  "password": "minhasenha123"
}'

curl -X GET http://localhost:8080/api/users/67c4738057d6c83c77e2a828
curl -X GET http://localhost:8080/api/users

api.PUT("/users/:id", controllers.UpdateUser)
api.DELETE("/users/:id", controllers.DeleteUser)


curl -X PUT http://localhost:8080/api/users/65d4b32f4c79f4b245e7d2f9a -H "Content-Type: application/json" -d '{
  "username": "rogerio_updated",
  "email": "rogerio_new@example.com",
  "password": "novaSenhaSegura"
}'

echo -n "minhasenha123" | sha256sum | awk '{print $1}'

no mongo:
db.createUser({
  username: "rogerio",
  "email": "rogerio@example.com"
  pwd: "4c94bd7240e61a20c60568f9aebe999e1d94b952ddc5dffeaa4db7257974a255"
})

curl -X DELETE http://localhost:8080/api/users/65d4b32f4c79f4b245e7d2f9a

curl -X GET http://localhost:8080/api/users/65d4b32f4c79f4b245e7d2f9a

====
rogerio_updated

curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{
  "email": "rogerio@example.com",
  "password": "minhasenha123"
}'

--------------------------------------------
--Fovoritos
curl -X POST http://localhost:8080/api/favorites -H "Content-Type: application/json" -d '{
  "user_id": "67c4738057d6c83c77e2a828",
  "post_id": "67c4695f496710e2d5f1beeb"
}'

curl -X DELETE http://localhost:8080/api/favorites/65d4b32f4c79f4b245e7d2f9a/67c3b2fe36728a583c1fefb8
//curl -X GET http://localhost:8080/api/favorites/67c4dc147f21dfdaed527983
curl -X GET http://localhost:8080/api/favorites/user/67c4738057d6c83c77e2a828

--- 

curl -X POST http://localhost:8080/api/follow -H "Content-Type: application/json" -d '{
  "user_id": "67c4738057d6c83c77e2a828",
  "follow_id": "67c4738057d6c83c77e2a828"
}'

curl -X GET http://localhost:8080/api/followers/user/67c4738057d6c83c77e2a828

=========================
#Melhorias

🔹 Soft Delete: Em vez de excluir, podemos apenas marcar o post como deletado (is_deleted: true).
🔹 Soft Update: Registrar um updated_at no post.

✅ Sistema de favoritos/seguidores -ok
✅ Busca avançada e filtros 
✅ Sistema de notificações
✅ Paginação para desempenho 
✅ Colocar ingra da ec2 

docker exec -it mongodb mongosh "mongodb://admin:admin@localhost:27017"
show dbs;
use blog;
show collections;
db.users.find().pretty();

docker build -t p-blog-api .
docker run -d -p 8080:8080 --name p-blog-api-container p-blog-api

docker build -t rogeriofontes/p-blog-api:v2 .
docker login
docker push rogeriofontes/p-blog-api:v2