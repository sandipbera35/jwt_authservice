# jwt_authservice
This Is an Implementation Of JWT User Authentication Service In Golang Fiber Framework with postgresql

# To Run Minio Object Storage Server Run Code Below 
Docker Should be installed before running minio
docker run --name minio  --publish 9000:9000  --publish 9001:9001  -e "MINIO_ROOT_USER=YOURUSERNAME" -e "MINIO_ROOT_PASSWORD=YOURPASSWORD" --volume d:/YourFolderPath: /data bitnami/minio:latest

OR 

Visit https://min.io/docs/minio/container/index.html

# Install Postgresql 

visit https://www.postgresql.org/download/

# Install Postman

visit https://www.postman.com/downloads/

Import The Postman Collection File From postman folder of the project 

