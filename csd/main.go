package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Zone struct {
	Name        string   `json:"name"`
	NameServers []string `json:"name_servers"`
	Owner       string   `json:"owner"`
}

// e.POST("/zones", saveZone)
func saveZone(c echo.Context) error {
	zone := new(Zone)
	if err := c.Bind(zone); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, zone)
}

// e.GET("/zones/:name", getZone)
func getZone(c echo.Context) error {
	name := c.Param("name")
	if name == "notfound" {
		return c.JSON(http.StatusNotFound, nil)
	} else {
		return c.JSON(http.StatusOK, Zone{
			Name:        name,
			NameServers: []string{"ns1.aws.example.com", "ns2.aws.example.net"},
			Owner:       "123456789",
		})
	}
}

// e.GET("/zones", getZones)
func getZones(c echo.Context) error {
	return c.JSON(http.StatusOK, []Zone{
		Zone{
			Name:        "jira",
			NameServers: []string{"ns1.aws.example.com", "ns2.aws.example.net"},
			Owner:       "123456789",
		},
		Zone{
			Name:        "confluence",
			NameServers: []string{"ns23.aws.example.com", "ns42.aws.example.net"},
			Owner:       "987654321",
		},
	})
}

// e.PUT("/zones/:name", updateZone)
func updateZone(c echo.Context) error {
	name := c.Param("name")
	if name == "notfound" {
		return c.JSON(http.StatusNotFound, nil)
	} else {
		update := new(Zone)
		if err := c.Bind(update); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusCreated, Zone{
			Name:        name,
			NameServers: update.NameServers,
			Owner:       update.Owner,
		})
	}
}

// e.DELETE("/zones/:name", deleteZone)
func deleteZone(c echo.Context) error {
	name := c.Param("name")
	if name == "notfound" {
		return c.JSON(http.StatusNotFound, nil)
	} else {
		return c.JSON(http.StatusOK, nil)
	}
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.POST("/zones", saveZone)
	e.GET("/zones", getZones)
	e.GET("/zones/:name", getZone)
	e.PUT("/zones/:name", updateZone)
	e.DELETE("/zones/:name", deleteZone)

	e.Logger.Fatal(e.Start(":8080"))
}
