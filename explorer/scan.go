package explorer

import (
	"context"
	"errors"
	"fmt"

	"sync"
	"time"

	"github.com/HcashOrg/hcd/chaincfg/chainhash"
	"github.com/dogecoinw/go-dogecoin/log"
	"github.com/james-ray/unielon-indexer/config"
	rpcclient "github.com/james-ray/unielon-indexer/package/github.com/HcashOrg/hcrpcclient"
	"github.com/james-ray/unielon-indexer/storage"
	"github.com/james-ray/unielon-indexer/verifys"
)

const (
	startInterval = 1 * time.Second

	wdogeFeeAddress  = "HsSDD4XmJ6YDoUMHqX3db7Wur1jEyfRExX5"
	wdogeCoolAddress = "HsLfRbdKLBaCmJUTN5YysnSasxvaKEudxPG"
)

var (
	chainNetworkErr = errors.New("Chain network error")
)

type Explorer struct {
	config     *config.Config
	node       *rpcclient.Client
	dbc        *storage.DBClient
	verify     *verifys.Verifys
	fromBlock  int64
	feeAddress string

	ctx context.Context
	wg  *sync.WaitGroup
}

func NewExplorer(ctx context.Context, wg *sync.WaitGroup, rpcClient *rpcclient.Client, dbc *storage.DBClient, fromBlock int64, feeAddress string) *Explorer {
	exp := &Explorer{
		node:       rpcClient,
		dbc:        dbc,
		verify:     verifys.NewVerifys(dbc, feeAddress),
		fromBlock:  fromBlock,
		ctx:        ctx,
		wg:         wg,
		feeAddress: feeAddress,
	}
	return exp
}

