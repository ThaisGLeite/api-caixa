package main

import (
	"api-caixa/database/driver"
	"api-caixa/logger"
	"api-caixa/routers"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dynamoClient *dynamodb.Client
	logs         *logger.Logrus
)

// inLambda verifica se o programa está rodando em um ambiente lambda ou não e retorna um booleano
func inLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	logs.Info("Servidor Ok")
	routers.ResponseOK(w, logs)
}

func handleGetCaixa(w http.ResponseWriter, r *http.Request) {
	routers.GetCaixa(w, r, logs, dynamoClient)
}

func handleGetCaixaAtual(w http.ResponseWriter, r *http.Request) {
	routers.GetCaixaAtual(w, r, logs, dynamoClient)
}

func handleFecharCaixa(w http.ResponseWriter, r *http.Request) {
	routers.Fechar(w, r, logs, dynamoClient)
}

// setupRouter configura as rotas da API

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/caixa/", handleGetCaixa)
	mux.HandleFunc("/caixaatual", handleGetCaixaAtual)
	mux.HandleFunc("/fechar/", handleFecharCaixa)
	return mux
}

func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// Para compilar o binario do sistema usamos:
//
//	GOARCH=arm64 GOOS=linux  CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap .
//
// para criar o zip do projeto comando:
//
//	zip lambda.zip bootstrap
//
// main.go
func main() {
	// Inicializa o logger
	logs = logger.NewGoAppTools()

	var err error
	dynamoClient, err = driver.ConfigAws()
	logs.Check(err)

	router := setupRouter()

	// Wrap your router with the SecureMiddleware
	secureMux := SecureMiddleware(router)

	if inLambda() {
		log.Fatal(gateway.ListenAndServe(":8080", secureMux))
	} else {
		logs.Info("Servidor Iniciado na porta 8080")
		log.Fatal(http.ListenAndServe(":8080", secureMux))
	}
}
