package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/kolo/xmlrpc"
)

// Odoo connection details struct
type OdooConfig struct {
	URL      string
	DB       string
	Username string
	Password string
	UID      int64
}
type OdooClient struct {
	URL      string
	DB       string
	Username string
	Password string
	UID      int
	Object   *xmlrpc.Client
}

type OodoInterface interface {
	GetSubscriptionTemplates() ([]map[string]interface{}, error)
	CreateUserSubscriptions(partnerId int64) (int64, error)
	ListModels() ([]string, error)
}

func New() *OdooClient {
	client, err := NewOdooClient()
	if err != nil {
		return nil
	}
	return client
}

// Fetch  Odoo credentials
func GetOdooConfigFromEnv() (*OdooConfig, error) {
	odooURL := os.Getenv("ODOO_URL")
	odooDB := os.Getenv("ODOO_DB_NAME")
	odooUser := os.Getenv("ODOO_USER_ID")
	odooPassword := os.Getenv("ODOO_PASSWORD")

	if odooURL == "" || odooDB == "" || odooUser == "" || odooPassword == "" {
		return nil, fmt.Errorf("missing one or more required environment variables")
	}

	return &OdooConfig{
		URL:      odooURL + "/xmlrpc/2",
		DB:       odooDB,
		Username: odooUser,
		Password: odooPassword,
	}, nil
}

func NewOdooClient() (*OdooClient, error) {

	config, err := GetOdooConfigFromEnv()
	if err != nil {
		return nil, err
	}

	// Authenticate
	common, err := xmlrpc.NewClient(config.URL+"/common", nil)
	if err != nil {
		return nil, err
	}

	var uid int
	err = common.Call("authenticate", []interface{}{config.DB, config.Username, config.Password, map[string]interface{}{}}, &uid)
	if err != nil || uid == 0 {
		return nil, err
	}

	object, err := xmlrpc.NewClient(config.URL+"/object", nil)
	if err != nil {
		return nil, err
	}

	return &OdooClient{
		URL:      config.URL,
		DB:       config.DB,
		Username: config.Username,
		Password: config.Password,
		UID:      uid,
		Object:   object,
	}, nil
}

