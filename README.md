# go-recipes-api

TL;DR - A RESTful api project built using Gin/Go following the patterns and book `Building distributed application in Gin` by `Mohamed Labouardy` published in 2021 and found on github [here](https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin)

## Project Configuration

Project is using the GitFlow approach. You can read more about that [here](https://nvie.com/posts/a-successful-git-branching-model/)

### Onboarding

If joining the development team on this project you should clone the repository to you local setup and run the `go mod download` command to install required dependencies after completing the cloning.

### Swagger/OpenAPI
- generate new document with: `swagger generate spec -o ./swagger.json`
- serve document with: `swagger serve ./swagger.json`

### Setup for MongoDB
Project uses a dockerised instance of MongoDB 4.4.3 and compass to access locally. The module mjst have access to the mongoDB driver and therefore you must run the `go get go.mongodb.org/mongo-driver/mongo` command to add the dependecy to our project. This command will modify the `go.mod` and `go.sum` files.

With a docker image pulled to the host system yo can run it with the following:
`docker run -d --name mongodb –v /home/data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3`

You can run the program once the mongoDB docker container is running with the `MONGO_URI="mongodb://USER:PASSWORD@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run main.go` command.

It is also possible to load the data without using Go code. Try the command: `mongoimport --username admin --password password --authenticationDatabase admin --db demo --collection recipes --file recipes.json --jsonArray`

To do it within Go code we can use the following init() function.
```go
func init() {
	recipes = make([]Recipe, 0)

	// temp data seeding
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)

	// mongoDB client connection
	ctx := context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully Connected to MongoDB")

	// data insertion
	// var listOfRecipes []interface{}
	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }

	// collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("inserted recipes:", len(insertManyResult.InsertedIDs))
}
```