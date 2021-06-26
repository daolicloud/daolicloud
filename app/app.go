package main

import (
    "errors"
    "encoding/hex"
    "os"
    "path"
    "strings"

    daolicloud "daolicloud"
    bc "daolicloud/blockchain"
    "daolicloud/utils"

    "go.dedis.ch/onet/v3/app"
    "go.dedis.ch/onet/v3/cfgpath"
    "go.dedis.ch/onet/v3/log"
    "go.dedis.ch/onet/v3/network"
    "golang.org/x/xerrors"
    "gopkg.in/urfave/cli.v1"
)

const (
	// DefaultName is the name of the binary we produce and is used to create a directory
	// folder with this name
	DefaultName = "daolicloud"
)

func main() {
    cliApp := cli.NewApp()
    cliApp.Name = "app"
    cliApp.Usage = "Used for building other apps."
    cliApp.Version = "0.1"
    groupsDef := "the group-definition-file"
    cliApp.Commands = []cli.Command{
        {
            Name:	"time",
	    Usage:	"measure the time to contact all nodes.",
	    Aliases:	[]string{"t"},
	    ArgsUsage:	groupsDef,
	    Action:	cmdTime,
	},
        {
            Name:	"counter",
	    Usage:	"return the counter",
	    Aliases:	[]string{"c"},
	    ArgsUsage:	groupsDef,
	    Action:	cmdCounter,
	},
    }
    cliApp.Commands = append(cliApp.Commands, cmds...)
    cliApp.Flags = []cli.Flag{
        cli.IntFlag{
            Name: "debug, d",
	    Value: 0,
	    Usage: "debug-level: 1 for terse, 5 for maximal,",
	},
	cli.StringFlag{
            Name: "config, c",
            Value: path.Join(cfgpath.GetConfigPath(DefaultName), app.DefaultGroupFile),
            Usage: "Configuration file of the server",
	},
    }
    cliApp.Before = func(c *cli.Context) error {
        log.SetDebugVisible(c.Int("debug"))
	    return nil
    }
    log.ErrFatal(cliApp.Run(os.Args))
}

// Use readGroup instead later.
func parseConfig(c *cli.Context) *app.Group {
    config := c.GlobalString("config")
    if _, err := os.Stat(config); os.IsNotExist(err) {
        log.Fatalf("[-] Configuration file does not exist. %s", config)
    }
    f, err := os.Open(config)
    log.ErrFatal(err, "Couldn't open group definition file")
    group, err := app.ReadGroupDescToml(f)
    log.ErrFatal(err, "Error while reading group definition file", err)
    if len(group.Roster.List) == 0 {
        log.ErrFatalf(err, "Empty entity or invalid group definition in: %s", config)
    }
    return group
}

// Add new peer node to peer storage
func addPeer(c *cli.Context) error {
    if c.NArg() < 1 {
	    return xerrors.New("please give the following arguments: " +
	        "host[:port] [host[:port]]...")
    }
    var peers []*network.ServerIdentity
    for i := 0; i < c.NArg(); i++ {
	peerURL := c.Args().Get(i)
	si, err := utils.ConvertPeerURL(peerURL)
	if err == nil {
            peers = append(peers, si)
        }
    }
    log.Info("Add peers:", peers)
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.Peer(group.Roster, &daolicloud.Peer{
            Command: "add",
	    PeerNodes: peers,
    })
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Info("Finished to add peers", resp.List)
    return nil
}

// Delete new peer node from peer storage
func delPeer(c *cli.Context) error {
    if c.NArg() < 1 {
	    return xerrors.New("please give the following arguments: " +
	        "host[:port] [host[:port]]...")
    }
    var peers []*network.ServerIdentity
    for i := 0; i < c.NArg(); i++ {
	peerURL := c.Args().Get(i)
	si, err := utils.ConvertPeerURL(peerURL)
	if err == nil {
            peers = append(peers, si)
        }
    }
    log.Info("Remove peers: ", peers)
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.Peer(group.Roster, &daolicloud.Peer{
            Command: "del",
	    PeerNodes: peers,
    })
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Info("Finished to remove peers", resp.List)
    return nil
}

// Show peer status
func showPeer(c *cli.Context) error {
    var peers []*network.ServerIdentity
    for i := 0; i < c.NArg(); i++ {
	peerURL := c.Args().Get(i)
	si, err := utils.ConvertPeerURL(peerURL)
	if err == nil {
            peers = append(peers, si)
        }
    }
    if len(peers) > 0 {
        log.Info("Show peers: ", peers)
    } else {
        log.Info("Show all peers")
    }
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.Peer(group.Roster, &daolicloud.Peer{
            Command: "show",
	    PeerNodes: peers,
    })
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Info("\t", resp.List)
    return nil
}

