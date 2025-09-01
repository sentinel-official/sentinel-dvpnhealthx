package core

import (
	"errors"
	"sync"

	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/sentinel-official/sentinel-go-sdk/core"
	"github.com/sentinel-official/sentinelhub/v12/types/v1"
	"github.com/sentinel-official/sentinelhub/v12/x/node/types/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Context struct {
	authzMsgTypes    []string
	chainID          string
	client           *core.Client
	db               *mongo.Database
	feegrantMsgTypes []string
	homeDir          string
	keyringBackend   string
	maxPrice         v1.Price
	rpcAddr          string

	sealed bool

	fm  sync.RWMutex
	txm sync.Mutex
}

func NewContext() *Context {
	return &Context{
		authzMsgTypes: []string{
			cosmossdk.MsgTypeURL(&v3.MsgStartSessionRequest{}),
		},
		feegrantMsgTypes: []string{
			cosmossdk.MsgTypeURL(&authz.MsgExec{}),
		},
	}
}

// checkSealed verifies if the context is sealed to prevent modification.
func (c *Context) checkSealed() {
	if c.sealed {
		panic(errors.New("context is sealed"))
	}
}

// Seal marks the context as sealed, preventing further modifications.
func (c *Context) Seal() *Context {
	c.sealed = true
	return c
}

func (c *Context) AuthzMsgTypes() []string {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.authzMsgTypes
}

func (c *Context) ChainID() string {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.chainID
}

func (c *Context) Client() *core.Client {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.client
}

func (c *Context) Database() *mongo.Database {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.db
}

func (c *Context) FeegrantMsgTypes() []string {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.feegrantMsgTypes
}

func (c *Context) HomeDir() string {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.homeDir
}

func (c *Context) MaxPrice() v1.Price {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.maxPrice
}

func (c *Context) RPCAddr() string {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.rpcAddr
}

func (c *Context) KeyringBackend() string {
	c.fm.RLock()
	defer c.fm.RUnlock()
	return c.keyringBackend
}

func (c *Context) WithChainID(chainID string) *Context {
	c.checkSealed()
	c.chainID = chainID
	return c
}

func (c *Context) WithClient(cc *core.Client) *Context {
	c.checkSealed()
	c.client = cc
	return c
}

func (c *Context) WithDatabase(db *mongo.Database) *Context {
	c.checkSealed()
	c.db = db
	return c
}

func (c *Context) WithMaxPrice(price v1.Price) *Context {
	c.checkSealed()
	c.maxPrice = price
	return c
}

func (c *Context) WithHomeDir(dir string) *Context {
	c.checkSealed()
	c.homeDir = dir
	return c
}

func (c *Context) WithKeyringBackend(backend string) *Context {
	c.checkSealed()
	c.keyringBackend = backend
	return c
}

func (c *Context) WithRPCAddr(addr string) *Context {
	c.checkSealed()
	c.rpcAddr = addr
	return c
}
