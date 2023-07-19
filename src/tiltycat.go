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
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func degToRad(degrees float64) (radians float64) {
	radians = degrees * (math.Pi / 180)
	return
}

type Point struct {
	x float64
	y float64
}

// Finds the new location of a point on a rectangle that has
// been rotated.
func calcLocationOfRotatedPoint(location, centreOfRotation Point, transSin, transCos float64) (newLocation Point) {
	return Point{
		x: (centreOfRotation.x + (location.x-centreOfRotation.x)*transCos + (location.y-centreOfRotation.y)*transSin),
		y: (centreOfRotation.y - (location.x-centreOfRotation.x)*transSin + (location.y-centreOfRotation.y)*transCos),
	}
}

// Calculates the bounding box of a square of width "width"
// rotated clockwise by "theta", in radians.
func boundingBox(theta float64, width float64) (minX, minY, maxX, maxY int) {
	s := math.Sin(theta)
	c := math.Cos(theta)

	centre := Point{
		x: width / 2,
		y: width / 2,
	}
	point1 := calcLocationOfRotatedPoint(Point{0, 0}, centre, s, c)
	point2 := calcLocationOfRotatedPoint(Point{0, width}, centre, s, c)
	point3 := calcLocationOfRotatedPoint(Point{width, 0}, centre, s, c)
	point4 := calcLocationOfRotatedPoint(Point{width, width}, centre, s, c)

	minX = int(math.Floor(math.Min(math.Min(point1.x, point2.x), math.Min(point3.x, point4.x))))
	minY = int(math.Floor(math.Min(math.Min(point1.y, point2.y), math.Min(point3.y, point4.y))))
	maxX = int(math.Ceil(math.Max(math.Max(point1.x, point2.x), math.Max(point3.x, point4.x))))
	maxY = int(math.Ceil(math.Max(math.Max(point1.y, point2.y), math.Max(point3.y, point4.y))))
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

	router.GET("/api/:angle", rotationHandler)
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	err := router.Run()
	log.Fatal(err)
}
