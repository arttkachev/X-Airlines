package controllers

import (
	"net/http"

	"github.com/arttkachev/X-Airlines/Backend/api/models"
	userService "github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func AddUser(c *gin.Context) {

	var user models.User
	err := c.ShouldBindJSON(&user) // ShouldBindJSON marshals the incoming request body into a struct passed in as an argument (it's user in our case)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	user.ID = xid.New().String() // xid lib set a unique ID
	users := userService.GetRepository()

	*users = append(*users, user)
	c.JSON(http.StatusOK, user) // sends a response with httpStatusOK and a newly created user as a JSON
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userService.GetRepository()) //marshals the users array to JSON
}