// Add new proxy addresses to proxy storage
func addProxy(c *cli.Context) error {
    if c.NArg() < 1 {
	    return xerrors.New("please give the following arguments: " +
	        "host[:port] [host[:port]]...")
    }
    var peers []string
    for i := 0; i < c.NArg(); i++ {
	    peerURL := c.Args().Get(i)
        if strings.Index(peerURL, "://") < 0 {
            peerURL = "tls://" + peerURL
        }
	    /*si, err := utils.ConvertPeerURL(peerURL)
	    if err == nil {
            peers = append(peers, si)
        } else {
            log.Warn(err.Error())
        }*/
	    peers = append(peers, peerURL)
    }
    log.Info("Add peers:", peers)
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.Peer(group.Roster, &daolicloud.Proxy{
        Command: "add",
	    PeerNodes: peers,
    })
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Info("Finished to add peers", resp.List)
    return nil
}

// Delete new peer node from peer storage
func delProxy(c *cli.Context) error {
    if c.NArg() < 1 {
	    return xerrors.New("please give the following arguments: " +
	        "host[:port] [host[:port]]...")
    }
    var peers []*string
    for i := 0; i < c.NArg(); i++ {
	    peerURL := c.Args().Get(i)
        if strings.Index(peerURL, "://") < 0 {
            peerURL = "tls://" + peerURL
        }
	    /*si, err := utils.ConvertPeerURL(peerURL)
	    if err == nil {
                peers = append(peers, si)
        }*/
        peers = append(peers, peerURL)
    }
    log.Info("Remove peers: ", peers)
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.Peer(group.Roster, &daolicloud.Proxy{
        Command: "del",
	    PeerNodes: peers,
    })
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Info("Finished to remove peers", resp.List)
    return nil
}

// Show peer status
func showProxy(c *cli.Context) error {
    var peers []*string
    /*for i := 0; i < c.NArg(); i++ {
	    peerURL := c.Args().Get(i)
	    si, err := utils.ConvertPeerURL(peerURL)
	    if err == nil {
            peers = append(peers, si)
        }
    }
    if len(peers) > 0 {
        log.Info("Show peers: ", peers)
    } else {
        log.Info("Show all peers")
    }*/
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.Peer(group.Roster, &daolicloud.Proxy{
        Command: "show",
	    PeerNodes: peers,
    })
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Info("\t", resp.List)
    return nil
}

// Create the genesis block to start new blockchain.
func createGenesisBlock(c *cli.Context) error {
    group := parseConfig(c)
    client := daolicloud.NewClient()
    resp, err := client.CreateGenesisBlock(group.Roster)
    if err != nil {
	    return errors.New("Error: " + err.Error())
    }
    log.Infof("Create new blockchain successful.\n\tIndex: %d\n\tHash: %#x", resp.Index, resp.Hash)
    return nil

}

// Returns the time needed to contact all nodes.
func cmdTime(c *cli.Context) error {
    log.Info("Time command")
    group := readGroup(c)
    client := daolicloud.NewClient()
    resp, err := client.Clock(group.Roster)
    if err != nil {
	    return errors.New("When asking the time: " + err.Error())
    }
    log.Infof("Children: %d - Time spent: %f", resp.Children, resp.Time)
    return nil
}

func getIDPointer(s string) (*bc.BlockID, error) {
    if strings.HasPrefix(strings.ToLower(s), "0x") {
        s = s[2:]
    }
    if s == "" {
        return nil, nil
    }
    bHash, err := hex.DecodeString(s)
    if err != nil {
        return nil, xerrors.Errorf("couldn't decode %s: %+v", s, err)
    }
    blockID := bc.BlockID(bHash)
    return &blockID, nil
}

func showBlock(c *cli.Context) error {
    blockID, err := getIDPointer(c.String("hash"))
    if err != nil {
        return xerrors.Errorf("couldn't get hash: %+v", err)
    }
    blockIndex := c.Int("index")
    if blockIndex >= 0 && blockID != nil {
        return xerrors.New("--index or --hash don't same time")
    }

    group := parseConfig(c)
    client := daolicloud.NewClient()

    var resp *bc.Block
    var resperr error
    if blockID != nil {
        resp, resperr = client.GetBlockByID(group.Roster, *blockID)
    } else if blockIndex >= 0 {
        resp, resperr = client.GetBlockByIndex(group.Roster, blockIndex)
    } else {
        resp, resperr = client.GetLatestBlock(group.Roster)
    }
    if resperr != nil {
        return xerrors.Errorf("couldn't get block: %+v", resperr.Error())
    }
    //log.Infof("\tIndex=%d\n\tHash=%#x", resp.Index, resp.Hash)
    log.Infof("%s", resp.String())
    return nil
}

// Returns the number of calls
func cmdCounter(c *cli.Context) error {
    log.Info("Counter command")
    group := readGroup(c)
    client := daolicloud.NewClient()
    counter, err := client.Count(group.Roster.RandomServerIdentity())
    if err != nil {
        return errors.New("When asking for counter: " + err.Error())
    }
    log.Info("Number of requests:", counter)
    return nil
}

func readGroup(c *cli.Context) *app.Group {
    if c.NArg() != 1 {
        log.Fatal("Please give the group-file as argument")
    }
    name := c.Args().First()
    f, err := os.Open(name)
    log.ErrFatal(err, "Couldn't open group definition file")
    group, err := app.ReadGroupDescToml(f)
    log.ErrFatal(err, "Error while reading group definition file", err)
    if len(group.Roster.List) == 0 {
        log.ErrFatalf(err, "Empty entity or invalid group definition in: %s", name)
    }
    return group
}
