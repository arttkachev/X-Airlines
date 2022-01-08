package auth

import (
	"context"
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arttkachev/X-Airlines/Backend/api/models/user"
	"github.com/arttkachev/X-Airlines/Backend/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

type AuthService struct{}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (handler *AuthService) SignUp(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user.ID = primitive.NewObjectID()
	var userService = services.GetUserService()
	collection := userService.Collection
	if user.Airlines == nil {
		user.Airlines = make([]primitive.ObjectID, 0)
	}
	h := sha256.New()
	_, err := collection.InsertOne(ctx, bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": string(h.Sum([]byte(user.Password))),
		"isAdmin":  user.IsAdmin,
		"balance":  user.Balance,
		"airlines": user.Airlines,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// clear cache
	log.Println("Remove user data from Redis")
	userService.RedisClient.Del("users")
	c.JSON(http.StatusOK, user)
}

func (handler *AuthService) SignIn(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userService := services.GetUserService()
	h := sha256.New()
	cur := userService.Collection.FindOne(ctx, bson.M{
		"name":     user.Name,
		"password": string(h.Sum([]byte(user.Password))),
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password"})
		return
	}
	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("name", user.Name)
	session.Set("token", sessionToken)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "User signed in"})

	// // Expiration time for a token is gonna be 10 mins
	// expirationTime := time.Now().Add(10 * time.Minute)
	// claims := &Claims{
	// 	Username: user.Name,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirationTime.Unix(),
	// 	},
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error()})
	// 	return
	// }
	// jwtOutput := JWTOutput{
	// 	Token:   tokenString,
	// 	Expires: expirationTime,
	// }
	// c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthService) SignOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"message": "Signed out"})
}

func (handler *AuthService) Refresh(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error()})
		return
	}
	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token"})
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is not expired yet"})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged in",
			})
			c.Abort()
		}
		c.Next()

		// tokenValue := c.GetHeader("Authorization")
		// claims := &Claims{}
		// tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("JWT_SECRET")), nil
		// })
		// if err != nil {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
		// if tkn == nil || !tkn.Valid {
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// }
		// c.Next()
	}
}
