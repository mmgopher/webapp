curl http://localhost:8080/rest/user/5aa63014d0667219f82106d4
curl -X DELETE http://localhost:8080/rest/user/5aa654b38487b6260817a79b
curl -X POST -H "Content-Type: application/json" -d "{\"login\":\"bond\",\"password\":\"bond\",\"first\":\"James\",\"last\":\"Bond\",\"age\":27}" http://localhost:8080/rest/user
curl -X PUT -H "Content-Type: application/json" -d "{\"id\":\"5aab967c8487b622d0accf69\",\"login\":\"bond\",\"first\":\"James\",\"last\":\"Bond\",\"age\":57}" http://localhost:8080/rest/user