func (e *Explorer) Start() {

	defer e.wg.Done()

	if e.fromBlock == 0 {
		forkBlockHash, err := e.dbc.LastBlock()
		if err != nil {
			e.fromBlock = 0
		} else {
			e.fromBlock = forkBlockHash
		}
	}

	startTicker := time.NewTicker(startInterval)
out:
	for {
		select {
		case <-startTicker.C:
			if err := e.scanTxHash("812f52cb15cf4c395647a153680073085e3b12c02649c7fd2622860f2883e622"); err != nil {
				log.Error("explorer", "Start", err.Error())
			}
		case <-e.ctx.Done():
			log.Warn("explorer", "Stop", "Done")
			break out
		}
	}
}
func (e *Explorer) scanTxHash(txHashString string) error {
	txHash, _ := chainhash.NewHashFromStr(txHashString)
	// rawTx, _ := e.node.GetRawTransaction(txHash)
	transactionVerbose, err := e.node.GetRawTransactionVerbose(txHash)
	if err != nil {
		log.Error("scanning", "GetRawTransactionVerboseBool", err, "txhash", transactionVerbose.Txid)
		return err
	}

	decode, pushedData, err := e.reDecode(transactionVerbose)
	if err != nil {
		log.Trace("scanning", "verifyReDecode", err, "txhash", txHashString)
	}
	fmt.Println("decode", decode)

	if decode.P == "drc-20" {
		card, err := e.drc20Decode(transactionVerbose, pushedData, e.fromBlock)
		if err != nil {
			log.Error("scanning", "drc20Decode", err, "txhash", transactionVerbose.Txid)
		}

		err = e.verify.VerifyDrc20(card)
		if err != nil {
			log.Error("scanning", "VerifyDrc20", err, "txhash", transactionVerbose.Txid)
			e.dbc.UpdateCardinalsInfoNewErrInfo(card.OrderId, err.Error())

		}

		err = e.deployOrMintOrTransfer(card)
		if err != nil {
			log.Error("scanning", "deployOrMintOrTransfer", err, "txhash", transactionVerbose.Txid)
			e.dbc.UpdateCardinalsInfoNewErrInfo(card.OrderId, err.Error())

		}
	} else if decode.P == "pair-v1" {
		swap, err := e.swapDecode(transactionVerbose, pushedData, e.fromBlock)
		if err != nil {
			log.Error("scanning", "swapDecode", err, "txhash", transactionVerbose.Txid)

		}

		err = e.verify.VerifySwap(swap)
		if err != nil {
			log.Error("scanning", "VerifySwap", err, "txhash", transactionVerbose.Txid)
			e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

		}

		if swap.Op == "create" || swap.Op == "add" {
			err = e.swapCreateOrAdd(swap)
			if err != nil {
				log.Error("scanning", "swapCreateOrAdd", err, "txhash", transactionVerbose.Txid)
				e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

			}
		}

		if swap.Op == "remove" {
			err = e.swapRemove(swap)
			if err != nil {
				log.Error("scanning", "swapRemove", err, "txhash", transactionVerbose.Txid)
				e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

			}
		}

		if swap.Op == "swap" {
			if err = e.swapNow(swap); err != nil {
				log.Error("scanning", "swapNow", err, "txhash", transactionVerbose.Txid)
				e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

			}
		}

	} else if decode.P == "wdoge" {
		wdoge, err := e.wdogeDecode(transactionVerbose, pushedData, e.fromBlock)
		if err != nil {
			log.Error("scanning", "wdogeDecode", err, "txhash", transactionVerbose.Txid)

		}

		err = e.verify.VerifyWDoge(wdoge)
		if err != nil {
			log.Error("scanning", "VerifyWDoge", err, "txhash", transactionVerbose.Txid)
			e.dbc.UpdateWDogeInfoErr(wdoge.OrderId, err.Error())

		}

		if wdoge.Op == "deposit" {
			if err = e.dogeDeposit(wdoge); err != nil {
				log.Error("scanning", "dogeDeposit", err.Error(), "txhash", transactionVerbose.Txid)
				e.dbc.UpdateWDogeInfoErr(wdoge.OrderId, err.Error())

			}
		}

		if wdoge.Op == "withdraw" {
			if err = e.dogeWithdraw(wdoge); err != nil {
				log.Error("scanning", "dogeWithdraw", err.Error(), "txhash", transactionVerbose.Txid)
				e.dbc.UpdateWDogeInfoErr(wdoge.OrderId, err.Error())

			}
		}
	}
	return nil
}
func (e *Explorer) scan() error {

	blockCount, err := e.node.GetBlockCount()
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("scan GetBlockCount err: %s", err.Error())
	}

	if blockCount-e.fromBlock > 10 {
		blockCount = e.fromBlock + 10
	}

	for ; e.fromBlock < blockCount; e.fromBlock++ {
		err := e.forkBack()
		if err != nil {
			return fmt.Errorf("scan forkBack err: %s", err.Error())
		}

		log.Info("explorer", "scanning start ", e.fromBlock)
		blockHash, err := e.node.GetBlockHash(e.fromBlock)
		if err != nil {
			return fmt.Errorf("scan GetBlockHash err: %s", err.Error())
		}

		block, err := e.node.GetBlock(blockHash)
		if err != nil {
			return fmt.Errorf("scan GetBlockVerboseBool err: %s", err.Error())
		}

		for _, tx := range block.Transactions {
			txHash := tx.TxHash()
			fmt.Printf("handle tx %s \n", txHash.String())
			transactionVerbose, err := e.node.GetRawTransactionVerbose(&txHash)
			if err != nil {
				log.Error("scanning", "GetRawTransactionVerboseBool", err, "txhash", transactionVerbose.Txid)
				return err
			}

			decode, pushedData, err := e.reDecode1(tx)
			if err != nil {
				log.Trace("scanning", "verifyReDecode", err, "txhash", tx.TxHash())

			}

			if decode.P == "drc-20" {
				card, err := e.drc20Decode(transactionVerbose, pushedData, e.fromBlock)
				if err != nil {
					log.Error("scanning", "drc20Decode", err, "txhash", transactionVerbose.Txid)

				}

				err = e.verify.VerifyDrc20(card)
				if err != nil {
					log.Error("scanning", "VerifyDrc20", err, "txhash", transactionVerbose.Txid)
					e.dbc.UpdateCardinalsInfoNewErrInfo(card.OrderId, err.Error())

				}

				err = e.deployOrMintOrTransfer(card)
				if err != nil {
					log.Error("scanning", "deployOrMintOrTransfer", err, "txhash", transactionVerbose.Txid)
					e.dbc.UpdateCardinalsInfoNewErrInfo(card.OrderId, err.Error())

				}
			} else if decode.P == "pair-v1" {
				swap, err := e.swapDecode(transactionVerbose, pushedData, e.fromBlock)
				if err != nil {
					log.Error("scanning", "swapDecode", err, "txhash", transactionVerbose.Txid)

				}

				err = e.verify.VerifySwap(swap)
				if err != nil {
					log.Error("scanning", "VerifySwap", err, "txhash", transactionVerbose.Txid)
					e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

				}

				if swap.Op == "create" || swap.Op == "add" {
					err = e.swapCreateOrAdd(swap)
					if err != nil {
						log.Error("scanning", "swapCreateOrAdd", err, "txhash", transactionVerbose.Txid)
						e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

					}
				}

				if swap.Op == "remove" {
					err = e.swapRemove(swap)
					if err != nil {
						log.Error("scanning", "swapRemove", err, "txhash", transactionVerbose.Txid)
						e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

					}
				}

				if swap.Op == "swap" {
					if err = e.swapNow(swap); err != nil {
						log.Error("scanning", "swapNow", err, "txhash", transactionVerbose.Txid)
						e.dbc.UpdateSwapInfoErr(swap.OrderId, err.Error())

					}
				}

			} else if decode.P == "wdoge" {
				wdoge, err := e.wdogeDecode(transactionVerbose, pushedData, e.fromBlock)
				if err != nil {
					log.Error("scanning", "wdogeDecode", err, "txhash", transactionVerbose.Txid)

				}

				err = e.verify.VerifyWDoge(wdoge)
				if err != nil {
					log.Error("scanning", "VerifyWDoge", err, "txhash", transactionVerbose.Txid)
					e.dbc.UpdateWDogeInfoErr(wdoge.OrderId, err.Error())

				}

				if wdoge.Op == "deposit" {
					if err = e.dogeDeposit(wdoge); err != nil {
						log.Error("scanning", "dogeDeposit", err.Error(), "txhash", transactionVerbose.Txid)
						e.dbc.UpdateWDogeInfoErr(wdoge.OrderId, err.Error())

					}
				}

				if wdoge.Op == "withdraw" {
					if err = e.dogeWithdraw(wdoge); err != nil {
						log.Error("scanning", "dogeWithdraw", err.Error(), "txhash", transactionVerbose.Txid)
						e.dbc.UpdateWDogeInfoErr(wdoge.OrderId, err.Error())

					}
				}
			}
		}

		err = e.dbc.UpdateBlock(e.fromBlock, blockHash.String())
		if err != nil {
			return fmt.Errorf("scan SetBlockHash err: %s", err.Error())
		}

		log.Info("explorer", "scanning end ", e.fromBlock)
	}
	return nil
}
