package controllers

import (
	"fmt"
	"strings"

	"corpweb/models"

	"github.com/astaxie/beego"
)

type FlagshipProductController struct {
	// beego.Controller
	baseController
}

// ToFlagshipProducts 跳转到首页当季主打产品管理页
func (this *FlagshipProductController) ToFlagshipProducts() {
	this.Data["menu_lv_1"] = "home"
	this.Data["menu_lv_2"] = "flagship"

	// fProds ,err:= GetAllFlagshipProducts(false)
	this.TplNames = "admin/home_flagship_product.html"
}

// GetFlagshipProducts ajax获取主打产品列表
func (this *FlagshipProductController) GetFlagshipProducts() {
	resp := AjaxFormResp{
		Result: RESULT_RESP_FAIL,
		Msg:    this.Tr("tips_action_fail"),
	}
	defer func() {
		this.Data["json"] = &resp
		this.ServeJson(true)
	}()

	curPageNo, _ := this.GetInt("curPageNo")
	pageSize, _ := this.GetInt("pageSize")

	cond := models.FlagshipProduct{}
	if err := this.ParseForm(&cond); err != nil {
		beego.Error(fmt.Printf("this.ParseForm() err: %s", err))
		resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		return
	}

	page, err := models.GetFlagshipProductPage(&cond, curPageNo, pageSize, false)
	if err != nil {
		resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		return
	}

	if page == nil {
		page = models.EmptyPage(curPageNo, pageSize)
	}

	resp.ExtObj = page
	resp.Result = RESULT_RESP_SUCC
	resp.Msg = this.Tr("tips_action_success")
}

// AjaxGetProductsButFlagships ajax方式获取产品（不包含主打产品）分页信息
func (this *FlagshipProductController) AjaxGetProductsButFlagships() {
	resp := AjaxFormResp{
		Result: RESULT_RESP_FAIL,
		Msg:    this.Tr("tips_action_fail"),
	}
	defer func() {
		this.Data["json"] = &resp
		this.ServeJson(true)
	}()

	curPageNo, _ := this.GetInt("curPageNo")
	pageSize, _ := this.GetInt("pageSize")

	cond := models.Product{}
	if err := this.ParseForm(&cond); err != nil {
		beego.Error(fmt.Printf("this.ParseForm() err: %s", err))
		resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		return
	}

	page, err := models.GetProductsButFlagshipsPage(&cond, curPageNo, pageSize, false)
	if err != nil {
		resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		return
	}

	if page == nil {
		page = models.EmptyPage(curPageNo, pageSize)
	}

	resp.ExtObj = page
	resp.Result = RESULT_RESP_SUCC
	resp.Msg = this.Tr("tips_action_success")

}

// AddFlagshipProducts ajax添加主打产品
func (this *FlagshipProductController) AddFlagshipProducts() {
	resp := AjaxFormResp{
		Result: RESULT_RESP_FAIL,
		Msg:    this.Tr("tips_action_fail"),
	}
	defer func() {
		this.Data["json"] = &resp
		this.ServeJson(true)
	}()

	prodIds := make([]string, 0, 3)
	this.Ctx.Input.Bind(&prodIds, "prodIds")
	// prodIds := this.GetStrings("prodIds")
	if len(prodIds) <= 0 {
		return
	}

	fProds := make([]*models.FlagshipProduct, 0, len(prodIds))
	for _, prodId := range prodIds {
		fProd := &models.FlagshipProduct{
			ProductId: prodId,
			SortNo:    0,
		}
		fProds = append(fProds, fProd)
	}

	_, err := models.AddFlagshipProducts(fProds...)
	if err != nil {
		beego.Error(err)
		if strings.Contains(fmt.Sprintf("%s", err), "Duplicate entry") {
			resp.Msg = this.Tr("tips_duplicate_entry")
		} else {
			resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		}
		return
	}
	resp.Result = RESULT_RESP_SUCC
	resp.Msg = this.Tr("tips_action_success")
}

// RmFlagshipProduct ajax删除主打产品
func (this *FlagshipProductController) RmFlagshipProduct() {
	resp := AjaxFormResp{
		Result: RESULT_RESP_FAIL,
		Msg:    this.Tr("tips_action_fail"),
	}
	defer func() {
		this.Data["json"] = &resp
		this.ServeJson(true)
	}()

	id := this.GetString("id")
	_, err := models.RmFlagshipProductById(id)
	if err != nil {
		beego.Error(fmt.Sprintf("models.RmFlagshipProductById(%s) err: %s\n", id, err))
		resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		return
	}
	resp.Result = RESULT_RESP_SUCC
	resp.Msg = this.Tr("tips_action_success")
}

// PushPin 置顶
func (this *FlagshipProductController) PushPin() {
	resp := AjaxFormResp{
		Result: RESULT_RESP_FAIL,
		Msg:    this.Tr("tips_action_fail"),
	}
	defer func() {
		this.Data["json"] = &resp
		this.ServeJson(true)
	}()

	fProdId := this.GetString("fProdId")
	affected, err := models.PushPinFlagshipProduct(fProdId)
	if err != nil {
		resp.Msg = this.Tr("tips_sys_err_and_contact_tech")
		return
	}

	resp.Result = RESULT_RESP_SUCC
	resp.Msg = this.Tr("tips_action_success")
	resp.ExtObj = affected
}
