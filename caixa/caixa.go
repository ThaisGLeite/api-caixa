package caixa

import (
	"api-caixa/database/query"
	"api-caixa/logar"
	"api-caixa/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Fechar(dynamoClient *dynamodb.Client, log logar.Logfile) {
	caixa := query.GetLatestCaixa(dynamoClient, log)
	pagamentos := query.GetPagamentosAfterDate(dynamoClient, log, caixa.Dia)
	caixaNovo := model.Caixa{
		Dia:              time.Now().Format("2006-01-02_15:04:05"),
		DinheiroAbertura: caixa.DinheiroFechamento,
	}
	var TotalDebito float64
	var TotalPersyCoins float64
	var TotalPicPay float64
	var TotalPix float64
	var TotalCredito float64
	TotalDinheiro := caixa.DinheiroAbertura
	for _, pagamento := range pagamentos {
		TotalDebito += pagamento.Debito
		TotalPersyCoins += pagamento.PersyCoins
		TotalPicPay += pagamento.PicPay
		TotalPix += pagamento.Pix
		TotalCredito += pagamento.Credito
		TotalDinheiro += pagamento.Dinheiro
	}
	caixaNovo.TotalDebito = TotalDebito
	caixaNovo.TotalPersyCoins = TotalPersyCoins
	caixaNovo.TotalPicPay = TotalPicPay
	caixaNovo.TotalPix = TotalPix
	caixaNovo.TotalCredito = TotalCredito
	caixaNovo.DinheiroFechamento = TotalDinheiro
	query.InsertCaixa(dynamoClient, log, caixaNovo)
}
