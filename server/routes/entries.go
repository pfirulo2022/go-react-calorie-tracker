package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pfirulo2022/go-react-calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var entryCollection *mongo.Collection = OpenCollection(Client, "calories")
var validate = validator.New()

// Agrego una entrada
func AddEntry(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	entry.ID = primitive.NewObjectID()
	result, insertErr := entryCollection.InsertOne(ctx, entry)
	if insertErr != nil {
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	c.JSON(http.StatusOK, result)

}

// Obtengo todas las entradas
func GetEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
	}
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
	}
	defer cancel()
	c.JSON(http.StatusOK, entries)

}

// Entradas por ingredientes
func GetEntriesByIngredient(c *gin.Context) {
	ingredient := c.Param("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, entries)

}

/* func GetEntriesByIngredient(c *gin.Context) {
	ingredient := c.Param("id")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient})

	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, entries)

} */

// Obtengo entrada por id
func GetEntryById(c *gin.Context) {
	EntryID := c.Params.ByName(("id"))
	docId, err := primitive.ObjectIDFromHex(EntryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

// func GetEntryById(c *gin.Context) {
// 	EntryID := c.Params.ByName(("id"))
// 	docId, _ := primitive.ObjectIDFromHex(EntryID)

// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var entry bson.M
// 	if err := entryCollection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)

// 	}

// 	defer cancel()
// 	fmt.Println(entry)
// 	c.JSON(http.StatusOK, entry)

// }

// Actualizando  un registro en la base de datos
func UpdateEntry(c *gin.Context) {

	entryID := c.Params.ByName(("id"))
	docID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid object ID"})
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entry models.Entry

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	result, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"dish":        entry.Dish,
			"fat":         entry.Fat,
			"ingredients": entry.Ingredients,
			"calories":    entry.Calories,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)
}

// Actualizo ingredientes
func UpdateIngredient(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, err := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Entry ID."})
		return
	}

	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	var ingredient Ingredient

	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, result.ModifiedCount)

}

// func UpdateIngredient(c *gin.Context) {
// 	entryID := c.Params.ByName("id")
// 	docID, err := primitive.ObjectIDFromHex(entryID)
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Entry ID."})
// 		return
// 	}

// 	type Ingredient struct {
// 		Ingredients *string `json:"ingredients"`
// 	}
// 	var ingredient Ingredient

// 	if err := c.BindJSON(&ingredient); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}
// 	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
// 		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, result.ModifiedCount)

// }

// Elimino una entrada

func DeleteEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, err := primitive.ObjectIDFromHex(entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.DeletedCount)
}

// func DeleteEntry(c *gin.Context) {
// 	entryID := c.Params.ByName("id")
// 	docID, _ := primitive.ObjectIDFromHex(entryID)

// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 	}

// 	defer cancel()
// 	c.JSON(http.StatusOK, result.DeletedCount)

// }
