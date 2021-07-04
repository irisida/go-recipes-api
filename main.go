// Recipes API
//
// This is a sample API covering fairly simple endpoints in a
// recipes domain.
//
// Schemes: http
// Host: localhost:8080
// BASEPATH: /
// Version: 1.0.0
// Contact: Ed Leonard <edward.leonard@gmail.com>
//
// consumes:
// -application/json
//
// Produces:
// -application/json
// swagger:meta
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)

	// temp data seeding
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//   200:
//   description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation POST /recipes recipes NewRecipe
// Creates a new recipe
// ---
// requestBody:
// - name: string
//   tags: [string]
//   ingredients: [string]
//   instructions: [string]
//   required: true
// produces:
// - application/json
// responses:
//   '200':
//   description: Successful operation
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /recipes/{id} recipes UpdateRecipe
// Update an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe to be updated
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
//   description: Successful operation
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// default the index position and conduct a search for the match
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	// update the record and also the record
	// timestamp to reflect it was changed.
	recipe.PublishedAt = time.Now()
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes DeleteRecipe
// Deletes an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe to be deleted
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
//   description: Successful operation
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe hs been deleted",
	})
}

// swagger:operation GET /recipes/search?tag={tag} recipes SearchRecipes
// Searches existing recipes and returns matches basedon the
// tag that is passed as the search criteria
// ---
// parameters:
// - name: tag
//   in: path
//   description: string to use as a search criteria against a recies tags
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '200':
//   description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}

		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {
	router := gin.Default()

	// routes
	router.GET("/recipes", ListRecipesHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.Run()
}
