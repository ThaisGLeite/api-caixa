package query

import (
	"api-caixa/logar"
	"api-caixa/model"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetCaixaByDate retorna o caixa de uma data específica
func GetLatestCaixa(dynamoClient *dynamodb.Client, log logar.Logfile) model.Caixa {
	params := &dynamodb.QueryInput{
		TableName:              aws.String("Caixa"),
		IndexName:              aws.String("DummyIndex"),
		KeyConditionExpression: aws.String("DummyKey = :v1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: "1"},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	}
	output, err := dynamoClient.Query(context.Background(), params)
	logar.Check(err, log)

	caixa := model.Caixa{}
	err = attributevalue.UnmarshalMap(output.Items[0], &caixa)
	logar.Check(err, log)

	return caixa
}

// GetCaixaByDate retorna o caixa de uma data específica
func GetPagamentosAfterDate(dynamoClient *dynamodb.Client, log logar.Logfile, periodo string) []model.Pagamento {
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
		logar.Check(err, log)

		for _, item := range output.Items {
			pagamento := model.Pagamento{}
			err := attributevalue.UnmarshalMap(item, &pagamento)
			logar.Check(err, log)
			pagamentos = append(pagamentos, pagamento)
		}
	}
	return pagamentos
}

// InsertCaixa insere um novo caixa no banco de dados
func InsertCaixa(dynamoClient *dynamodb.Client, log logar.Logfile, caixa model.Caixa) {
	av, err := attributevalue.MarshalMap(caixa)
	logar.Check(err, log)

	params := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Caixa"),
	}
	_, err = dynamoClient.PutItem(context.Background(), params)
	logar.Check(err, log)
}
