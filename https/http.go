package https

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"simple-gateway/infra"
	"simple-gateway/internal/adapters/handler"
	"simple-gateway/internal/adapters/repository"
	"simple-gateway/internal/core/services"
	"simple-gateway/middleware/clientHttp"
)

func InitRoutes() {
	r := gin.New()

	bindVerification(r)
	bindConfirmation(r)

	appPort := viper.GetString("service.port")
	err := r.Run(":" + appPort)
	if err != nil {
		log.Error().Err(err).Msg("Start run error")
	}
}

func bindVerification(r *gin.Engine) {
	validate := validator.New()
	client := clientHttp.NewClient(clientHttp.NewClientConfig("integration.prompay"))
	prompay := repository.NewPromPay(client)
	transactionDb := repository.NewTransactionRepo(infra.TransactionDB)
	verificationSvc := services.NewVerificationSvc(prompay, validate, transactionDb)
	verificationHdl := handler.NewVerificationHdl(verificationSvc)

	r.POST("/verification", verificationHdl.Verification)
}

func bindConfirmation(r *gin.Engine) {
	validate := validator.New()
	clientPrompay := clientHttp.NewClient(clientHttp.NewClientConfig("integration.prompay"))
	prompay := repository.NewPromPay(clientPrompay)
	clientNoti := clientHttp.NewClient(clientHttp.NewClientConfig("integration.notification"))
	noti := repository.NewNotification(clientNoti)
	transactionDb := repository.NewTransactionRepo(infra.TransactionDB)
	confirmSvc := services.NewConfirmationSvc(prompay, validate, transactionDb, noti)
	confirmHdl := handler.NewConfirmationHdl(confirmSvc)

	r.POST("/confirmation", confirmHdl.Confirmation)
}
