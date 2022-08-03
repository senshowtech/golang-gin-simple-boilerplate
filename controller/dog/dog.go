package dog

import (
	config "goGin/config"
	dog "goGin/model/dog"

	"github.com/gin-gonic/gin"
)

func GetDogs(c *gin.Context) {
	var dogs []dog.Dog
	config.Database.Find(&dogs)
	c.JSON(200, gin.H{
		"data": dogs,
	})
}
