package routers

import (
	"api-caixa/caixa"
	"api-caixa/database/query"
	"api-caixa/logger"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ResponseOK retorna uma mensagem de ok
func ResponseOK(w http.ResponseWriter, log *logger.Logrus) {
	writeJSON(w, http.StatusOK, "Servidor up")

}

// writeJSON is a helper function to write a JSON response to the http.ResponseWriter.
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)

}

// GetCaixa retorna o ultimo caixa registrado no banco
func GetCaixa(w http.ResponseWriter, r *http.Request, log *logger.Logrus, dynamoClient *dynamodb.Client) {
	seq, err := query.GetCaixaSeq(dynamoClient, log)
	log.Check(err)
	caixa := query.GetLatestCaixa(seq, dynamoClient, log)
	writeJSON(w, http.StatusOK, caixa)
}

func Fechar(w http.ResponseWriter, r *http.Request, log *logger.Logrus, dynamoClient *dynamodb.Client) {
	pagamentosReport := caixa.Fechar(dynamoClient, log)
	writeJSON(w, http.StatusOK, pagamentosReport)
}

func GetCaixaAtual(w http.ResponseWriter, r *http.Request, log *logger.Logrus, dynamoClient *dynamodb.Client) {
	seq, err := query.GetCaixaSeq(dynamoClient, log)
	log.Check(err)
	caixa := query.GetCaixaAtual(seq, dynamoClient, log)
	writeJSON(w, http.StatusOK, caixa)
}
