package main

import (
	"FinProject/docs"
	modelz "FinProject/models"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Profile struct {
	Nama        string    `json:"Name" example:"Jemmi"`
	Date        time.Time `example:"0001-01-01T07:00:00+07:00"`
	DateInput   string    `json:"DateInput" example:"2022-07-05"`
	Description []string  `json:"Description" example:"Bekerja"`
}

type Response struct {
	Code    string      `json:"code" example:"00"`
	Message string      `json:"message" example:"Succesful"`
	Data    interface{} `json:"data"`
}
type Profiles struct {
	Profile []Profile
}

func init() {
	auths = map[string]string{
		"dicky123": "pratama123",
		"Jemmi123": "Epriadi123",
	}
}

var auths map[string]string

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := Response{}
		user, password, ok := c.Request.BasicAuth()
		fmt.Println("User => ", user)
		fmt.Println("Password => ", password)
		if !ok {
			res.Code = "02"
			res.Message = "Invalid authentication"
			c.JSON(http.StatusOK, res)
			c.Abort()
			return
		}
		if _, ok := auths[user]; !ok {
			res.Code = "02"
			res.Message = "Invalid authentication"
			c.JSON(http.StatusOK, res)
			c.Abort()
			return
		}
		if auths[user] != password {
			res.Code = "02"
			res.Message = "Invalid authentication"
			c.JSON(http.StatusOK, res)
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	docs.SwaggerInfo.Title = "ToDo REST API"
	docs.SwaggerInfo.Description = "Simple Example of ToDo API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:7070"
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	modelz.ConnectDataBase()

	// /profile/dashboard/get
	profile := r.Group("/todos")
	profile.Use(CheckAuth())
	profile.GET("/", getAPI)
	profile.GET("/:id", getById)
	profile.POST("/", postAPI)
	profile.PUT("/:id", updateById)
	profile.DELETE("/:id", deleteById)

	err := r.Run(":7070")
	if err != nil {
		fmt.Println("error exit: ", err)
	}
}

var res = Response{
	Code:    "00",
	Message: "Success",
}

// GetProfile godoc
// @Summary Get User Todo List
// @Description get user ToDo List
// @ID get-profile
// @Produce  json
// @Success 200 {object} Response{data=modelz.Profile}
// @Router /ToDos [get]
func getAPI(c *gin.Context) {
	var dataProfile []modelz.Profile
	modelz.DB.Find(&dataProfile)
	res.Data = dataProfile
	c.JSON(http.StatusOK, dataProfile)
}

func (box *Profiles) AddItem(item Profile) []Profile {
	box.Profile = append(box.Profile, item)
	return box.Profile
}

// CreateToDo godoc
// @Summary Create a ToDo
// @Description for save user profile data
// @ID save-profile
// @Accept  json
// @Produce  json
// @Param requestbody body Profile true "request body json"
// @Success 200 {object} Response{}
// @Router /ToDo [post]
func postAPI(c *gin.Context) {
	var req Profile

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to profile struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"

		c.JSON(http.StatusBadRequest, res)
		return
	}

	//NOTE: Ketika melakukan input tanggal, tolong ikuti format ini.
	format := "2006-01-02"

	s1, _ := time.Parse(format, req.DateInput)
	requestDescription := strings.Join(req.Description, ",")
	dataProfile := modelz.Profile{Nama: req.Nama, Description: requestDescription, Date: s1}
	modelz.DB.Create(&dataProfile)

	// byteslice, err:= json.Marshal(req)
	res.Data = dataProfile
	c.JSON(http.StatusOK, res)
	return
}

// UpdateToDo godoc
// @Summary Update a ToDo
// @Description for update user ToDo data
// @ID update-Todo
// @Accept  json
// @Produce  json
// @Param requestbody body Profile true "request body json"
// @Success 200 {object} Response{}
// @Router /ToDo/id [put]
func updateById(c *gin.Context) {
	var dataProfile modelz.Profile
	if err := modelz.DB.First(&dataProfile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile not found!"})
		return
	}

	var req Profile
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Failed to parse request to profile struct: ", err)
		res.Code = "02"
		res.Message = "Failed parsing request"

		c.JSON(http.StatusBadRequest, res)
		return
	}
	requestDescription := strings.Join(req.Description, ",")
	modelz.DB.Model(&dataProfile).Where("id = ?", c.Param("id")).Updates(modelz.Profile{Nama: req.Nama, DateInput: req.DateInput, Description: requestDescription})
	res.Data = dataProfile
	c.JSON(http.StatusOK, res)
}

// GetProfile godoc
// @Summary Get User Todo List
// @Description get user ToDo List
// @ID getById-profile
// @Produce  json
// @Success 200 {object} Response{data=modelz.Profile}
// @Router /ToDos/id [get]
func getById(c *gin.Context) {
	var dataProfile modelz.Profile
	if err := modelz.DB.First(&dataProfile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile not found!"})
		return
	}
	res.Data = dataProfile
	c.JSON(http.StatusOK, res)
}

// DeleteToDo godoc
// @Summary Delete a ToDo
// @Description for deleting user profile data
// @ID delete-ToDo
// @Accept  json
// @Produce  json
// @Param requestbody body Profile true "request body json"
// @Success 200 {object} Response{}
// @Router /ToDo/id [delete]
func deleteById(c *gin.Context) {
	var dataProfile modelz.Profile
	if err := modelz.DB.First(&dataProfile, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile not found!"})
		return
	}
	modelz.DB.Delete(&dataProfile, c.Param("id"))
	c.JSON(http.StatusOK, "Delete succesfull")
}
