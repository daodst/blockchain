package common

import (
	"fmt"
	"freemasonry.cc/blockchain/client"
	"freemasonry.cc/blockchain/core"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

func DposCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dpos",
		Short: "dpos create and unjail",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(dposCreateCmd())
	return cmd
}

//dpos
func dposCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create dpos node",
		Run: func(cmd *cobra.Command, args []string) {
			identity := ""                                                                      //dpos 
			remark := ""                                                                        //dpos 
			name := cmd.Flag("name").Value.String()                                             //dpos
			website := cmd.Flag("website").Value.String()                                       //dpos
			contact := cmd.Flag("contact").Value.String()                                       //dpos
			selfDelegationStr := cmd.Flag("self-delegation").Value.String()                     
			minSelfDelegationStr := cmd.Flag("min-self-delegation").Value.String()              //， dpos
			commissionRateStr := cmd.Flag("commission-rate").Value.String()                     
			commissionMaxRateStr := cmd.Flag("commission-max-rate").Value.String()              
			commissionMaxChangeRateStr := cmd.Flag("commission-max-change-rate").Value.String() 
			if name == "" {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --name params cannot be empty")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}

			
			selfDelegationAmount, err := strconv.ParseFloat(selfDelegationStr, 64)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --self-delegation params value illegal")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}

			selfDelegation := core.NewLedgerCoin(selfDelegationAmount)

			
			minSelfDelegationAmount, err := strconv.ParseInt(minSelfDelegationStr, 10, 64)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --min-self-delegation params value illegal")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}
			minSelfDelegation := sdk.NewInt(minSelfDelegationAmount)
			totalMinSelfDelegationAmount := int64(selfDelegationAmount / 2) //50%
			if minSelfDelegationAmount < totalMinSelfDelegationAmount {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --min-self-delegation should be greater than ", totalMinSelfDelegationAmount)
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}

			
			commissionRate, err := sdk.NewDecFromStr(commissionRateStr)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --commission-rate params value illegal")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}
			
			commissionMaxRate, err := sdk.NewDecFromStr(commissionMaxRateStr)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --commission-max-rate params value illegal")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}
			
			commissionMaxChangeRate, err := sdk.NewDecFromStr(commissionMaxChangeRateStr)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | --commission-max-change-rate params value illegal")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}

			txClient := client.NewTxClient()
			accClient := client.NewAccountClient(&txClient)
			dposClient := client.NewGatewayClinet()

			fmt.Println("Please enter 12 mnemonic words separated by spaces")
			
			var mn1, mn2, mn3, mn4, mn5, mn6, mn7, mn8, mn9, mn10, mn11, mn12, mn13, mn14, mn15, mn16, mn17, mn18, mn19, mn20, mn21, mn22, mn23, mn24 string
			for {
				_, err := fmt.Scanf("%s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s", &mn1, &mn2, &mn3, &mn4, &mn5, &mn6, &mn7, &mn8, &mn9, &mn10, &mn11, &mn12, &mn13, &mn14, &mn15, &mn16, &mn17, &mn18, &mn19, &mn20, &mn21, &mn22, &mn23, &mn24)
				if err == nil {
					break
				} else {
					fmt.Println("Please enter 12 mnemonic words separated by spaces | ", err.Error())
					continue
				}
			}
			mnemonicWords := strings.Join([]string{mn1, mn2, mn3, mn4, mn5, mn6, mn7, mn8, mn9, mn10, mn11, mn12, mn13, mn14, mn15, mn16, mn17, mn18, mn19, mn20, mn21, mn22, mn23, mn24}, " ")
			wallet, err := accClient.CreateAccountFromSeed(mnemonicWords)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | mnemonic words params value illegal")
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}

			createAcc, err := sdk.AccAddressFromBech32(wallet.Address)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | AccAddressFromBech32 | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				cmd.Help()
				return
			}
			
			validatorAddress := sdk.ValAddress(createAcc).String()

			fmt.Println("Confirm parameters:")
			fmt.Println("--------------------------")
			fmt.Println("dpos node name:", name)
			fmt.Println("dpos node website:", website)
			fmt.Println("dpos node contact:", contact)
			fmt.Println("dpos manage address:", wallet.Address)
			fmt.Println("dpos operation address:", validatorAddress)
			fmt.Println("self-delegation:", selfDelegationAmount, "fm")
			fmt.Println("min-self-delegation:", minSelfDelegationAmount, "fm")
			fmt.Println("commission-rate:", commissionRate)
			fmt.Println("commission-max-rate:", commissionMaxRate)
			fmt.Println("commission-max-change-rate:", commissionMaxChangeRate)
			fmt.Println("--------------------------")
			fmt.Println("Is it correct?")
			fmt.Println("please input y send or n exit")
			
			var confimStr string
			for {
				_, err := fmt.Scanln(&confimStr)
				if err != nil || (confimStr != "y" && confimStr != "n") {
					fmt.Println("please input y send or n exit")
					continue
				}
				if confimStr == "y" || confimStr == "n" {
					break
				}
			}

			if confimStr == "n" {
				fmt.Println("exit")
				return
			}

			fmt.Println("begin send to chain")

			statusInfo, err := dposClient.StatusInfo()
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | StatusInfo | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}

			valPubkey, err := cryptocodec.FromTmPubKeyInterface(statusInfo.ValidatorInfo.PubKey)
			//bech32ValidatorPubkey, err := nodeClient.ParseBech32ValConsPubkey(base64ValPubkey)
			if err != nil {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | FromTmPubKeyInterface | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}

			commission := stakingtypes.NewCommissionRates(commissionRate, commissionMaxRate, commissionMaxChangeRate) 

			description := stakingtypes.NewDescription(name, identity, website, contact, remark)

			resp, err := txClient.RegisterValidator(
				wallet.Address, validatorAddress, valPubkey, selfDelegation,
				description, commission, minSelfDelegation, wallet.PrivateKey, 0.001)
			
			if err == nil {
				txResp := resp.Data.(core.BroadcastTxResponse)
				fmt.Println("TxHash:", txResp.TxHash)
				for {
					result, notFound, err := txClient.FindByHex(txResp.TxHash)
					if notFound {
						time.Sleep(time.Second)
						fmt.Println("Waiting for block out confirmation。。。")
						continue
					}
					if err != nil {
						fmt.Println("----------------------------------------------------------------------")
						fmt.Println("error | RegisterValidator | ", err.Error())
						fmt.Println("----------------------------------------------------------------------")
						return
					}
					if result.TxResult.Code == 0 {
						fmt.Println("----------------------------------------------------------------------")
						fmt.Println("tx height:", result.Height)
						fmt.Println("Dpos node created successfully")
						fmt.Println("----------------------------------------------------------------------")
						break
					} else {
						fmt.Println("----------------------------------------------------------------------")
						fmt.Println("Failed to create dpos node")
						fmt.Println(result.TxResult.Log)
						fmt.Println("----------------------------------------------------------------------")
						break
					}
				}

			} else {
				fmt.Println("----------------------------------------------------------------------")
				fmt.Println("error | RegisterValidator | ", err.Error())
				fmt.Println("----------------------------------------------------------------------")
				return
			}
		},
	}
	cmd.Flags().StringP("name", "", "", "Dpos Node Name")
	cmd.Flags().StringP("website", "", "", "Dpos Node Website")
	cmd.Flags().StringP("contact", "", "", "Dpos Node Contact")
	cmd.Flags().Float64("self-delegation", 0, "FSV amount committed when creating dpos node")
	cmd.Flags().Float64("min-self-delegation", 0, "The minimum pledge amount of the dpos node manager. If it is lower than this value, the point will be imprisoned")
	cmd.Flags().Float64("commission-rate", 0.1, "commission rate")
	cmd.Flags().Float64("commission-max-rate", 0.5, "Maximum percentage of commission that can be adjusted each time")
	cmd.Flags().Float64("commission-max-change-rate", 0.01, "The amount of commission that can be adjusted each time")
	return cmd
}
