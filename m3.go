package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ogier/pflag"

	"github.com/awishformore/m3/adaptor/logger"
	"github.com/awishformore/m3/contract"
)

const (
	version string = "0.1.0"
)

func main() {

	// define configuration parameters
	testnet := pflag.BoolP("testnet", "t", true, "use testnet network")
	level := pflag.StringP("level", "l", "INFO", "log level")
	market := pflag.StringP("market", "m", "0x5661e7bc2403c7cc08df539e4a8e2972ec256d11", "Maker Market contract address")
	pflag.Parse()

	// initialize logger
	lvl, err := logger.ParseLevel(*level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAILED TO INITIALIZE LOGGER (%v)\n", err)
		os.Exit(1)
	}
	lgr := logger.New(lvl)

	lgr.Infof("starting m3 daemon version %v", version)

	// get current user
	usr, err := user.Current()
	if err != nil {
		lgr.Criticalf("could not identify current user (%v)", err)
		os.Exit(1)
	}

	// build path for IPC socket with geth
	var net string
	if *testnet {
		net = "testnet"
	} else {
		net = "real"
	}
	dir := path.Join(usr.HomeDir, ".ethereum", net, "geth.ipc")

	// initialize connection to the geth node
	conn, err := rpc.NewIPCClient(dir)
	if err != nil {
		lgr.Criticalf("could not initialize IPC connection (%v)", err)
		os.Exit(1)
	}
	defer conn.Close()
	be := backends.NewRPCBackend(conn)

	// bind maker market contract
	otc, err := contract.NewToken(common.HexToAddress(*market), be)
	if err != nil {
		lgr.Criticalf("could not bind to market contract (%v)", err)
		os.Exit(1)
	}

	_ = otc

	lgr.Infof("shutting down m3 daemon")

	os.Exit(0)
}
