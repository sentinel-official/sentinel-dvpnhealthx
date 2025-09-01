package core

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

func (c *Context) UpsertGrants(ctx context.Context, addr types.AccAddress) error {
	granterAddr, err := c.Client().MsgFromAddr()
	if err != nil {
		return fmt.Errorf("getting granter addr: %w", err)
	}

	authzMsgs, err := func() (msgs []types.Msg, err error) {
		for _, msgType := range c.AuthzMsgTypes() {
			grants, _, err := c.Client().AuthzGrants(ctx, granterAddr, addr, msgType, nil)
			if err != nil {
				return nil, fmt.Errorf("querying authz grants: %w", err)
			}

			if grants == nil {
				msg, err := authz.NewMsgGrant(granterAddr, addr, authz.NewGenericAuthorization(msgType), nil)
				if err != nil {
					return nil, fmt.Errorf("creating authz msg grant: %w", err)
				}

				msgs = append(msgs, msg)
			}
		}

		return msgs, nil
	}()

	if err != nil {
		return fmt.Errorf("generating authz msgs: %w", err)
	}

	feegrantMsgs, err := func() (msgs []types.Msg, err error) {
		grant, err := c.Client().FeegrantAllowance(ctx, granterAddr, addr)
		if err != nil {
			return nil, fmt.Errorf("querying feegrants allowance: %w", err)
		}

		if grant == nil {
			basicAllowance := &feegrant.BasicAllowance{
				SpendLimit: nil,
				Expiration: nil,
			}

			allowedMsgAllowance, err := feegrant.NewAllowedMsgAllowance(basicAllowance, c.FeegrantMsgTypes())
			if err != nil {
				return nil, fmt.Errorf("creating allowed msg allowance: %w", err)
			}

			msg, err := feegrant.NewMsgGrantAllowance(allowedMsgAllowance, granterAddr, addr)
			if err != nil {
				return nil, fmt.Errorf("creating msg grant allowance: %w", err)
			}

			msgs = append(msgs, msg)
		}

		return msgs, nil
	}()

	if err != nil {
		return fmt.Errorf("generating feegrant msgs: %w", err)
	}

	var msgs []types.Msg
	msgs = append(msgs, authzMsgs...)
	msgs = append(msgs, feegrantMsgs...)

	if err := c.BroadcastTx(ctx, msgs...); err != nil {
		return fmt.Errorf("broadcasting tx with %d grant msgs: %w", len(msgs), err)
	}

	return nil
}

func (c *Context) UpsertKey(_ context.Context, name string) (types.AccAddress, error) {
	key, err := c.Client().Key(name)
	if err != nil {
		return nil, fmt.Errorf("getting key %q: %w", name, err)
	}
	if key == nil {
		_, key, err = c.Client().CreateKey(name, "", "", "")
		if err != nil {
			return nil, fmt.Errorf("creating key %q: %w", name, err)
		}
	}

	addr, err := key.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("getting addr from key %q: %w", name, err)
	}

	return addr, nil
}

func (c *Context) Inspect(_ context.Context, addr string) ([]byte, error) {
	granterAddr, err := c.Client().MsgFromAddr()
	if err != nil {
		return nil, fmt.Errorf("getting granter addr: %w", err)
	}

	args := []string{
		"run",
		"--privileged",
		"--rm",
		"--tty",
		"--volume", fmt.Sprintf("%s:/root/.sentinel-dvpncli", c.HomeDir()),
		"sentinel-dvpncli:latest",
		"inspect",
		"--log.level", "none",
		"--keyring.backend", c.KeyringBackend(),
		"--rpc.addrs", c.RPCAddr(),
		"--rpc.chain-id", c.ChainID(),
		"--tx.authz-granter-addr", granterAddr.String(),
		"--tx.fee-granter-addr", granterAddr.String(),
		"--tx.from-name", addr,
		"--max-price", c.MaxPrice().String(),
		addr,
	}

	cmd := exec.Command("docker", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("running command: %w", err)
	}

	return output, nil
}
