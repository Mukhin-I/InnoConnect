Hello *hello!* **HELLO!!!!**

Here is ```README.md``` and I will show you how to run this project using docker!

It's simple: 
* create ```.env``` file or rename ```.env.example``` -> ```.env```
* then open here a ```terminal```
* run ```docker compose --env-file .env -f backend/deployment/docker/docker-compose.yml up -d --build```
* And it starts!!!
* you can verify it by writing ```docker ps -a``` in the terminal and see (healthy) mark at the corresponding container names
* or **STOP** it using ```docker compose --env-file .env -f backend/deployment/docker/docker-compose.yml down```

*This is IT*! 
Good **luck** in your future work!

Best withes, MJ and AF