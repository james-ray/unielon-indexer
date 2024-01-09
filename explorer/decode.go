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

func PushedData(script []byte) ([][]byte, error) {
	const scriptVersion = 0

	var data [][]byte
	tokenizer := txscript.MakeScriptTokenizer(scriptVersion, script)
	for tokenizer.Next() {
		if tokenizer.Data() != nil {
			data = append(data, tokenizer.Data())
		} else if tokenizer.Opcode() == txscript.OP_0 {
			data = append(data, nil)
		}
	}
	if err := tokenizer.Err(); err != nil {
		return nil, err
	}
	return data, nil
}
func (e *Explorer) reDecode(tx *hcjson.TxRawResult) (*utils.BaseParams, []byte, error) {

	in := tx.Vin[0]

	if in.ScriptSig == nil {
		return nil, nil, errors.New("ScriptSig is nil")
	}

	scriptbytes, err := hex.DecodeString(in.ScriptSig.Hex)
	if err != nil {
		return nil, nil, fmt.Errorf("hex.DecodeString err: %s", err.Error())
	}

	pkScript, err := PushedData(scriptbytes)
	if err != nil {
		return nil, nil, fmt.Errorf("PushedData err: %s", err.Error())
	}

	if len(pkScript) < 2 {
		return nil, nil, errors.New("pkScript length < 3")
	}

	pushedData, err := PushedData(pkScript[1])
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

	pkScript, err := PushedData(scriptbytes)
	if err != nil {
		return nil, nil, fmt.Errorf("PushedData err: %s", err.Error())
	}

	if len(pkScript) < 2 {
		return nil, nil, errors.New("pkScript length < 2")
	}

	pushedData, err := PushedData(pkScript[1])
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
