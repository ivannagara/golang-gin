how to setup MongoDB =>

1) Run the mongoDB image with the following docker command :
{
    docker run -d --name mongodb –v /home/data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo
}
the username will be 'admin' and the password will be 'password',
then you should input this username and password in the MongoDB Compass
application with the login with Authentication; then type in the username
and password. 


2) Then, you should install the mongoDB dependencies in your GoLang code like below: 
[a]
{
    go get go.mongodb.org/mongo-driver/mongo
}

[b]
{
    module github.com/mlabouardy/recipes-api
    go 1.15
    require (
        github.com/gin-gonic/gin v1.6.3
        github.com/rs/xid v1.2.1
        go.mongodb.org/mongo-driver v1.4.5
    )
}

[c] -- main.go file --
{
    package main
    import (
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "go.mongodb.org/mongo-driver/mongo/readpref"
    )
}

3) After that, copy the URI String from MongoDB and include it in your Environment Variable, and run it with your GoLang code in the terminal by typing =>
{
    MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" go run main.go
}
the MONGO_URI is defined in the code already, and is read by using :
{
    os.Getenv("MONGO_URI")
}

4) If you want to specify the database name, you can do this =>
{
    MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run main.go
}
the database name will be named 'demo'.
The MONGO_DATABASE is also specified in the GoLang code of
{
    os.Getenv("MONGO_DATABASE")
}

5) You should also initiate the REDIS image as a running container
because the code now also uses the REDIS container as well