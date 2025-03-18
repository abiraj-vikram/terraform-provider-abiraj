package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func setParam(params map[string]any, key string, value attr.Value) {
	if !value.IsNull() && !value.IsUnknown() {
		switch v := value.(type) {
		case types.String:
			if v.ValueString() != "" {
				params[key] = v.ValueString()
			}
		case types.Int64:
			params[key] = v.ValueInt64()
		case types.Int32:
			params[key] = v.ValueInt32()
		case types.Bool:
			params[key] = v.ValueBool()
		}
	}
}

func get_request(params map[string]any, api_url string) ([]byte, error) {
	client := &http.Client{}

	api_url = SecurdenServerURL + api_url

	api_request, err := http.NewRequest("GET", api_url, nil)
	if err != nil {
		return nil, err
	}
	api_request.Header.Set(const_authtoken, SecurdenAuthToken)

	q := api_request.URL.Query()
	for key, value := range params {
		switch v := value.(type) {
		case string:
			if v != "" {
				q.Add(key, v)
			}
		case int, int64, float64:
			q.Add(key, fmt.Sprintf("%v", v))
		case bool:
			q.Add(key, strconv.FormatBool(v))
		default:

		}
	}
	api_request.URL.RawQuery = q.Encode()
	resp, err := client.Do(api_request)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func fetchSSLCertificate(serverURL string) (*x509.Certificate, error) {
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	host := parsedURL.Host
	if parsedURL.Port() == "" {
		// Default HTTPS port if none is provided
		host += ":443"
	}

	// Dial the server using TLS
	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: true, // Just fetching the cert
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Retrieve certificates
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	return certs[0], nil
}

// Function to create an HTTP client with the fetched certificate
func createSecureClient(cert *x509.Certificate) *http.Client {
	certPool := x509.NewCertPool()
	certPool.AddCert(cert)

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func createInsecureClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Disable SSL verification
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return &http.Client{
		Transport: transport,
	}
}

func readPEMFile(filePath string) (*x509.Certificate, error) {
	pemData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM file")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

func raise_request(params map[string]any, apiURL string, method string) ([]byte, error) {
	pattern := regexp.MustCompile("^https")
	client := &http.Client{}
	if pattern.MatchString(SecurdenServerURL) {
		if len(SecurdenCertificate) == 0 {
			cert, err := fetchSSLCertificate(SecurdenServerURL)
			if err != nil {
				client = createInsecureClient()
			} else {
				client = createSecureClient(cert)
			}
		} else if filepath.IsAbs(SecurdenCertificate) {
			cert, err := readPEMFile(SecurdenCertificate)
			if err != nil {
				return nil, fmt.Errorf("Failed to read certificate: %v", err)
			}
			client = createSecureClient(cert)
		} else {

			return nil, fmt.Errorf("Please provide valid certificate path")
		}
	}
	apiURL = SecurdenServerURL + apiURL

	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	if method == http.MethodDelete {
		q := reqURL.Query()
		for key, value := range params {
			switch v := value.(type) {
			case []int64:
				for _, num := range v {
					q.Add(key, fmt.Sprintf("%d", num))
				}
			default:
				q.Add(key, fmt.Sprintf("%v", v))
			}
		}
		reqURL.RawQuery = q.Encode()
	}

	var apiRequest *http.Request

	if method == http.MethodDelete {
		apiRequest, err = http.NewRequest(method, reqURL.String(), nil)
	} else {

		if ids, exists := params["account_ids"]; exists {
			if tfIDs, ok := ids.([]types.Int64); ok {
				var int64IDs []int64
				for _, id := range tfIDs {
					if !id.IsNull() && !id.IsUnknown() {
						int64IDs = append(int64IDs, id.ValueInt64())
					}
				}
				params["account_ids"] = int64IDs
			}
		}

		requestBody, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize request body: %v", err)
		}

		apiRequest, err = http.NewRequest(method, apiURL, bytes.NewBuffer(requestBody))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set(const_authtoken, SecurdenAuthToken)
	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return body, nil
}

func get_account_dynamic(ctx context.Context, account_id, account_name, account_title, account_type, key_field string) (AccountModel, int, string) {
	var account AccountModel
	params := make(map[string]any)
	setParam(params, "account_id", types.StringValue(account_id))
	setParam(params, "account_name", types.StringValue(account_name))
	setParam(params, "account_title", types.StringValue(account_title))
	setParam(params, "account_type", types.StringValue(account_type))
	setParam(params, "key_field", types.StringValue(key_field))
	body, err := get_request(params, "/secretsmanagement/get_account")
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Error parsing response JSON: %v", err)
	}
	statusCode, ok := response["status_code"].(float64)
	if !ok {
		return account, 500, "Missing or invalid status_code in response"
	}
	if statusCode != 200 {
		if errMsg, exists := response["error"].(map[string]interface{}); exists {
			if msg, ok := errMsg["message"].(string); ok {
				return account, int(statusCode), msg
			}
		}
		if msg, ok := response["message"].(string); ok {
			return account, int(statusCode), msg
		}
		return account, int(statusCode), "Unknown error"
	}
	accountData := make(map[string]attr.Value)
	for key, value := range response {
		switch v := value.(type) {
		case string:
			accountData[key] = types.StringValue(v)
		case float64:
			accountData[key] = types.StringValue(fmt.Sprintf("%v", v))
		case map[string]interface{}:
			nestedData := make(map[string]attr.Value)
			for nestedKey, nestedValue := range v {
				if nestedStr, ok := nestedValue.(string); ok {
					nestedData[nestedKey] = types.StringValue(nestedStr)
				} else {
					nestedData[nestedKey] = types.StringValue(fmt.Sprintf("%v", nestedValue))
				}
			}
			accountData[key], _ = types.MapValue(types.StringType, nestedData)
		default:
			accountData[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
	}
	account.Account, _ = types.MapValue(types.StringType, accountData)
	return account, int(statusCode), "Success"
}

func get_account(ctx context.Context, account_id, account_name, account_title, key_field string) (AccountModel, int, string) {
	var account AccountModel
	params := make(map[string]any)
	setParam(params, "account_id", types.StringValue(account_id))
	setParam(params, "account_name", types.StringValue(account_name))
	setParam(params, "account_title", types.StringValue(account_title))
	setParam(params, "key_field", types.StringValue(key_field))

	var Response struct {
		AccountID         int64  `json:"account_id"`
		AccountName       string `json:"account_name"`
		AccountTitle      string `json:"account_title"`
		Password          string `json:"password"`
		KeyValue          string `json:"key_value"`
		PrivateKey        string `json:"private_key"`
		PuTTYPrivateKey   string `json:"putty_private_key"`
		Passphrase        string `json:"passphrase"`
		PPKPassphrase     string `json:"ppk_passphrase"`
		Address           string `json:"address"`
		ClientID          string `json:"client_id"`
		ClientSecret      string `json:"client_secret"`
		AccountAlias      string `json:"account_alias"`
		AccountFile       string `json:"account_file"`
		OracleSID         string `json:"oracle_sid"`
		OracleServiceName string `json:"oracle_service_name"`
		DefaultDatabase   string `json:"default_database"`
		Port              string `json:"port"`
		StatusCode        int    `json:"status_code"`
		Message           string `json:"message"`
		Error             struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	body, err := get_request(params, "/api/get_account_details_dict")

	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}
	json.Unmarshal(body, &Response)
	if Response.StatusCode != 200 {
		if Response.Error.Message != "" {
			return account, Response.StatusCode, Response.Error.Message
		}
		return account, Response.StatusCode, Response.Message
	}
	account.AccountID = types.Int64Value(Response.AccountID)
	account.AccountName = types.StringValue(Response.AccountName)
	account.AccountTitle = types.StringValue(Response.AccountTitle)
	account.Password = types.StringValue(Response.Password)
	account.KeyValue = types.StringValue(Response.KeyValue)
	account.PrivateKey = types.StringValue(Response.PrivateKey)
	account.PuTTYPrivateKey = types.StringValue(Response.PuTTYPrivateKey)
	account.Passphrase = types.StringValue(Response.Passphrase)
	account.PPKPassphrase = types.StringValue(Response.PPKPassphrase)
	account.Address = types.StringValue(Response.Address)
	account.ClientID = types.StringValue(Response.ClientID)
	account.ClientSecret = types.StringValue(Response.ClientSecret)
	account.AccountAlias = types.StringValue(Response.AccountAlias)
	account.AccountFile = types.StringValue(Response.AccountFile)
	account.OracleSID = types.StringValue(Response.OracleSID)
	account.OracleServiceName = types.StringValue(Response.OracleServiceName)
	account.DefaultDatabase = types.StringValue(Response.DefaultDatabase)
	account.Port = types.StringValue(Response.Port)
	return account, Response.StatusCode, Response.Message
}

func get_accounts(ctx context.Context, params map[string]any) (map[string]map[string]string, int, string) {
	var accounts_data = make(map[string]any)
	var null map[string]map[string]string

	body, err := raise_request(params, "/secretsmanagement/get_accounts", "POST")
	if err != nil {
		return null, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	err = json.Unmarshal(body, &accounts_data)
	if err != nil {
		return null, 500, fmt.Sprintf("Error parsing response: %v", err)
	}

	processedAccounts := make(map[string]map[string]string)

	for key, value := range accounts_data {
		accountMap, ok := value.(map[string]any)
		if !ok {
			continue
		}

		processedEntry := make(map[string]string)
		for k, v := range accountMap {
			if v == nil {
				processedEntry[k] = ""
			} else {
				processedEntry[k] = fmt.Sprintf("%v", v)
			}
		}

		processedAccounts[key] = processedEntry
	}

	return processedAccounts, 200, "Success"
}

func get_passwords(ctx context.Context, accountIDs []string) (types.Map, int, string) {
	var accountIDsInt64 []int64
	for _, id := range accountIDs {
		accountID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return types.Map{}, 400, fmt.Sprintf("Invalid account ID format: %v", err)
		}
		accountIDsInt64 = append(accountIDsInt64, accountID)
	}

	params := map[string]interface{}{
		"account_ids": accountIDsInt64,
	}

	body, err := raise_request(params, "/api/get_multiple_accounts_passwords", "POST")
	if err != nil {
		return types.Map{}, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	var response struct {
		Passwords  map[string]string `json:"passwords"`
		StatusCode int               `json:"status_code"`
		Message    string            `json:"message"`
		Error      struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return types.Map{}, 500, fmt.Sprintf("Failed to parse response: %v", err)
	}

	if response.StatusCode != 200 {
		errorMessage := response.Message
		if response.Error.Message != "" {
			errorMessage = response.Error.Message
		}
		return types.Map{}, response.StatusCode, errorMessage
	}

	passwordsMap := make(map[string]attr.Value, len(response.Passwords))
	for k, v := range response.Passwords {
		passwordsMap[k] = types.StringValue(v)
	}

	passwords, diags := types.MapValue(types.StringType, passwordsMap)
	if diags.HasError() {
		return types.Map{}, 500, fmt.Sprintf("Error setting map value: %v", diags)
	}

	return passwords, response.StatusCode, "Success"
}

func add_account_function(ctx context.Context, params map[string]any) (AddAccountModel, int, string) {
	var account AddAccountModel
	body, err := raise_request(params, "/api/add_account", "POST")
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	var response struct {
		ID         int64  `json:"ID"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Error      struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Failed to parse response: %v", err)
	}
	if response.StatusCode != 200 && response.StatusCode != 0 {
		errorMessage := response.Message
		if response.Error.Message != "" {
			errorMessage = response.Error.Message
		}
		return account, response.StatusCode, errorMessage
	}
	account.ID = types.Int64Value(response.ID)
	account.Message = types.StringValue(response.Message)
	return account, response.StatusCode, response.Message
}

func delete_accounts_function(ctx context.Context, params map[string]any) (DeleteAccountsModel, int, string) {
	var account DeleteAccountsModel
	body, err := raise_request(params, "/api/delete_accounts", "DELETE")
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}

	var response map[string]any
	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Error parsing response: %v", err)
	}

	if message, ok := response["message"].(string); ok {
		account.Message = types.StringValue(message)
	}

	if deletedIDs, ok := response["IDs deleted successfully"].([]any); ok {
		var idList []types.Int64
		for _, id := range deletedIDs {
			if idFloat, ok := id.(float64); ok {
				idList = append(idList, types.Int64Value(int64(idFloat)))
			}
		}
		account.DeletedAccounts = idList
	}

	return account, 200, account.Message.ValueString()
}

func edit_account_function(ctx context.Context, params map[string]any) (EditAccountModel, int, string) {
	var account EditAccountModel
	body, err := raise_request(params, "/api/edit_account", "PUT")
	if err != nil {
		return account, 500, fmt.Sprintf("Error in API call: %v", err)
	}
	var response struct {
		ID         int64  `json:"ID"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Error      struct {
			Code    interface{} `json:"code"`
			Message string      `json:"message"`
		} `json:"error"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return account, 500, fmt.Sprintf("Failed to parse response: %v", err)
	}
	if response.StatusCode != 200 && response.StatusCode != 0 {
		errorMessage := response.Message
		if response.Error.Message != "" {
			errorMessage = response.Error.Message
		}
		return account, response.StatusCode, errorMessage
	}
	account.Message = types.StringValue(response.Message)
	return account, response.StatusCode, response.Message
}
