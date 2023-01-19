package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/manifoldco/promptui"
	"github.com/the-zeitgeist/voter/constants"
	"github.com/the-zeitgeist/voter/models"
)

func InitiateVoteChain() (*models.VoteChain, error) {
	configFilePath := constants.ConfigFile
	config, err := os.ReadFile(configFilePath)

	if err != nil {
		if err.Error() != fmt.Sprintf("open %s: no such file or directory", configFilePath) {
			return nil, err
		}

		prompt := promptui.Select{
			Label: "This node is not configured yet, what would you like to do?",
			Items: []string{"Join", "Create"},
		}

		index, _, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		var vc *models.VoteChain

		if index == 0 {
			vc, err = models.JoinVoteChain()
		} else {
			vc, err = models.NewVoteChain()
		}

		if err != nil {
			fmt.Printf("Could not join chain: %s\n", err.Error())
			return nil, err
		}

		return vc, nil
	}

	var vc models.VoteChain
	err = json.Unmarshal(config, &vc)
	if err != nil {
		return nil, err
	}

	vc.IsProcessing = false
	return &vc, nil
}

func ValidateChainHandler(vc *models.VoteChain) func(c *gin.Context) {
	return func(c *gin.Context) {
		isValid := vc.ValidateChain()

		c.JSON(http.StatusOK, gin.H{"isValid": isValid, "chain": vc.Chain})
	}
}

func ResultHandler(vc *models.VoteChain) func(c *gin.Context) {
	return func(c *gin.Context) {
		res := vc.Result()

		c.JSON(http.StatusOK, gin.H{"result": res})
	}
}
