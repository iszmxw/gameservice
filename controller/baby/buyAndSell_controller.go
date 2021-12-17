package baby

import (
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
	"redisData/controller"
	"redisData/model"
	"redisData/pkg/mysql"
)

type Controller struct {
}

// 买卖数据

func (h *Controller) BuyAndSellHandler(c *gin.Context) {
	var (
		params     model.ParamsBuyAndSellQuery // 接收请求参数
		BuyAndSell []model.RespBabyOrder       // 查询数据
	)
	// 绑定参数
	_ = c.Bind(&params)
	where := cmap.New().Items()
	if len(params.Name) > 0 {
		where["name"] = params.Name
	}
	if len(params.TokenId) > 0 {
		where["token_id"] = params.TokenId
	}
	if len(params.Status) > 0 {
		where["status"] = params.Status
	}
	mysql.DB.Debug().Model(model.BabyOrder{}).Where(where).Find(&BuyAndSell)
	controller.ResponseSuccess(c, BuyAndSell)
}
