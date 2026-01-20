package gormutil

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// BaseHandler 基础处理层, 只用于快速测试，不建议在生产环境使用
//
//	T: 实体模型类型
type BaseHandler[T any] struct {
	Svc *BaseService[T]
}

// NewBaseHandler 构造函数，完成整个链路的组装
func NewBaseHandler[T any](db *gorm.DB) *BaseHandler[T] {
	return &BaseHandler[T]{
		Svc: NewBaseService[T](db),
	}
}

// RegisterRoutes 批量注册 RESTful 路由
// 注意：GET /path 必须放在 GET /path/:id 之前
func (h *BaseHandler[T]) RegisterRoutes(r *echo.Group, path string) {
	// CURD
	r.PUT(path+"/:id", h.Update)
	r.GET(path+"/:id", h.GetByID)
	r.POST(path, h.Create)
	r.DELETE(path+"/:id", h.Delete)
	// pages
	r.GET(path, h.List) // 列表
	// excel
	//r.GET(path+"/export", h.Export)  // 导出
	//r.POST(path+"/import", h.Import) // 导入

	// 批量操作路由
	r.POST(path+"/batch", h.BatchCreate)
	r.PUT(path+"/batch", h.BatchUpdate)
	r.DELETE(path+"/batch", h.BatchDelete)
}

// ================= CRUD 方法 =================

// Create 创建
func (h *BaseHandler[T]) Create(e echo.Context) error {
	var entity T
	if err := h.Svc.Create(e.Request().Context(), &entity); err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "创建失败: " + err.Error()})
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "data": entity})
	return nil
}

// GetByID 根据ID查询
func (h *BaseHandler[T]) GetByID(e echo.Context) error {
	id := e.Param("id")

	entity, err := h.Svc.GetByID(e.Request().Context(), id)
	if err != nil {
		// 这里简单处理为 404，生产环境可判断 gorm.ErrRecordNotFound
		e.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "msg": "数据不存在"})
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "data": entity})
	return nil
}

// Update 更新
func (h *BaseHandler[T]) Update(e echo.Context) error {
	var entity T
	if err := e.Bind(&entity); err != nil {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "参数错误: " + err.Error()})
		return err
	}

	if err := h.Svc.Update(e.Request().Context(), &entity); err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "更新失败: " + err.Error()})
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "更新成功"})
	return nil
}

// Delete 删除
func (h *BaseHandler[T]) Delete(e echo.Context) error {
	id := e.Param("id")

	if err := h.Svc.Delete(e.Request().Context(), id); err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "删除失败: " + err.Error()})
		return err
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "删除成功"})
	return nil
}

// List 分页列表
func (h *BaseHandler[T]) List(e echo.Context) error {
	page, _ := strconv.Atoi(e.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(e.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	// 调用 Repo 层的分页方法
	list, total, err := h.Svc.List(e.Request().Context(), page, pageSize)
	if err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "查询失败: " + err.Error()})
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
	return nil
}

// 批量创建

// CreateBatch 批量创建 (接收 JSON 数组)
func (h *BaseHandler[T]) CreateBatch(e echo.Context) error {
	var entities []*T
	if err := e.Bind(&entities); err != nil {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": err.Error()})
		return nil
	}

	if err := h.Svc.CreateBatch(e.Request().Context(), entities, 100); err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": err.Error()})
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "成功导入 " + strconv.Itoa(len(entities)) + " 条数据"})
	return nil
}

// BatchCreate 批量创建
func (h *BaseHandler[T]) BatchCreate(e echo.Context) error {
	var entities []*T
	if err := e.Bind(&entities); err != nil {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "参数格式错误"})
		return nil
	}

	if err := h.Svc.CreateBatch(e.Request().Context(), entities, 100); err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "批量创建失败"})
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "成功导入 " + strconv.Itoa(len(entities)) + " 条数据"})
	return nil
}

// BatchUpdate 批量更新
// Body 结构示例: {"ids": ["id1", "id2"], "data": {"status": 1}}
func (h *BaseHandler[T]) BatchUpdate(e echo.Context) error {
	var req struct {
		IDs  []string       `json:"ids"`
		Data map[string]any `json:"data"`
	}

	if err := e.Bind(&req); err != nil {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "参数错误"})
		return nil
	}

	if len(req.IDs) == 0 || len(req.Data) == 0 {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "ids 或 data 不能为空"})
		return nil
	}

	// 调用 Service 层的批量更新方法
	//if err := h.Svc.UpdateBatch(e.Request().Context(), req.IDs, req.Data, 100); err != nil {
	//	e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "批量更新失败"})
	//	return nil
	//}
	//e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "成功更新 " + strconv.Itoa(len(req.IDs)) + " 条数据"})
	return nil
}

// BatchDelete 批量删除
// Body 结构示例: {"ids": ["id1", "id2"]}
func (h *BaseHandler[T]) BatchDelete(e echo.Context) error {
	var req struct {
		IDs []string `json:"ids"`
	}

	if err := e.Bind(&req); err != nil {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "参数错误"})
		return nil
	}

	if len(req.IDs) == 0 {
		e.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": "ids 不能为空"})
		return nil
	}

	if err := h.Svc.DeleteBatch(e.Request().Context(), req.IDs); err != nil {
		e.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": "批量删除失败"})
		return nil
	}

	e.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "成功删除 " + strconv.Itoa(len(req.IDs)) + " 条数据"})
	return nil
}

// 导入导出

// ================= 导出方法 =================

// Export 导出 Excel
func (h *BaseHandler[T]) Export(e echo.Context) error {
	// 1. 查询数据 (这里为了演示只查前 1000 条，实际应根据 Query 参数筛选)
	list, _, err := h.Svc.ListByCondition(e.Request().Context(), nil, 1, 1000)
	if err != nil {
		e.String(http.StatusInternalServerError, err.Error())
		return nil
	}

	// 2. 创建 Excel
	f := excelize.NewFile()
	sheetName := "Sheet1"
	_ = f.SetSheetName("Sheet1", sheetName)

	// 3. 写入表头和数据
	// 注意：这里利用反射动态获取字段名作为表头，适用于任意结构体
	if len(list) > 0 {
		t := reflect.TypeOf(list[0])
		// 写入表头
		for i := 0; i < t.NumField(); i++ {
			tag := t.Field(i).Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, tag)
		}

		// 写入数据
		for row, item := range list {
			v := reflect.ValueOf(item).Elem()
			for col := 0; col < v.NumField(); col++ {
				// 获取 json tag 决定是否导出该列
				tag := t.Field(col).Tag.Get("json")
				if tag == "" || tag == "-" {
					continue
				}

				val := v.Field(col).Interface()
				cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
				f.SetCellValue(sheetName, cell, val)
			}
		}
	}
	// 4. 输出
	e.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	e.Response().Header().Set("Content-Disposition", "attachment; filename=export.xlsx")
	f.Write(e.Response().Writer)

	return nil
}
