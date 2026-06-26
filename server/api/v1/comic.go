package v1

import (
	"strconv"

	"blog/internal/dto"
	"blog/internal/service"
	"blog/pkg/response"

	"github.com/gin-gonic/gin"
)

type ComicController struct {
	service *service.ComicService
}

func NewComicController(service *service.ComicService) *ComicController {
	return &ComicController{service: service}
}

// ListForWeb 前端番剧列表
func (ctrl *ComicController) ListForWeb(c *gin.Context) {
	comics, err := ctrl.service.ListForWeb()
	if err != nil {
		response.Failed(c, err.Error())
		return
	}
	response.Success(c, comics)
}

// AdminList 管理端番剧列表
func (ctrl *ComicController) AdminList(c *gin.Context) {
	var req dto.ListComicRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Failed(c, "参数错误")
		return
	}
	comics, total, err := ctrl.service.List(&req)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}
	response.Success(c, gin.H{"list": comics, "total": total})
}

// Get 管理端获取番剧详情
func (ctrl *ComicController) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, "参数错误")
		return
	}
	comic, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		response.Failed(c, "番剧不存在")
		return
	}
	response.Success(c, comic)
}

// Create 管理端创建番剧
func (ctrl *ComicController) Create(c *gin.Context) {
	var req dto.CreateComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, "参数错误")
		return
	}
	comic, err := ctrl.service.Create(&req)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}
	response.Success(c, comic)
}

// Update 管理端更新番剧
func (ctrl *ComicController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, "参数错误")
		return
	}
	var req dto.UpdateComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, "参数错误")
		return
	}
	comic, err := ctrl.service.Update(uint(id), &req)
	if err != nil {
		response.Failed(c, err.Error())
		return
	}
	response.Success(c, comic)
}

// Delete 管理端删除番剧
func (ctrl *ComicController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Failed(c, "参数错误")
		return
	}
	if err := ctrl.service.Delete(uint(id)); err != nil {
		response.Failed(c, err.Error())
		return
	}
	response.Success(c, nil)
}
