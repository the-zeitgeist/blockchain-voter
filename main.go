package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/the-zeitgeist/voter/controllers"
)

func main() {
	fmt.Printf("Voter\n\n")

	vc, err := controllers.InitiateVoteChain()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome Voter!",
		})
	})

	r.POST("/transactions", controllers.AddTransactionHandler(vc))
	r.GET("/transactions", controllers.GetTransactionHandler(vc))
	r.GET("/chains/valid", controllers.ValidateChainHandler(vc))
	r.GET("/chains/result", controllers.ResultHandler(vc))

	fmt.Println("Initiating voter server...")
	r.Run()
}
