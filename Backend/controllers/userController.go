package controllers

import (
	"net/http"
	"strings"

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

	// add a new user
	*userService.GetRepository() = append(*userService.GetRepository(), user)
	c.JSON(http.StatusOK, user) // sends a response with httpStatusOK and a newly created user as a JSON
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userService.GetRepository()) //marshals the users array to JSON
}

func GetUserByAirline(c *gin.Context) {
	users := make([]models.User, 0)
	airline := c.Query("airlines")
	// loop through all users
	for i := 0; i < len(*userService.GetRepository()); i++ {
		// check user's airlines
		for _, t := range (*userService.GetRepository())[i].Airlines {
			// if requested airline is found in the list of this user, then this user is a right one for filtering
			if strings.EqualFold(t, airline) {
				// add to list of filtered users
				users = append(users, (*userService.GetRepository())[i])
			}
		}
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Users not found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context) {
	// fetch the user id from the request URL
	id := c.Param("id")

	// convert the request body into User struct. It assigns request body (new user info to the user var)
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	// loop through the list of users
	for i := 0; i < len(*userService.GetRepository()); i++ {
		// check if it's a right one by id
		if (*userService.GetRepository())[i].ID == id {
			// if found, set a new user info (per/field to save ID) and send him to the client
			(*userService.GetRepository())[i].Airlines = user.Airlines
			(*userService.GetRepository())[i].Balance = user.Balance
			(*userService.GetRepository())[i].Email = user.Email
			(*userService.GetRepository())[i].IsAdmin = user.IsAdmin
			(*userService.GetRepository())[i].Name = user.Name
			(*userService.GetRepository())[i].Password = user.Password
			user = (*userService.GetRepository())[i]
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found"})

}

func DeleteUser(c *gin.Context) {
	// fetch the user id from the request URL
	id := c.Param("id")

	// loop through the list of users
	for i := 0; i < len(*userService.GetRepository()); i++ {
		// check if it's a right one by id
		if (*userService.GetRepository())[i].ID == id {
			// if found, delete him from the slice by setting on his place all users going next (rearrangement)
			//[:i] captures all what was before i
			// [i+1:]... captures all what goes after i+1
			*userService.GetRepository() = append((*userService.GetRepository())[:i], (*userService.GetRepository())[i+1:]...)

			// send a message to a client
			c.JSON(http.StatusOK, gin.H{
				"message": "User has been deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found"})
}
