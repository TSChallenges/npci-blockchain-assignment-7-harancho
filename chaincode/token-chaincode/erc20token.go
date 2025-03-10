package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    int    `json:"decimals"`
	TotalSupply int    `json:"totalSupply"`
}

type TokenContract struct {
	contractapi.Contract
}

type Balance struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type ApprovedAllowance struct {
	Spender string `json:"spender"`
	Owner   string `json:"owner"`
	Amount  int    `json:"amount"`
}

// TODO: InitLedger for token initialization
func (tc *TokenContract) InitToken(ctx contractapi.TransactionContextInterface) error {

	token := Token{
		Name:        "BNB-Token",
		Symbol:      "BNB",
		Decimals:    18,
		TotalSupply: 200000000,
	}

	tokenByte, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("TOKEN_INFO", tokenByte)
	if err != nil {
		return err
	}

	admin := Balance{
		Address: "User A",
		Amount:  200000000,
	}

	fmt.Println(admin)
	adminByte, err := json.Marshal(admin)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(admin.Address, adminByte)
	if err != nil {
		return err
	}

	userB := Balance{
		Address: "User B",
	}

	userBByte, err := json.Marshal(userB)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(userB.Address, userBByte)
	if err != nil {
		return err
	}

	userC := Balance{
		Address: "User C",
	}

	userCByte, err := json.Marshal(userC)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(userC.Address, userCByte)
	if err != nil {
		return err
	}

	userD := Balance{
		Address: "User D",
	}

	userDByte, err := json.Marshal(userD)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(userD.Address, userDByte)
	if err != nil {
		return err
	}
	return nil
}

