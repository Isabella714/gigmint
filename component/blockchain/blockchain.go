package blockchain

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	gethRpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/Isabella714/gigmint/component"
)

var instance *BlockChain

func init() {
	instance = NewBlockChain()

	component.RegisterApplicationInitEventListener(instance)
	component.RegisterApplicationStopEventListener(instance)
}

type BlockChainConfig struct {
	ChainUri    string `mapstructure:"chain_uri"`
	ChainId     uint64 `mapstructure:"chain_id"`
	OptPriKey   string `mapstructure:"private_key"`
	GasLimit    uint64 `mapstructure:"gas_limit"`
	MaxGasPrice uint64 `mapstructure:"max_gas_price"`
}

type BlockChain struct {
	ctx          *component.ApplicationContext
	config       *BlockChainConfig
	instance     *ethclient.Client
	transactOpts *bind.TransactOpts
}

func NewBlockChain() *BlockChain {
	return &BlockChain{}
}

func (c *BlockChain) Instantiate() error {
	err := c.ctx.UnmarshalKey("blockchain", &c.config)
	if err != nil {
		return err
	}

	if c.config == nil {
		return errors.New("blockchain config isn't found")
	}

	rpcClient, err := gethRpc.DialContext(context.Background(), c.config.ChainUri)
	if err != nil {
		return errors.New("init blockchain client error")
	}
	c.instance = ethclient.NewClient(rpcClient)

	privateKey, err := gethCrypto.HexToECDSA(c.config.OptPriKey)
	if err != nil {
		return errors.New("init blockchain client error")
	}

	c.transactOpts, err = bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(c.config.ChainId))
	if err != nil {
		return err
	}
	c.transactOpts.GasLimit = c.config.GasLimit
	c.transactOpts.GasPrice = new(big.Int).SetUint64(c.config.MaxGasPrice)

	return nil
}

func Get(ctx context.Context) *ethclient.Client {
	return instance.Get(ctx)
}

func (c *BlockChain) Get(ctx context.Context) *ethclient.Client {
	return c.instance
}

func Config(ctx context.Context) *BlockChainConfig {
	return instance.Config(ctx)
}

func (c *BlockChain) Config(ctx context.Context) *BlockChainConfig {
	return c.config
}

func TransactOpts(ctx context.Context) *bind.TransactOpts {
	return instance.TransactOpts(ctx)
}

func (c *BlockChain) TransactOpts(ctx context.Context) *bind.TransactOpts {
	return c.transactOpts
}

func (c *BlockChain) Close() error {

	return nil
}

func (c *BlockChain) BeforeInit() error {
	return nil
}

func (c *BlockChain) AfterInit(applicationContext *component.ApplicationContext) error {
	c.ctx = applicationContext

	return c.Instantiate()
}

func (c *BlockChain) BeforeStop() {
	return
}

func (c *BlockChain) AfterStop() {
	_ = c.Close()

	return
}
