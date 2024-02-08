package caixa

import (
	"api-caixa/database/query"
	"api-caixa/logger"
	"api-caixa/model"
	"api-caixa/utils"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Fechar(dynamoClient *dynamodb.Client, log *logger.Logrus) []model.PagamentoReport {
	pagamentoReport := make([]model.PagamentoReport, 0)
	seq, err := query.GetCaixaSeq(dynamoClient, log)
	log.Check(err)
	caixa := query.GetLatestCaixa(seq, dynamoClient, log)
	pagamentos := query.GetPagamentosbySeq(dynamoClient, log, seq+1)
	caixaNovo := model.Caixa{
		Dia:              time.Now().Add(time.Hour * -3).Format("2006-01-02_15:04:05"),
		DinheiroAbertura: caixa.DinheiroFechamento,
	}
	var TotalDebito float64
	var TotalPersyCoins float64
	var TotalPicPay float64
	var TotalPix float64
	var TotalCredito float64
	TotalDinheiro := caixa.DinheiroAbertura
	for _, pagamento := range pagamentos {
		formaPagamento := make([]string, 0)
		valor := 0.0
		if pagamento.Debito > 0 {
			TotalDebito += pagamento.Debito
			valor += pagamento.Debito
			formaPagamento = append(formaPagamento, "Debito: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento += "Debito : " + utils.FormatCurrency(pagamento.Debito)
		} else if pagamento.PersyCoins > 0 {
			TotalPersyCoins += pagamento.PersyCoins
			valor += pagamento.PersyCoins
			formaPagamento = append(formaPagamento, "PersyCoins: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "PersyCoins " + utils.FormatCurrency(pagamento.PersyCoins)
		} else if pagamento.PicPay > 0 {
			TotalPicPay += pagamento.PicPay
			valor += pagamento.PicPay
			formaPagamento = append(formaPagamento, "PicPay: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "PicPay " + utils.FormatCurrency(pagamento.PicPay)
		} else if pagamento.Pix > 0 {
			TotalPix += pagamento.Pix
			valor += pagamento.Pix
			formaPagamento = append(formaPagamento, "Pix: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "Pix " + utils.FormatCurrency(pagamento.Pix)
		} else if pagamento.Credito > 0 {
			TotalCredito += pagamento.Credito
			valor += pagamento.Credito
			formaPagamento = append(formaPagamento, "Credito: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "Credito " + utils.FormatCurrency(pagamento.Credito)
		} else if pagamento.Dinheiro > 0 {
			TotalDinheiro += pagamento.Dinheiro
			valor += pagamento.Dinheiro
			formaPagamento = append(formaPagamento, "Dinheiro: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "Dinheiro " + utils.FormatCurrency(pagamento.Dinheiro)
		}
		pagamentoReport = append(pagamentoReport, model.PagamentoReport{
			Cliente:         pagamento.Cliente,
			FormasPagamento: formaPagamento,
			Valor:           pagamento.Dinheiro,
			Data:            pagamento.Data,
		})
	}
	caixaNovo.TotalDebito = TotalDebito
	caixaNovo.TotalPersyCoins = TotalPersyCoins
	caixaNovo.TotalPicPay = TotalPicPay
	caixaNovo.TotalPix = TotalPix
	caixaNovo.TotalCredito = TotalCredito
	caixaNovo.DinheiroFechamento = TotalDinheiro
	caixaNovo.PagamentoReport = pagamentoReport
	query.InsertCaixa(dynamoClient, log, caixaNovo)

	return pagamentoReport
}

func CaixaAtual(dynamoClient *dynamodb.Client, log *logger.Logrus) []model.PagamentoReport {
	pagamentoReport := make([]model.PagamentoReport, 0)
	seq, err := query.GetCaixaSeq(dynamoClient, log)
	log.Check(err)
	caixa := query.GetLatestCaixa(seq, dynamoClient, log)
	pagamentos := query.GetPagamentosbySeq(dynamoClient, log, seq+1)
	caixaNovo := model.Caixa{
		Dia:              time.Now().Add(time.Hour * -3).Format("2006-01-02_15:04:05"),
		DinheiroAbertura: caixa.DinheiroFechamento,
	}
	var TotalDebito float64
	var TotalPersyCoins float64
	var TotalPicPay float64
	var TotalPix float64
	var TotalCredito float64
	TotalDinheiro := caixa.DinheiroAbertura
	for _, pagamento := range pagamentos {
		formaPagamento := make([]string, 0)
		valor := 0.0
		if pagamento.Debito > 0 {
			TotalDebito += pagamento.Debito
			valor += pagamento.Debito
			formaPagamento = append(formaPagamento, "Debito: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento += "Debito : " + utils.FormatCurrency(pagamento.Debito)
		} else if pagamento.PersyCoins > 0 {
			TotalPersyCoins += pagamento.PersyCoins
			valor += pagamento.PersyCoins
			formaPagamento = append(formaPagamento, "PersyCoins: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "PersyCoins " + utils.FormatCurrency(pagamento.PersyCoins)
		} else if pagamento.PicPay > 0 {
			TotalPicPay += pagamento.PicPay
			valor += pagamento.PicPay
			formaPagamento = append(formaPagamento, "PicPay: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "PicPay " + utils.FormatCurrency(pagamento.PicPay)
		} else if pagamento.Pix > 0 {
			TotalPix += pagamento.Pix
			valor += pagamento.Pix
			formaPagamento = append(formaPagamento, "Pix: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "Pix " + utils.FormatCurrency(pagamento.Pix)
		} else if pagamento.Credito > 0 {
			TotalCredito += pagamento.Credito
			valor += pagamento.Credito
			formaPagamento = append(formaPagamento, "Credito: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "Credito " + utils.FormatCurrency(pagamento.Credito)
		} else if pagamento.Dinheiro > 0 {
			TotalDinheiro += pagamento.Dinheiro
			valor += pagamento.Dinheiro
			formaPagamento = append(formaPagamento, "Dinheiro: "+utils.FormatCurrency(pagamento.Debito))
			//formaPagamento = "Dinheiro " + utils.FormatCurrency(pagamento.Dinheiro)
		}
		pagamentoReport = append(pagamentoReport, model.PagamentoReport{
			Cliente:         pagamento.Cliente,
			FormasPagamento: formaPagamento,
			Valor:           pagamento.Dinheiro,
			Data:            pagamento.Data,
		})
	}
	caixaNovo.TotalDebito = TotalDebito
	caixaNovo.TotalPersyCoins = TotalPersyCoins
	caixaNovo.TotalPicPay = TotalPicPay
	caixaNovo.TotalPix = TotalPix
	caixaNovo.TotalCredito = TotalCredito
	caixaNovo.DinheiroFechamento = TotalDinheiro
	caixaNovo.PagamentoReport = pagamentoReport

	return pagamentoReport
}
