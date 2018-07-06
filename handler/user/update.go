package user

import (
	"apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/errno"
	"apiserver/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

//Update update a exist user account info.

func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	userId, _ := strconv.Atoi(c.Param("id"))

	//Binding the user data
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	// update the record based on the user id.
	user, err := model.GetUserByID(userId)
	if err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
	}
	user.Username = u.Username
	user.Password = u.Password
	//Validate the data.
	//Encrypt the user password.
	if err := user.Encrypt(); err != nil {
		handler.SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	//Save changed fields
	if err := user.Update(); err != nil {
		log.Info(err.Error())
		handler.SendResponse(c, errno.ErrDatabase, nil)
		return
	}
	handler.SendResponse(c, nil, nil)
}
