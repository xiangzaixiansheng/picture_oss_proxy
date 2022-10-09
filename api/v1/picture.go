package v1

import (
	"picture-oss-proxy/service"

	"github.com/gin-gonic/gin"
)

func GetPicture(c *gin.Context) {
	var pictureService service.PictureService
	res := pictureService.Get()
	c.JSON(200, res)
}
