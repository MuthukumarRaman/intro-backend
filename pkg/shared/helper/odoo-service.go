package helper

import (
	"fmt"

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

// Fetch demo Odoo credentials
func GetOdooDemoCredentials() (*OdooConfig, error) {
	client, err := xmlrpc.NewClient("https://demo.odoo.com/start", nil)
	if err != nil {
		return nil, err
	}

	info := map[string]string{}
	if err := client.Call("start", nil, &info); err != nil {
		return nil, err
	}

	return &OdooConfig{
		URL:      info["host"] + "/xmlrpc/2",
		DB:       info["database"],
		Username: info["user"],
		Password: info["password"],
	}, nil
}

func AuthenticateOdoo(config *OdooConfig) error {

	// /xmlrpc/2/common

	finalURL := fmt.Sprintf("%s/common", config.URL)
	fmt.Println("Connecting to:", finalURL)

	client, err := xmlrpc.NewClient(finalURL, nil)
	if err != nil {
		return err
	}

	fmt.Println("Client initialized:", client)
	fmt.Println("DB:", config.DB)
	fmt.Println("Username:", config.Username)
	fmt.Println("Password:", config.Password)

	common := map[string]any{}
	if err := client.Call("version", nil, &common); err != nil {
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

	client, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", config.URL), nil)
	if err != nil {
		return 0, err
	}
	fmt.Println("connectde")

	// var id int
	// err = client.Call("execute_kw", []interface{}{
	// 	config.DB, config.UID, config.Password,
	// 	"res.partner", "create",
	// 	[]interface{}{
	// 		map[string]interface{}{
	// 			"name":         name,
	// 			"email":        email,
	// 			"phone":        phone,
	// 			"company_type": "company",
	// 		},
	// 	},
	// }, &id)

	var id int64
	if err := client.Call("execute_kw", []any{
		config.DB, config.UID, config.Password,
		"res.partner", "create",
		[]map[string]string{
			{
				"name": name,
			},
		},
	}, &id); err != nil {
		return 0, err

	}

	return id, err
}

// Fetch all partners
func GetPartners(config *OdooConfig) ([]map[string]interface{}, error) {
	client, err := xmlrpc.NewClient(config.URL+"/object", nil)
	if err != nil {
		return nil, err
	}

	var partners []map[string]interface{}
	err = client.Call("execute_kw", []interface{}{
		config.DB, config.UID, config.Password,
		"res.partner", "search_read",
		[][]interface{}{},
		map[string]interface{}{"fields": []string{"id", "name", "email"}},
	}, &partners)

	return partners, err
}
