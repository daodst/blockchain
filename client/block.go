package client

import (
	"context"
	"encoding/hex"
	"freemasonry.cc/blockchain/core"
	abci "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"strings"
)

type Block struct {
	ChainId            string   //id
	Height             int64    
	Time               int64    
	LastCommitHash     string   //hash
	Datahash           string   //hash
	ValidatorsHash     string   //hash
	NextValidatorsHash string   //hash
	ConsensusHash      string   //Hash
	Apphash            string   //hash
	LastResultsHash    string   //Hash
	EvidenceHash       string   //hash
	ProposerAddress    string   
	Txs                []string //hash
	Signatures         []Signature
	LastBlockId        string //ID
	BlockId            string //ID
}

type Signature struct {
	ValidatorAddress string 
	TimeStamp        string 
	Sign             string 
}

type BlockClient struct {
	ServerUrl string
	logPrefix string
}

/*
*

*/
func (this *BlockClient) Block(height int64) (blockData *coretypes.ResultBlock, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)

	var paramsHeight *int64
	paramsHeight = &height
	if height == 0 {
		paramsHeight = nil
	}

	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	
	return node.Block(context.Background(), paramsHeight)
}

/*
*
ctxClient
*/
func (this *BlockClient) Find(height int64) (blockData *Block, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("height", height)
	blockData = &Block{}
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	
	if height == 0 {
		nodeStatus, err := node.Status(context.Background())
		if err != nil {
			log.WithError(err).Error("node.Status")
			return nil, err
		}
		height = nodeStatus.SyncInfo.LatestBlockHeight
	}

	
	blockInfo, err := node.Block(context.Background(), &height)
	if err != nil {
		log.WithError(err).Error("node.Block")
		return nil, err
	}
	blockData.Height = blockInfo.Block.Height
	blockData.Datahash = blockInfo.Block.DataHash.String()
	blockData.ChainId = blockInfo.Block.ChainID
	blockData.Time = blockInfo.Block.Time.Unix()
	blockData.Apphash = blockInfo.Block.AppHash.String()
	blockData.ConsensusHash = blockInfo.Block.ConsensusHash.String()
	blockData.EvidenceHash = blockInfo.Block.EvidenceHash.String()
	blockData.LastCommitHash = blockInfo.Block.LastCommitHash.String()
	blockData.LastResultsHash = blockInfo.Block.LastResultsHash.String()
	blockData.ValidatorsHash = blockInfo.Block.ValidatorsHash.String()
	blockData.NextValidatorsHash = blockInfo.Block.NextValidatorsHash.String()
	blockData.ProposerAddress = blockInfo.Block.ProposerAddress.String()
	blockData.LastResultsHash = blockInfo.Block.LastResultsHash.String()
	blockData.LastBlockId = blockInfo.Block.LastBlockID.Hash.String()
	blockData.BlockId = blockInfo.BlockID.Hash.String()
	for _, s := range blockInfo.Block.LastCommit.Signatures {
		signature := new(Signature)
		signature.ValidatorAddress = s.ValidatorAddress.String()
		signature.Sign = string(s.Signature)
		signature.TimeStamp = s.Timestamp.String()
		blockData.Signatures = append(blockData.Signatures, *signature)
	}
	for i := 0; i < len(blockInfo.Block.Txs); i++ {
		resTx, err := node.Tx(context.Background(), blockInfo.Block.Txs[i].Hash(), true)
		if err != nil {
			log.WithError(err).Error("node.Tx")
			return nil, err
		}
		blockData.Txs = append(blockData.Txs, strings.ToUpper(hex.EncodeToString(resTx.Hash)))
	}

	return
}


func (this *BlockClient) FindBlockResults(height *int64) (events []abci.Event, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	blockResults, err := node.BlockResults(context.Background(), height)
	if err != nil {
		log.WithError(err).Error("node.BlockResults")
		return nil, err
	}
	for _, endEvent := range blockResults.EndBlockEvents {
		blockResults.BeginBlockEvents = append(blockResults.BeginBlockEvents, endEvent)
	}
	return blockResults.BeginBlockEvents, nil
}

func (this *BlockClient) GetSyncInfo() (blockData *coretypes.SyncInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return nil, err
	}
	nodeStatus, err := node.Status(context.Background())
	if err != nil {
		log.WithError(err).Error("node.Status")
		return nil, err
	}
	return &nodeStatus.SyncInfo, nil
}


func (this *BlockClient) StatusInfo() (statusInfo *coretypes.ResultStatus, err error) {
	node, err := clientCtx.GetNode()
	return node.Status(context.Background())
}


func (this *BlockClient) NetInfo() (statusInfo *coretypes.ResultNetInfo, err error) {
	node, err := clientCtx.GetNode()
	return node.NetInfo(context.Background())
}
