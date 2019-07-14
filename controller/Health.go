package controller

import (
	"api/model"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strings"
	"time"
)

func (ctrl Controller) HelloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, model.Response{
		Message: fmt.Sprintf("Hello World at %s", time.Now().Format(time.Stamp)),
	})
}

func (ctrl Controller) Receiver(c echo.Context) error {
	var i struct {
		Random model.Random `json:"random" validate:"required"`
	}

	// bind input value to variable i
	if err := c.Bind(&i); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// normalize input values
	i.Random.Name = strings.TrimSpace(i.Random.Name)

	// validation input
	if err := c.Validate(i); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	return c.JSON(http.StatusOK, model.Response{
		Message: fmt.Sprintf("Your name is %s, and number is %d", i.Random.Name, i.Random.RandomNumber),
	})
}
