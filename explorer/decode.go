package explorer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HcashOrg/hcd/hcjson"
	"github.com/HcashOrg/hcd/txscript"
	"github.com/HcashOrg/hcd/wire"
	"github.com/james-ray/unielon-indexer/utils"
)

func (e *Explorer) reDecode(tx *hcjson.TxRawResult) (*utils.BaseParams, []byte, error) {

	in := tx.Vin[0]

	if in.ScriptSig == nil {
		return nil, nil, errors.New("ScriptSig is nil")
	}

	scriptbytes, err := hex.DecodeString(in.ScriptSig.Hex)
	if err != nil {
		return nil, nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}

	pkScript, err := txscript.PushedData(scriptbytes)
	if err != nil {
		return nil, nil, fmt.Errorf("PushedData err: %s", err.Error())
	}

	if len(pkScript) < 3 {
		return nil, nil, errors.New("pkScript length < 3")
	}

	pushedData, err := txscript.PushedData(pkScript[2])
	if err != nil {
		return nil, nil, fmt.Errorf("PushedData err: %s", err.Error())
	}

	if len(pushedData) < 4 {
		return nil, nil, errors.New("len(pushedData) < 4")
	}

	param := &utils.BaseParams{}
	err = json.Unmarshal(pushedData[3], param)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Unmarshal err: %s", err.Error())
	}

	return param, pushedData[3], nil

}

/*func (e *Explorer) convertMsgTxToTxRawResult(tx *wire.MsgTx) *hcjson.TxRawResult {
	txRawResult := hcjson.TxRawResult{}
	txRawResult.Vin = tx.TxIn
}*/

func (e *Explorer) reDecode1(tx *wire.MsgTx) (*utils.BaseParams, []byte, error) {

	in := tx.TxIn[0]

	if in.SignatureScript == nil {
		return nil, nil, errors.New("ScriptSig is nil")
	}

	scriptbytes := in.SignatureScript

	pkScript, err := txscript.PushedData(scriptbytes)
	if err != nil {
		return nil, nil, fmt.Errorf("PushedData err: %s", err.Error())
	}

	if len(pkScript) < 2 {
		return nil, nil, errors.New("pkScript length < 2")
	}

	pushedData, err := txscript.PushedData(pkScript[1])
	if err != nil {
		return nil, nil, fmt.Errorf("PushedData err: %s", err.Error())
	}

	if len(pushedData) < 4 {
		return nil, nil, errors.New("len(pushedData) < 4")
	}

	param := &utils.BaseParams{}
	err = json.Unmarshal(pushedData[3], param)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Unmarshal err: %s", err.Error())
	}

	return param, pushedData[3], nil

}
