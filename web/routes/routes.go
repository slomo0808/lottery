package routes

import (
	"github.com/kataras/iris/mvc"
	"lottery/bootstrap"
	"lottery/services"
	"lottery/web/controllers"
	"lottery/web/middleware"
)

func Configure(b *bootstrap.Bootstrapper) {
	userService := services.NewUserService()
	codeService := services.NewCodeService()
	resultService := services.NewResultService()
	giftService := services.NewGiftService()
	userdayService := services.NewUserdayService()
	blackipService := services.NewBlackipService()

	index := mvc.New(b.Party("/"))
	index.Register(
		userService,
		codeService,
		resultService,
		giftService,
		userdayService,
		blackipService)
	index.Handle(new(controllers.IndexController))

	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(
		userService,
		codeService,
		resultService,
		giftService,
		userdayService,
		blackipService)
	admin.Handle(new(controllers.AdminController))

	gift := admin.Party("/gift")
	gift.Register(giftService)
	gift.Handle(new(controllers.AdminGiftController))

	code := admin.Party("/code")
	code.Register(codeService, giftService)
	code.Handle(new(controllers.AdminCodeController))

	result := admin.Party("/result")
	result.Register(resultService)
	result.Handle(new(controllers.AdminResultController))

	user := admin.Party("/user")
	user.Register(userService)
	user.Handle(new(controllers.AdminUserController))

	blackip := admin.Party("/blackip")
	blackip.Register(blackipService)
	blackip.Handle(new(controllers.AdminBlackipController))
}
