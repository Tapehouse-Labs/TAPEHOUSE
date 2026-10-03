// SPDX-License-Identifier: MIT OR Apache-2.0

package keeper

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"math/big"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tapehouse/tapehouse/services/sdk"
)

// Sender sends the keepers' transactions from one key: each simulated first and sent only if it would succeed, one at
// a time so that their nonces follow, each awaited until mined.
type Sender struct {
	client  *sdk.Client
	chain   Chain
	key     *ecdsa.PrivateKey
	from    common.Address
	chainID *big.Int
	log     *slog.Logger
	mu      sync.Mutex
}

// NewSender returns a sender from key on chainID.
func NewSender(client *sdk.Client, chain Chain, key *ecdsa.PrivateKey, chainID *big.Int, log *slog.Logger) *Sender {
	return &Sender{client: client, chain: chain, key: key, from: crypto.PubkeyToAddress(key.PublicKey), chainID: chainID,
		log: log}
}

// From returns the address the keepers send from.
func (s *Sender) From() common.Address {
	return s.from
}

// Act simulates tx from the keeper's address and, if it would succeed, sends it and waits for its receipt. A revert,
// simulated or in the gas estimate, comes back as the *sdk.Revert sdk.DecodeRevert gives; one in the block is an error.
func (s *Sender) Act(ctx context.Context, action string, tx sdk.Tx) (*types.Receipt, error) {
	if err := s.client.Simulate(&bind.CallOpts{From: s.from, Context: ctx}, tx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	opts := bind.NewKeyedTransactor(s.key, s.chainID)
	opts.Context = ctx
	sent, err := s.client.Send(opts, tx)
	if err != nil {
		return nil, err
	}
	receipt, err := bind.WaitMined(ctx, s.chain, sent.Hash())
	if err != nil {
		return nil, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("%s reverted in block %d, transaction %s", action, receipt.BlockNumber, sent.Hash())
	}
	s.log.Info(action, "tx", sent.Hash().Hex(), "block", receipt.BlockNumber.Uint64(), "gasUsed", receipt.GasUsed)
	return receipt, nil
}

// reverted reports whether err is a revert named one of names.
func reverted(err error, names ...string) bool {
	revert, ok := sdk.DecodeRevert(err)
	return ok && slices.Contains(names, revert.Name)
}

// Keys reads the keepers' key and the halt signer's from the environment: each a Web3 Secret Storage keystore and the
// file holding its password, TAPEHOUSE_KEEPER_KEYSTORE and TAPEHOUSE_KEEPER_PASSWORD_FILE, TAPEHOUSE_HALT_KEYSTORE and
// TAPEHOUSE_HALT_PASSWORD_FILE, or a hex key for a dev node, TAPEHOUSE_KEEPER_KEY and TAPEHOUSE_HALT_KEY. A password
// is never read from the environment itself. The keepers' key is required; without the halt signer's the halt keeper
// does not run.
func Keys(getenv func(string) string) (keeper, halt *ecdsa.PrivateKey, err error) {
	if keeper, err = key(getenv, "KEEPER"); err == nil && keeper == nil {
		err = errors.New("the keepers' key is required: TAPEHOUSE_KEEPER_KEYSTORE or TAPEHOUSE_KEEPER_KEY")
	}
	if err != nil {
		return nil, nil, err
	}
	halt, err = key(getenv, "HALT")
	return keeper, halt, err
}

func key(getenv func(string) string, role string) (*ecdsa.PrivateKey, error) {
	prefix := "TAPEHOUSE_" + role + "_"
	if path := getenv(prefix + "KEYSTORE"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%sKEYSTORE: %w", prefix, err)
		}
		password, err := passwordFile(getenv, prefix+"PASSWORD_FILE")
		if err != nil {
			return nil, err
		}
		decrypted, err := keystore.DecryptKey(data, password)
		if err != nil {
			return nil, fmt.Errorf("%sKEYSTORE: %w", prefix, err)
		}
		return decrypted.PrivateKey, nil
	}
	if hex := getenv(prefix + "KEY"); hex != "" {
		parsed, err := crypto.HexToECDSA(strings.TrimPrefix(hex, "0x"))
		if err != nil {
			return nil, fmt.Errorf("%sKEY is not a private key", prefix)
		}
		return parsed, nil
	}
	return nil, nil
}

// passwordFile reads the password from the file the variable names, without its one trailing newline. Its errors name
// the variable, never the file's content.
func passwordFile(getenv func(string) string, variable string) (string, error) {
	path := getenv(variable)
	if path == "" {
		return "", fmt.Errorf("%s is required with the keystore", variable)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			err = pathErr.Err
		}
		return "", fmt.Errorf("%s could not be read: %w", variable, err)
	}
	password := strings.TrimSuffix(strings.TrimSuffix(string(data), "\n"), "\r")
	if password == "" {
		return "", fmt.Errorf("%s is empty", variable)
	}
	return password, nil
}

// opts returns call options at the latest block for ctx.
func (k *Keeper) opts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{Context: ctx, From: k.sender.From()}
}

// act sends tx for action, logging what the chain refused: names it expects are not failures.
func (k *Keeper) act(ctx context.Context, action string, tx sdk.Tx, err error, expected ...string) (bool, error) {
	if err != nil {
		return false, err
	}
	if _, err = k.sender.Act(ctx, action, tx); err == nil {
		return true, nil
	}
	if reverted(err, expected...) {
		k.log.Debug(action+" not needed", "reason", err.Error())
		return false, nil
	}
	return false, fmt.Errorf("%s: %w", action, err)
}
