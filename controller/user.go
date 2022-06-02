package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验
	p := new(models.ParamSingUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SingUp with invalid param", zap.Error(err))

		//判断err是不是 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}

		//c.JSON(http.StatusOK, gin.H{
		//	//"msg": "请求参数有误",
		//	//"msg": errs.Translate(trans), //翻译错误
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})

		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//手动对请求参数进行详细的业务规则校验
	//if len(p.Username) == 0 || len(p.Password) == 0 ||
	//	len(p.RePassword) == 0 || p.RePassword != p.Password {
	//	//请求参数有误，直接返回响应
	//	zap.L().Error("SingUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	fmt.Println(p)

	//2.业务处理
	if err := logic.SingUp(p); err != nil {
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		zap.L().Error("login.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExit) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success...",
	//})
	ResponseSuccess(c, nil)

}

func LoginHandler(c *gin.Context) {
	//1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))

		//判断err是不是 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			ResponseError(c, CodeInvalidParam)
			return
		}

		//c.JSON(http.StatusOK, gin.H{
		//	//"msg": "请求参数有误",
		//	//"msg": errs.Translate(trans), //翻译错误
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed",
			zap.String("username", p.Username),
			zap.Error(err),
		)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或者密码错误",
		//})
		if errors.Is(err, mysql.ErrorUserNotExit) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "登录成功!",
	//})
	//ResponseSuccess(c, token)
	ResponseSuccess(c, gin.H{
		//"user_id":   user.UserID, //id值大于1<<53-1  int64类型的最大值是1<<63-1
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
