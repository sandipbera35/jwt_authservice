# jwt_authservice
This Is an Implementation Of JWT User Authentication Service In Golang Fiber Framework with postgresql
# Install Docker 
visit https://docs.docker.com/engine/install/
# To Run Minio Object Storage Server Run Code Below In Terminal Or CommandPrompt
Docker Should be installed before running minio 

docker run --name minio  --publish 9000:9000  --publish 9001:9001  -e "MINIO_ROOT_USER=YOURUSERNAME" -e "MINIO_ROOT_PASSWORD=YOURPASSWORD" --volume d:/YourFolderPath: /data bitnami/minio:latest

OR 

Visit https://min.io/docs/minio/container/index.html

# Install Postgresql 

visit https://www.postgresql.org/download/

# Install Postman

visit https://www.postman.com/downloads/

Import The Postman Collection File From postman folder of the project 


# Some Screenshots hare ...


# Register New user
![image](https://github.com/user-attachments/assets/5c24d474-faf6-479c-830f-48ee2fc1d701)

# User PassWords saved in database as encripted form

![image](https://github.com/user-attachments/assets/ebfc796c-e456-4f6d-8249-ea5938d2019c)

# LogIn Or Get Token

![image](https://github.com/user-attachments/assets/a53501c1-a5c8-479c-8ee1-debff888eecf)


# Get Profile With JWT Token 

![image](https://github.com/user-attachments/assets/64fe8846-bb3f-48d2-9e59-fb397e4dd275)

# Add OR Update profile picture

![image](https://github.com/user-attachments/assets/215840b1-0573-498f-b90f-8a0d601c417d)

# Add OR Update Cover Picture

![image](https://github.com/user-attachments/assets/51d97b0b-8108-40fe-8496-0558d502c7ab)

# ADD SUPER USER

![image](https://github.com/user-attachments/assets/75a2edaa-7566-4af2-919e-afa7434517c8)








