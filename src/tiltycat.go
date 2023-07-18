package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/gin-gonic/gin"
)

func degToRad(degrees float64) (radians float64) {
	radians = degrees * (math.Pi / 180)
	return
}

func boundingBox(theta float64, width float64) (minX, minY, maxX, maxY int) {
	s := math.Sin(theta)
	c := math.Cos(theta)

	var x0, y0 = width / 2, width / 2
	var x1, y1 float64 = (x0 + (0-x0)*c + (0-y0)*s), (y0 - (0-x0)*s + (0-y0)*c)
	var x2, y2 float64 = (x0 + (0-x0)*c + (width-y0)*s), (y0 - (0-x0)*s + (width-y0)*c)
	var x3, y3 float64 = (x0 + (width-x0)*c + (0-y0)*s), (y0 - (width-x0)*s + (0-y0)*c)
	var x4, y4 float64 = (x0 + (width-x0)*c + (width-y0)*s), (y0 - (width-x0)*s + (width-y0)*c)

	minX = int(math.Floor(math.Min(math.Min(x1, x2), math.Min(x3, x4))))
	minY = int(math.Floor(math.Min(math.Min(y1, y2), math.Min(y3, y4))))
	maxX = int(math.Ceil(math.Max(math.Max(x1, x2), math.Max(x3, x4))))
	maxY = int(math.Ceil(math.Max(math.Max(y1, y2), math.Max(y3, y4))))
	return
}

func rotate(angle float64) (*image.RGBA, error) {
	imagePath, err := os.Open("img/tiltycat.png")
	if err != nil {
		return nil, err
	}

	defer imagePath.Close()
	srcImage, _, err := image.Decode(imagePath)
	if err != nil {
		return nil, err
	}

	rotated := image.NewRGBA(image.Rect(boundingBox(degToRad(angle), 128.0)))
	graphics.Rotate(rotated, srcImage, &graphics.RotateOptions{Angle: degToRad(angle)})

	return rotated, nil
}

func rotationHandler(c *gin.Context) {
	parameterNoPng := strings.Split(c.Param("angle"), ".png")[0]
	angle, err := strconv.ParseFloat(parameterNoPng, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Malformed URL, or invalid angle supplied.")
	}
	image, err := rotate(angle)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusInternalServerError, "Internal server error.")
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusInternalServerError, "Internal server error.")
	}

	imageBytes := buf.Bytes()

	c.Header("Content-Disposition", "inline")
	c.Data(http.StatusOK, "image/png", imageBytes)
}

func main() {
	router := gin.Default()

	router.GET("/:angle", rotationHandler)

	err := router.Run()
	log.Fatal(err)
}
