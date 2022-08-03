package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var privateThings = map[string]map[int64]string{
	"mike": {
		0: "MIKE: private string",
		1: "MIKE: secret thing",
		2: "MIKE: sneaky secret",
	},
	"rama": {
		0: "RAMA: private string",
		1: "RAMA: secret thing",
		2: "RAMA: sneaky secret",
	},
}

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type SignedResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func Index(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "index"})
}

func Private(c *gin.Context) {
	uidStr := c.Param("uid")
	pidInt, _ := strconv.ParseInt(c.Param("pid"), 10, 64)

	secret, ok := privateThings[uidStr][pidInt]

	if ok {
		c.JSON(200, gin.H{"msg": secret})
		return
	}

	c.JSON(200, gin.H{"msg": "unknown pid"})
}

func Login(c *gin.Context) {
	type login struct {
		Username string `json:"username,omitempty"`
	}

	loginParams := login{}
	c.ShouldBindJSON(&loginParams)

	if loginParams.Username == "mike" || loginParams.Username == "rama" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": loginParams.Username,
			"nbf":  time.Date(2022, 01, 01, 12, 0, 0, 0, time.UTC).Unix(),
		})

		tokenStr, err := token.SignedString([]byte("supersaucysecret"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, SignedResponse{
			Token:   tokenStr,
			Message: "logged in",
		})
		return
	}

	c.JSON(http.StatusBadRequest, UnsignedResponse{
		Message: "bad username",
	})
}

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func ParseToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte("supersaucysecret"), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}

func JwtTokenCheck(c *gin.Context) {
	jwtToken, err := ExtractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := ParseToken(jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "bad jwt token",
		})
		return
	}

	_, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		c.AbortWithStatusJSON(http.StatusInternalServerError, UnsignedResponse{
			Message: "unable to parse claims",
		})
		return
	}
	c.Next()
}

func PrivateACLCheck(c *gin.Context) {
	jwtToken, err := ExtractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := ParseToken(jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "bad jwt token",
		})
		return
	}

	claims, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		c.AbortWithStatusJSON(http.StatusInternalServerError, UnsignedResponse{
			Message: "unable to parse claims",
		})
		return
	}

	claimedUser, OK := claims["user"].(string)
	println(claimedUser)
	c.Next()
}