func (client *OdooClient) GetSubscriptionTemplates() ([]map[string]interface{}, error) {
	var templates []map[string]interface{}
	err := client.Object.Call("execute_kw", []interface{}{
		client.DB, client.UID, client.Password,
		"sale.subscription.plan", "search_read",
		[]interface{}{[]interface{}{}},
		map[string]interface{}{
			"fields": []string{"id", "name"},
		},
	}, &templates)

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func (client *OdooClient) GetCompanies() ([]map[string]interface{}, error) {
	var companies []map[string]interface{}
	err := client.Object.Call("execute_kw", []interface{}{
		client.DB, client.UID, client.Password,
		"res.company", "search_read",
		[]interface{}{[]interface{}{}},
		map[string]interface{}{
			"fields": []string{"id", "name"},
		},
	}, &companies)

	if err != nil {
		return nil, err
	}

	return companies, nil
}

func AuthenticateOdoo(config *OdooConfig) error {

	// /xmlrpc/2/common

	finalURL := fmt.Sprintf("%s/common", "https://kriyatec.odoo.com/xmlrpc/2")
	fmt.Println("Connecting to:", finalURL)

	client, err := xmlrpc.NewClient(finalURL, nil)
	if err != nil {
		return err
	}

	common := map[string]any{}
	if err := client.Call("version", nil, &common); err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(common)

	var uid int64
	if err := client.Call("authenticate", []any{
		config.DB, config.Username, config.Password,
		map[string]any{},
	}, &uid); err != nil {
		return err
	}

	if uid == 0 {
		return fmt.Errorf("authentication failed")
	}

	config.UID = uid
	return nil
}

// Create a new partner in Odoo
func CreatePartner(config *OdooConfig, name, email, phone string) (int64, error) {

	client, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", "https://kriyatec.odoo.com"), nil)
	if err != nil {
		return 0, err
	}
	fmt.Println("connected", config.UID)

	partnerData := map[string]interface{}{
		"name":         name,
		"email":        email,
		"phone":        phone,
		"company_type": "individual",
	}

	var partnerId int64
	err = client.Call("execute_kw", []interface{}{
		"kriyatec", config.UID, "Siva@2100000",
		"res.partner", "create",
		[]interface{}{partnerData},
	}, &partnerId)

	return partnerId, err
}

func (client *OdooClient) CreatePartner(name string, email string, phone string) (int64, error) {

	partnerData := map[string]interface{}{
		"name":         name,
		"email":        email,
		"phone":        phone,
		"company_type": "person",
	}

	var partnerId int64
	err := client.Object.Call("execute_kw", []interface{}{
		client.DB, client.UID, client.Password,
		"res.partner", "create",
		[]interface{}{partnerData},
	}, &partnerId)

	if err != nil {
		return 0, err
	}

	return partnerId, nil
}

func (client *OdooClient) CreateUserSubscriptions(partnerId int64) (int64, error) {

	// subscriptionData := map[string]interface{}{
	// 	"name":                    "User Subscription",
	// 	"partner_id":              partnerId,
	// 	"template_id":             1,
	// 	"recurring_interval_unit": 1,
	// 	"company_id":              1,
	// }

	// Example subscription data
	subscriptionData := map[string]interface{}{
		"name":                "New Subscription",
		"partner_id":          partnerId, // Replace with a valid customer ID
		"date_order":          time.Now(),
		"partner_invoice_id":  1,
		"partner_shipping_id": 1,
		"company_id":          1,
	}

	var subscriptionId int64
	err := client.Object.Call("execute_kw", []interface{}{
		client.DB, client.UID, client.Password,
		"sale.order", "create",
		[]interface{}{subscriptionData},
	}, &subscriptionId)
	if err != nil {
		return 0, err
	}

	return subscriptionId, nil
}

// Fetch all partners
func GetPartners(config *OdooConfig) ([]map[string]interface{}, error) {
	client, err := xmlrpc.NewClient("https://kriyatec.odoo.com/xmlrpc/2/object", nil)
	if err != nil {
		return nil, err
	}

	var partners []map[string]interface{}
	err = client.Call("execute_kw", []interface{}{
		"kriyatec", config.UID, "Siva@2100000",
		"res.partner", "search_read",
		[][]interface{}{},
		map[string]interface{}{"fields": []string{"id", "name", "email"}},
	}, &partners)

	return partners, err
}

// // Create a new partner in Odoo
func CreateUserWithSubscriptions(config *OdooConfig, name, email, phone string) (int64, error) {

	client, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", "https://kriyatec.odoo.com"), nil)
	if err != nil {
		return 0, err
	}
	fmt.Println("connected", config.UID)

	subscriptionData := map[string]interface{}{
		"name":                    "User Subscription",
		"partner_id":              1,
		"template_id":             1,
		"recurring_interval_unit": "monthly",
		"company_id":              1,
	}

	var subscriptionId int64
	err = client.Call("execute_kw", []interface{}{
		config.DB, config.UID, config.Password,
		"res.partner", "create",
		[]interface{}{subscriptionData},
	}, &subscriptionId)

	return subscriptionId, err
}

// func (client *OdooClient) ListModels() ([]string, error) {
// 	var models []string
// 	err := client.Object.Call("execute_kw", []interface{}{
// 		client.DB, client.UID, client.Password,
// 		"ir.model", "search_read",
// 		[]interface{}{[]interface{}{}},
// 		map[string]interface{}{"fields": []string{"model"}},
// 	}, &models)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return models, nil
// }

func (client *OdooClient) ListModels() ([]string, error) {
	var models []map[string]interface{}

	// Call search_read on ir.model to get the model names
	err := client.Object.Call("execute_kw", []interface{}{
		client.DB, client.UID, client.Password,
		"ir.model", "search_read",
		[]interface{}{[]interface{}{}},
		map[string]interface{}{"fields": []string{"model"}},
	}, &models)

	if err != nil {
		return nil, err
	}

	// Extract the model names from the results
	var modelNames []string
	for _, model := range models {
		if modelName, ok := model["model"].(string); ok {
			modelNames = append(modelNames, modelName)
		}
	}

	return modelNames, nil
}
