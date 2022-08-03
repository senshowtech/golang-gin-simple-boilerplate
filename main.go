// package main

// import (
// 	config "goGin/config"
// 	dogController "goGin/controller/dog"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()
// 	config.Connect()
// 	r.GET("/test", dogController.GetDogs)
// 	if err := http.ListenAndServe(":5000", r); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	middleware "goGin/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	router.POST("/login", middleware.Login)
	router.GET("/", middleware.Index)

	privateRouter := router.Group("/private")
	privateRouter.Use(middleware.JwtTokenCheck)
	privateRouter.Use(middleware.PrivateACLCheck).GET("/test", middleware.Index)
	privateRouter.Use(middleware.PrivateACLCheck).GET("/:uid/:pid", middleware.Private)

	router.Run(":5000")
}
