package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielkhtse/supreme-adventure/common/models"
	"github.com/danielkhtse/supreme-adventure/common/response"
	"github.com/danielkhtse/supreme-adventure/common/types"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AccountClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAccountClient(baseURL string) *AccountClient {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error loading .env file")
	}

	accountServiceURL := os.Getenv("ACCOUNT_SERVICE_URL")

	return &AccountClient{
		baseURL:    accountServiceURL,
		httpClient: &http.Client{},
	}
}

func (c *AccountClient) GetAccount(accountID types.AccountID) (*models.Account, error) {
	url := fmt.Sprintf("%s/accounts/%d", c.baseURL, accountID)

	logrus.WithFields(logrus.Fields{
		"url":    url,
		"method": "GET",
	}).Debug("sending request to account service")

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch account")
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}
	defer resp.Body.Close()

	logrus.WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
	}).Debug("received response from account service")

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
		}).Error("failed to fetch account")
		return nil, fmt.Errorf("failed to fetch account, status code: %d", resp.StatusCode)
	}

	var response models.Account
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.WithError(err).Error("failed to decode account response")
		return nil, fmt.Errorf("failed to decode account response: %w", err)
	}
	account := response

	logrus.WithFields(logrus.Fields{
		"account_id": account.ID,
	}).Debug("successfully fetched account")

	return &account, nil
}

func (c *AccountClient) TransferFunds(sourceAccountID types.AccountID, destAccountID types.AccountID, amount types.AccountBalance) (err error) {
	url := fmt.Sprintf("%s/accounts/%d/balance/transfer", c.baseURL, sourceAccountID)

	requestBody := struct {
		DestAccountID types.AccountID      `json:"dest_account_id"`
		Amount        types.AccountBalance `json:"amount"`
	}{
		DestAccountID: destAccountID,
		Amount:        amount,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal request body")
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"url":             url,
		"method":          "PUT",
		"source_account":  sourceAccountID,
		"dest_account":    destAccountID,
		"transfer_amount": amount,
	}).Debug("sending transfer request to account service")

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logrus.WithError(err).Error("failed to create request")
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("failed to send transfer request")
		return fmt.Errorf("failed to send transfer request: %w", err)
	}
	defer resp.Body.Close()

	logrus.WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
	}).Debug("received response from account service")

	if resp.StatusCode != http.StatusOK {
		var response response.StandardResponse[string]
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			logrus.WithError(err).Error("failed to decode error response")
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		if response.Message != "" {
			logrus.WithFields(logrus.Fields{
				"error_message": response.Message,
			}).Error("transfer failed with error message")
			return fmt.Errorf("%s", response.Message)
		}
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
		}).Error("transfer failed")
		return fmt.Errorf("transfer failed with status code: %d", resp.StatusCode)
	}

	logrus.WithFields(logrus.Fields{
		"source_account":  sourceAccountID,
		"dest_account":    destAccountID,
		"transfer_amount": amount,
	}).Info("successfully completed transfer")

	return nil
}
