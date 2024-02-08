package query

import (
	"api-caixa/logger"
	"api-caixa/model"
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetCaixaByDate retorna o caixa de uma data específica
func GetLatestCaixa(seqVal int, dynamoClient *dynamodb.Client, log *logger.Logrus) model.Caixa {
	// Cria expressão de busca com base no numero do caixa
	query := expression.Name("Seq").Equal(expression.Value(seqVal))

	// Cosntroi uma expressão
	expr, err := expression.NewBuilder().WithFilter(query).Build()

	//Checa para erros
	log.Check(err)

	// Prepara os parametros para o ScanInput do dynamoDB
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("Caixa"),
	}

	//Faz a chamada na API do DynamoDB
	result, err := dynamoClient.Scan(context.TODO(), params)
	log.Check(err)

	var caixa model.Caixa

	err = attributevalue.UnmarshalMap(result.Items[0], &caixa)

	log.Check(err)

	return caixa

}

// GetCaixaByDate retorna o caixa de uma data específica
func GetCaixaAtual(seqVal int, dynamoClient *dynamodb.Client, log *logger.Logrus) []model.Pagamento {
	// Cria expressão de busca com base no numero do caixa
	query := expression.Name("Seq").Equal(expression.Value(seqVal + 1))

	// Cosntroi uma expressão
	expr, err := expression.NewBuilder().WithFilter(query).Build()

	//Checa para erros
	log.Check(err)

	// Prepara os parametros para o ScanInput do dynamoDB
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("Pagamentos"),
	}

	//Faz a chamada na API do DynamoDB
	result, err := dynamoClient.Scan(context.TODO(), params)
	log.Check(err)

	var pagamentos []model.Pagamento

	for _, item := range result.Items {
		pagamento := model.Pagamento{}
		err := attributevalue.UnmarshalMap(item, &pagamento)
		log.Check(err)
		pagamentos = append(pagamentos, pagamento)
	}

	log.Check(err)

	return pagamentos

}

// GetCaixaByDate retorna o caixa de uma data específica
func GetPagamentosAfterDate(dynamoClient *dynamodb.Client, log *logger.Logrus, periodo string) []model.Pagamento {
	params := &dynamodb.ScanInput{
		TableName:        aws.String("Pagamentos"),
		FilterExpression: aws.String("#data > :periodo"),
		ExpressionAttributeNames: map[string]string{
			"#data": "Data",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":periodo": &types.AttributeValueMemberS{Value: periodo},
		},
		ProjectionExpression: aws.String("Cliente, Troco, Credito, Debito, Dinheiro, PicPay, Pix, Persycoins, #data, Pedidos"),
	}
	paginator := dynamodb.NewScanPaginator(dynamoClient, params)

	var pagamentos []model.Pagamento
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		log.Check(err)

		for _, item := range output.Items {
			pagamento := model.Pagamento{}
			err := attributevalue.UnmarshalMap(item, &pagamento)
			log.Check(err)
			pagamentos = append(pagamentos, pagamento)
		}
	}
	return pagamentos
}

// GetCaixaByDate retorna o caixa de uma data específica
func GetPagamentosbySeq(dynamoClient *dynamodb.Client, log *logger.Logrus, seq int) []model.Pagamento {
	query := expression.Name("Seq").Equal(expression.Value(seq))

	// Cosntroi uma expressão
	expr, err := expression.NewBuilder().WithFilter(query).Build()

	//Checa para erros
	log.Check(err)

	// Prepara os parametros para o ScanInput do dynamoDB
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("Caixa"),
	}

	paginator := dynamodb.NewScanPaginator(dynamoClient, params)

	var pagamentos []model.Pagamento

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		log.Check(err)

		for _, item := range output.Items {
			pagamento := model.Pagamento{}
			err := attributevalue.UnmarshalMap(item, &pagamento)
			log.Check(err)
			pagamentos = append(pagamentos, pagamento)
		}
	}
	return pagamentos
}

// InsertCaixa insere um novo caixa no banco de dados
func InsertCaixa(dynamoClient *dynamodb.Client, log *logger.Logrus, caixa model.Caixa) {
	av, err := attributevalue.MarshalMap(caixa)
	log.Check(err)

	params := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Caixa"),
	}
	_, err = dynamoClient.PutItem(context.Background(), params)
	log.Check(err)
}

func GetCaixaSeq(dynamoClient *dynamodb.Client, log *logger.Logrus) (int, error) {

	// Set up the scan input
	input := &dynamodb.ScanInput{
		TableName: aws.String("CaixaSeq"),
		Limit:     aws.Int32(1),
	}

	// Perform the scan
	result, err := dynamoClient.Scan(context.Background(), input)
	if err != nil {
		return 0, fmt.Errorf("failed to perform scan: %v", err)
	}

	if len(result.Items) == 0 {
		return 0, fmt.Errorf("no items found in table")
	}

	// Assuming the attribute storing the number is named "NumberValue"
	numberValue := result.Items[0]["Seq"]
	if numberValue == nil {
		return 0, fmt.Errorf("item does not have a 'Seq' attribute")
	}

	// Convert the value to int
	numberStr, ok := numberValue.(*types.AttributeValueMemberN)
	if !ok {
		return 0, fmt.Errorf("'Seq' attribute is not a number")
	}

	number, err := strconv.Atoi(numberStr.Value)
	if err != nil {
		return 0, fmt.Errorf("failed to convert 'Seq' to int: %v", err)
	}

	return number - 1, nil
}
