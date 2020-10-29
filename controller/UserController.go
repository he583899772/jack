package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"jack/common"
	"jack/dto"
	"jack/model"
	"jack/response"
	"jack/util"
	"net/http"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	//使用map 获取请求的参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	//
	var requestUser = model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	ctx.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码必须是6位以上")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码必须是6位以上"})
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	if isTelephoneExist(DB,telephone) {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户已存在")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已存在"})
		return
	}
	hasedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx,http.StatusInternalServerError,500,nil,"加密错误")
		//ctx.JSON(http.StatusInternalServerError,gin.H{"code":500,"msg":"加密错误"})
		return
	}
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hasedPassword),
	}
	DB.Create(&newUser)
	//插入数据

	//发放token
	token ,err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx,http.StatusInternalServerError,500,nil,"系统错误")
	}
	response.Success(ctx,gin.H{"token":token},"注册成功")
}

func Login(ctx *gin.Context)  {
	//获取参数
	DB := common.GetDB()
	var requestUser = model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	ctx.Bind(&requestUser)
	//获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码必须是6位以上")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码必须是6位以上"})
		return
	}
	//判断密码是否正确
	var user model.User
	DB.Where("telephone = ? ",telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户不存在")
		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password));err != nil {
		response.Response(ctx,http.StatusBadRequest,400,nil,"密码错误")
		//ctx.JSON(http.StatusBadRequest,gin.H{"code":400,"msg":"密码错误"})
		return
	}
	//发放token
	token ,err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx,http.StatusInternalServerError,500,nil,"系统错误")
	}
	//返回结果
	//ctx.JSON(200, gin.H{
	//	"code": 200,
	//	"data": gin.H{"token":token},
	//	"msg": "登陆成功",
	//})
	response.Success(ctx,gin.H{"token":token},"登陆成功")
}

func Info(ctx *gin.Context)  {
	user,_ := ctx.Get("user")

	ctx.JSON(http.StatusOK,gin.H{"ceode":200,"data":gin.H{"user":dto.ToUserDto(user.(model.User))}})
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user model.User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID != 0{
		return true
	}
	return false
}