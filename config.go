package moneymakerplaid

import (
	"context"
	"fmt"
	"github.com/plaid/plaid-go/plaid"
	"log"
	"os"
)

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

type Configuration struct {
	Client       *plaid.APIClient
	Config       *plaid.Configuration
	Products     string
	CountryCodes string
	RedirectUrl  string
}

type PlaidHandler struct {
	Configuration *Configuration
}

type Handler interface {
	GetAccountsForItem(privateToken string) (*plaid.AccountsGetResponse, error)
}

func NewConfiguration() *Configuration {

	plaidClientId := getOrExit("PLAID_CLIENT_ID")
	plaidSecret := getOrExit("PLAID_SECRET")
	plaidEnv := getOrDefault("PLAID_ENV", "sandbox")
	plaidProducts := getOrDefault("PLAID_PRODUCTS", "transactions")
	plaidCountryCodes := getOrDefault("PLAID_COUNTRY_CODES", "US")
	plaidRedirectUri := getOrDefault("PLAID_REDIRECT_URI", "")

	config := plaid.NewConfiguration()
	config.AddDefaultHeader("PLAID-CLIENT-ID", plaidClientId)
	config.AddDefaultHeader("PLAID-SECRET", plaidSecret)
	config.UseEnvironment(environments[plaidEnv])
	// config.Debug = true

	client := plaid.NewAPIClient(config)

	return &Configuration{
		Client:       client,
		Config:       config,
		Products:     plaidProducts,
		CountryCodes: plaidCountryCodes,
		RedirectUrl:  plaidRedirectUri,
	}
}

func (config *Configuration) NewPlaidHandler() Handler {
	return &PlaidHandler{
		Configuration: config,
	}
}

func (handler *PlaidHandler) GetAccountsForItem(privateToken string) (*plaid.AccountsGetResponse, error) {

	plaidClient := handler.Configuration.Client

	accountsGetRequest := *plaid.NewAccountsGetRequest(privateToken)

	ctx := context.Background()
	accountsGetResponse, _, err := plaidClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		accountsGetRequest,
	).Execute()
	if err != nil {
		log.Printf("Unable to get account details \n%+v\n", err)
		return nil, err
	}
	log.Printf("Retrieved Auth Response \n%+v\n", accountsGetResponse)

	return &accountsGetResponse, nil
}

func getOrExit(envVar string) string {
	val := os.Getenv(envVar)
	if val == "" {
		log.Fatal(fmt.Printf("%s is not set. Make sure to fill out the .env file", envVar))
	}
	return val
}

func getOrDefault(envVar string, defaultVal string) string {
	val := os.Getenv(envVar)
	if val == "" {
		return defaultVal
	}
	return val
}
