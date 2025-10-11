# Server code

## Run the server image with Docker
1. Do install docker and docker compose
2. Create the directory for the server and where the sql init will be
```shell
mkdir server
mkdir server/db/init
```
3. Create the .env file with your values and the docker-compose.yml, that is the same as in the project
```shell
cd server
vim .env
vim docker-compose.yml
```
4. Create the init.sql 
```shell
cd db/init
vim init.sql
```
and input the content of the ./db/init/01-init.sql
5. Lastly run the docker
```shell
docker compose up
```