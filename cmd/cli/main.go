package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	userSvcURL   = "http://localhost:5454"
	ledgerSvcURL = "http://localhost:5455"
)

func main() {
	var rootCmd = &cobra.Command{Use: "cashapp-cli"}

	var createCmd = &cobra.Command{
		Use:   "create [tag]",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createUser(args[0])
		},
	}

	var balanceCmd = &cobra.Command{
		Use:   "balance [tag]",
		Short: "Check balance for a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkBalance(args[0])
		},
	}

	var sendCmd = &cobra.Command{
		Use:   "send [from_tag] [to_tag] [amount] [description]",
		Short: "Send money",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			fromTag := args[0]
			toTag := args[1]
			amountStr := args[2]
			desc := args[3]
			sendMoney(fromTag, toTag, amountStr, desc)
		},
	}

	var seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed test data",
		Run: func(cmd *cobra.Command, args []string) {
			seedData()
		},
	}

	rootCmd.AddCommand(createCmd, balanceCmd, sendCmd, seedCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createUser(tag string) {
	payload := map[string]string{"tag": tag}
	resp := post(userSvcURL+"/users", payload)
	fmt.Println("Response:", resp)
}

func checkBalance(tag string) {
	userID, walletID := resolveUser(tag)
	if userID == 0 {
		fmt.Printf("User %s not found\n", tag)
		return
	}

	resp := get(fmt.Sprintf("%s/wallets/%d/balance", ledgerSvcURL, walletID))
	fmt.Println("Balance:", resp)
}

func sendMoney(fromTag, toTag, amountStr, desc string) {
	fromUserID, _ := resolveUser(fromTag)
	toUserID, _ := resolveUser(toTag)

	if fromUserID == 0 || toUserID == 0 {
		fmt.Println("Could not resolve users")
		return
	}

	var amount int64
	fmt.Sscanf(amountStr, "%d", &amount)

	payload := map[string]interface{}{
		"from":        fromUserID,
		"to":          toUserID,
		"amount":      amount,
		"description": desc,
	}

	resp := post(ledgerSvcURL+"/payments", payload)
	fmt.Println("Payment Response:", resp)
}

func seedData() {
	users := []string{"alice", "bob", "charlie"}
	for _, u := range users {
		fmt.Printf("Creating user %s...\n", u)
		createUser(u)
	}

	// Wait a bit or simplly proceed
	fmt.Println("Sending 100 from alice to bob...")
	sendMoney("alice", "bob", "100", "test transfer")
}

func resolveUser(tag string) (int, int) {
	url := fmt.Sprintf("%s/users/%s", userSvcURL, tag)
	respStr := get(url)

	// Quick and dirty JSON parsing
	var resp struct {
		Meta struct {
			Data struct {
				User struct {
					ID      int `json:"id"`
					Wallets []struct {
						ID int `json:"id"`
					} `json:"wallets"`
				} `json:"user"`
			} `json:"data"`
		} `json:"meta"`
	}

	if err := json.Unmarshal([]byte(respStr), &resp); err != nil {
		log.Printf("Failed to parse user response: %v", err)
		return 0, 0
	}

	if resp.Meta.Data.User.ID == 0 {
		return 0, 0
	}

	walletID := 0
	if len(resp.Meta.Data.User.Wallets) > 0 {
		walletID = resp.Meta.Data.User.Wallets[0].ID
	}

	return resp.Meta.Data.User.ID, walletID
}

func post(url string, data interface{}) string {
	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
