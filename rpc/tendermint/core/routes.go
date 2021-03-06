package core

import (
	"fmt"

	acm "github.com/eris-ltd/eris-db/account"
	"github.com/eris-ltd/eris-db/definitions"
	ctypes "github.com/eris-ltd/eris-db/rpc/tendermint/core/types"
	"github.com/eris-ltd/eris-db/txs"
	rpc "github.com/tendermint/go-rpc/server"
	rpctypes "github.com/tendermint/go-rpc/types"
)

// TODO: [ben] encapsulate Routes into a struct for a given TendermintPipe

// Magic! Should probably be configurable, but not shouldn't be so huge we
// end up DoSing ourselves.
const maxBlockLookback = 20

// TODO: eliminate redundancy between here and reading code from core/
type TendermintRoutes struct {
	tendermintPipe definitions.TendermintPipe
}

func (tmRoutes *TendermintRoutes) GetRoutes() map[string]*rpc.RPCFunc {
	var routes = map[string]*rpc.RPCFunc{
		"subscribe":               rpc.NewWSRPCFunc(tmRoutes.Subscribe, "event"),
		"unsubscribe":             rpc.NewWSRPCFunc(tmRoutes.Unsubscribe, "event"),
		"status":                  rpc.NewRPCFunc(tmRoutes.StatusResult, ""),
		"net_info":                rpc.NewRPCFunc(tmRoutes.NetInfoResult, ""),
		"genesis":                 rpc.NewRPCFunc(tmRoutes.GenesisResult, ""),
		"get_account":             rpc.NewRPCFunc(tmRoutes.GetAccountResult, "address"),
		"get_storage":             rpc.NewRPCFunc(tmRoutes.GetStorageResult, "address,key"),
		"call":                    rpc.NewRPCFunc(tmRoutes.CallResult, "fromAddress,toAddress,data"),
		"call_code":               rpc.NewRPCFunc(tmRoutes.CallCodeResult, "fromAddress,code,data"),
		"dump_storage":            rpc.NewRPCFunc(tmRoutes.DumpStorageResult, "address"),
		"list_accounts":           rpc.NewRPCFunc(tmRoutes.ListAccountsResult, ""),
		"get_name":                rpc.NewRPCFunc(tmRoutes.GetNameResult, "name"),
		"list_names":              rpc.NewRPCFunc(tmRoutes.ListNamesResult, ""),
		"broadcast_tx":            rpc.NewRPCFunc(tmRoutes.BroadcastTxResult, "tx"),
		"unsafe/gen_priv_account": rpc.NewRPCFunc(tmRoutes.GenPrivAccountResult, ""),
		"unsafe/sign_tx":          rpc.NewRPCFunc(tmRoutes.SignTxResult, "tx,privAccounts"),

		// TODO: hookup
		"blockchain": rpc.NewRPCFunc(tmRoutes.BlockchainInfo, "minHeight,maxHeight"),
		//	"get_block":               rpc.NewRPCFunc(GetBlock, "height"),
		//"list_validators":         rpc.NewRPCFunc(ListValidators, ""),
		// "dump_consensus_state":    rpc.NewRPCFunc(DumpConsensusState, ""),
		// "list_unconfirmed_txs":    rpc.NewRPCFunc(ListUnconfirmedTxs, ""),
		// subscribe/unsubscribe are reserved for websocket events.
	}
	return routes
}

func (tmRoutes *TendermintRoutes) Subscribe(wsCtx rpctypes.WSRPCContext,
	event string) (ctypes.ErisDBResult, error) {
	// NOTE: RPCResponses of subscribed events have id suffix "#event"
	result, err := tmRoutes.tendermintPipe.Subscribe(wsCtx.GetRemoteAddr(), event,
		func(result ctypes.ErisDBResult) {
			// NOTE: EventSwitch callbacks must be nonblocking
			wsCtx.TryWriteRPCResponse(
				rpctypes.NewRPCResponse(wsCtx.Request.ID+"#event", &result, ""))
		})
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (tmRoutes *TendermintRoutes) Unsubscribe(wsCtx rpctypes.WSRPCContext,
	event string) (ctypes.ErisDBResult, error) {
	result, err := tmRoutes.tendermintPipe.Unsubscribe(wsCtx.GetRemoteAddr(),
		event)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (tmRoutes *TendermintRoutes) StatusResult() (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.Status(); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) NetInfoResult() (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.NetInfo(); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) GenesisResult() (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.Genesis(); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) GetAccountResult(address []byte) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.GetAccount(address); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) GetStorageResult(address, key []byte) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.GetStorage(address, key); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) CallResult(fromAddress, toAddress,
	data []byte) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.Call(fromAddress, toAddress, data); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) CallCodeResult(fromAddress, code,
	data []byte) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.CallCode(fromAddress, code, data); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) DumpStorageResult(address []byte) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.DumpStorage(address); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) ListAccountsResult() (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.ListAccounts(); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) GetNameResult(name string) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.GetName(name); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) ListNamesResult() (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.ListNames(); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) GenPrivAccountResult() (ctypes.ErisDBResult, error) {
	//if r, err := tmRoutes.tendermintPipe.GenPrivAccount(); err != nil {
	//	return nil, err
	//} else {
	//	return r, nil
	//}
	return nil, fmt.Errorf("Unimplemented as poor practice to generate private account over unencrypted RPC")
}

func (tmRoutes *TendermintRoutes) SignTxResult(tx txs.Tx,
	privAccounts []*acm.PrivAccount) (ctypes.ErisDBResult, error) {
	// if r, err := tmRoutes.tendermintPipe.SignTx(tx, privAccounts); err != nil {
	// 	return nil, err
	// } else {
	// 	return r, nil
	// }
	return nil, fmt.Errorf("Unimplemented as poor practice to pass private account over unencrypted RPC")
}

func (tmRoutes *TendermintRoutes) BroadcastTxResult(tx txs.Tx) (ctypes.ErisDBResult, error) {
	if r, err := tmRoutes.tendermintPipe.BroadcastTxSync(tx); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (tmRoutes *TendermintRoutes) BlockchainInfo(minHeight,
	maxHeight int) (ctypes.ErisDBResult, error) {
	r, err := tmRoutes.tendermintPipe.BlockchainInfo(minHeight, maxHeight,
		maxBlockLookback)
	if err != nil {
		return nil, err
	} else {
		return r, nil
	}

}
