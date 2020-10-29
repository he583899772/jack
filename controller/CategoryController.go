package controller

import (
	"github.com/gin-gonic/gin"
	"jack/model"
	"jack/repository"
	"jack/response"
	"jack/vo"
	"log"
	"strconv"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.ICategoryRepository
}

func (c CategoryController) Create(ctx *gin.Context) {
	// 获取参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		log.Println(err.Error())
		response.Fail(ctx,nil,"数据验证错误")
		return
	}
	category,err := c.Repository.Create(requestCategory.Name)
	if err != nil {
		response.Fail(ctx,nil,"数据验证错误")
		return
	}
	response.Success(ctx,gin.H{"category":category},"创建成功")
}

func (c CategoryController) Update(ctx *gin.Context) {
	// 获取参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		log.Println(err.Error())
		response.Fail(ctx,nil,"数据验证错误")
		return
	}
	categoryId,_ := strconv.Atoi(ctx.Params.ByName("id"))
	var updateCategory *model.Category
	category ,err := c.Repository.SelectById(categoryId)
	if err != nil {
		panic(err)
	}
	updateCategory,err = c.Repository.Update(*category,requestCategory.Name)
	if err != nil{
		panic(err)
	}
	response.Success(ctx,gin.H{"category": updateCategory},"更新成功")
}

func (c CategoryController) Show(ctx *gin.Context) {
	categoryId,_ := strconv.Atoi(ctx.Params.ByName("id"))
	var category *model.Category
	category ,err := c.Repository.SelectById(categoryId)
	if err != nil {
		panic(err)
	}
	response.Success(ctx,gin.H{"category": category},"")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	categoryId,_ := strconv.Atoi(ctx.Params.ByName("id"))
	c.Repository.DeleteById(categoryId)
	response.Success(ctx,nil,"删除成功")
}

func NewCategoryController() ICategoryController {
	categoryController := CategoryController{
		Repository: repository.NewCategoryRepository(),
	}
	categoryController.Repository.(
		repository.CategoryRepository).DB.AutoMigrate(model.Category{
	})
	return categoryController
}