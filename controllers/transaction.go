package controllers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/the-zeitgeist/voter/models"
)

var h hash.Hash = sha256.New()

func AddTransactionHandler(vc *models.VoteChain) func(c *gin.Context) {
	return func(c *gin.Context) {
		var tx models.Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		if tx.Candidate == "" || tx.Voter == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": errors.New("malformed Transaction").Error()})
			return
		}

		validCandidate := false
		for _, v := range vc.Candidates {
			if v.Id == tx.Candidate {
				validCandidate = true
				break
			}
		}

		if !validCandidate {
			c.JSON(http.StatusBadRequest, gin.H{"Error": errors.New("invalid candidate was provided").Error()})
			return
		}

		if tx.Id == "" {
			h.Reset()
			h.Write([]byte(tx.Voter))
			id := fmt.Sprintf("%x", (h.Sum(nil)))
			tx.Id = id[0 : len(id)/2]
			tx.Voter = id
			tx.Timestamp = time.Now().Unix()
		}

		if _, ok := vc.Txs[tx.Id]; ok {
			c.JSON(http.StatusBadRequest, gin.H{"Error": errors.New("voter has a pending Transaction").Error()})
			return
		}

		for _, b := range vc.Chain {
			if _, ok := b.Transactions[tx.Id]; ok {
				c.JSON(http.StatusBadRequest, gin.H{"Error": errors.New("voter has already been proccessed").Error()})
				return
			}
		}

		vc.AddTx(tx)
		c.JSON(http.StatusCreated, gin.H{"Transaction": tx.Id})
	}
}

func GetTransactionHandler(vc *models.VoteChain) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Transactions": vc.Txs})
	}
}