// TODO: MintTokens for minting tokens
func (tc *TokenContract) MintToken(ctx contractapi.TransactionContextInterface, address string, amount int) error {
	
	if address != "User A" {
		return fmt.Errorf("Only admincan mint tokens")
	}

	existingAddressBalanceByteData, err := ctx.GetStub().GetState(address)
	if err != nil {
		return err
	}

	var existingAddressBalanceData Balance
	if existingAddressBalanceByteData == nil {
		existingAddressBalanceData.Address = address
		existingAddressBalanceData.Amount = amount
	} else {
		err := json.Unmarshal(existingAddressBalanceByteData, &existingAddressBalanceData)
		if err != nil {
			return err
		}

		existingAddressBalanceData.Amount = existingAddressBalanceData.Amount + amount
	}

	adminBalanceJSON, err := json.Marshal(existingAddressBalanceData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(existingAddressBalanceData.Address, adminBalanceJSON)
	if err != nil {
		return err
	}

	tokenByteData, err := ctx.GetStub().GetState("TOKEN_INFO")
	if err != nil {
		return err
	}

	var tokenData Token
	err = json.Unmarshal(tokenByteData, &tokenData)
	if err != nil {
		return err
	}

	tokenData.TotalSupply = tokenData.TotalSupply + amount

	tokenByteData, err = json.Marshal(tokenData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("TOKEN_INFO", tokenByteData)
	if err != nil {
		return err
	}

	return nil
}

// TODO: TransferTokens for transferring tokens
func (tc *TokenContract) TransferTokens(ctx contractapi.TransactionContextInterface, sender string, receiver string, amount int) error {

	senderByteData, err := ctx.GetStub().GetState(sender)
	if err != nil {
		return err
	}

	if senderByteData == nil {
		fmt.Println("sender address dosen't exists")
		return fmt.Errorf("sender address dosen't exists")
	}

	receiverByteData, err := ctx.GetStub().GetState(receiver)
	if err != nil {
		return err
	}

	if receiverByteData == nil {
		fmt.Println("receiver address dosen't exists")
		return fmt.Errorf("receiver address dosen't exists")
	}

	var senderData Balance
	err = json.Unmarshal(senderByteData, &senderData)
	if err != nil {
		return err
	}

	var receiverData Balance
	err = json.Unmarshal(receiverByteData, &receiverData)
	if err != nil {
		return err
	}

	if senderData.Amount < amount {
		fmt.Println("sender has insufficient balance")
		return fmt.Errorf("sender has insufficient balance")
	}

	senderData.Amount = senderData.Amount - amount
	receiverData.Amount = receiverData.Amount + amount

	senderByteData, err = json.Marshal(senderData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(sender, senderByteData)
	if err != nil {
		return err
	}

	receiverByteData, err = json.Marshal(receiverData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(receiver, receiverByteData)
	if err != nil {
		return err
	}

	return nil
}

// TODO: GetBalance to check the balance
func (tc *TokenContract) GetBalance(ctx contractapi.TransactionContextInterface, address string) (int, error) {

	var addressData Balance

	addressByteData, err := ctx.GetStub().GetState(address)
	if err != nil {
		return addressData.Amount, err
	}

	if addressByteData == nil {
		fmt.Println("address dosen't exists")
		return addressData.Amount, fmt.Errorf("address dosen't exists")
	}

	err = json.Unmarshal(addressByteData, &addressData)
	if err != nil {
		return addressData.Amount, err
	}

	fmt.Println(addressData)
	return addressData.Amount, nil
}

func (tc *TokenContract) GetTotalSupply(ctx contractapi.TransactionContextInterface) (int, error) {

	var tokenData Token

	tokenByteData, err := ctx.GetStub().GetState("TOKEN_INFO")
	if err != nil {
		return tokenData.TotalSupply, err
	}

	if tokenByteData == nil {
		fmt.Println("address dosen't exists")
		return tokenData.TotalSupply, fmt.Errorf("token dosen't exists")
	}

	err = json.Unmarshal(tokenByteData, &tokenData)
	if err != nil {
		return tokenData.TotalSupply, err
	}

	return tokenData.TotalSupply, nil
}

// TODO: ApproveSpender for approving spending
func (tc *TokenContract) ApproveSpender(ctx contractapi.TransactionContextInterface, owner string, spender string, amount int) error {

	ownerByteData, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return err
	}

	if ownerByteData == nil {
		fmt.Println("owner address dosen't exists")
		return fmt.Errorf("owner address dosen't exists")
	}

	spenderByteData, err := ctx.GetStub().GetState(spender)
	if err != nil {
		return err
	}

	if spenderByteData == nil {
		fmt.Println("spender address dosen't exists")
		return fmt.Errorf("spender address dosen't exists")
	}

	key := fmt.Sprintf("%s_%s", owner, spender)
	allowanceByteData, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}

	if allowanceByteData != nil {
		fmt.Println("allowance details for owner and spender already exists")
		return fmt.Errorf("allowance details for owner and spender already exists")
	}

	var allowance = ApprovedAllowance{
		Owner:   owner,
		Spender: spender,
		Amount:  amount,
	}

	allowanceByteData, err = json.Marshal(allowance)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(key, allowanceByteData)
	if err != nil {
		return err
	}

	return nil
}

// TODO: BurnTokens for burning tokens
func (tc *TokenContract) BurnToken(ctx contractapi.TransactionContextInterface, address string, amount int) error {
	if address != "User A" {
		return fmt.Errorf("Only admincan mint tokens")
	}

	existingAddressByteData, err := ctx.GetStub().GetState(address)
	if err != nil {
		return err
	}

	var existingAddressData Balance
	if existingAddressByteData == nil {
		fmt.Println("address desen't exists")
		return fmt.Errorf("address desen't exists")
	}

	err = json.Unmarshal(existingAddressByteData, &existingAddressData)
	if err != nil {
		return err
	}

	existingAddressData.Amount = existingAddressData.Amount - amount

	adminBalanceJSON, err := json.Marshal(existingAddressData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(existingAddressData.Address, adminBalanceJSON)
	if err != nil {
		return err
	}

	tokenByteData, err := ctx.GetStub().GetState("TOKEN_INFO")
	if err != nil {
		return err
	}

	var tokenData Token
	err = json.Unmarshal(tokenByteData, &tokenData)
	if err != nil {
		return err
	}

	tokenData.TotalSupply = tokenData.TotalSupply - amount

	tokenByteData, err = json.Marshal(tokenData)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("TOKEN_INFO", tokenByteData)
	if err != nil {
		return err
	}

	return nil
}
