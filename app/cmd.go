package main

import (
    "fmt"

    "gopkg.in/urfave/cli.v1"
)

var cmds = cli.Commands{
    {
        Name:	"peer",
	    Usage:	"Provides cli interface for peers.",
	    Aliases:	[]string{"p"},
	    ArgsUsage: "[add|del|show] [<arg>...]",
	    Description: fmt.Sprint(`
            app peer add HOST:PORT ...
            app peer del HOST:PORT ...
            app peer show
	    `),
	    Subcommands: cli.Commands{
            {
                Name: "add",
	            Usage:  "Add new peers",
	            ArgsUsage: "Peer Add",
	            Action: addPeer,
		    },
            {
                Name: "del",
	            Usage:  "Remove new peers",
	            ArgsUsage: "Peer Del",
	            Action: delPeer,
		    },
            {
                Name: "show",
	            Usage:  "Show all peers status",
	            ArgsUsage: "Peer Status",
	            Action: showPeer,
		    },
	    },
    },
    {
        Name:           "proxy",
        Usage:          "Public proxy addresses management.",
        Aliases:        []string{"a"},
        ArgsUsage:      "[add|del|show] [arg>...]",
	    Description:    fmt.Sprint(`
            app proxy add HOST:PORT ...
            app proxy del HOST:PORT ...
            app proxy show
	    `),
        Subcommands:    cli.Commands{
            {
                Name: "add",
	            Usage:  "Add new proxy addresses",
	            ArgsUsage: "Proxy Addresses Add",
	            Action: addProxy,
		    },
            {
                Name: "del",
	            Usage:  "Remove new proxy addresses",
	            ArgsUsage: "Proxy Addresses Del",
	            Action: delProxy,
		    },
            {
                Name: "show",
	            Usage:  "Show all proxy addresses status",
	            ArgsUsage: "Proxy Addresses Status",
	            Action: showProxy,
		    },
        },
    },
    {
        Name:	"genesisblock",
	    Usage:	"Initialize new blockchain.",
	    Aliases:	[]string{"g"},
	    ArgsUsage: "[<arg>...]",
	    Action: createGenesisBlock,
    },
    {
        Name:	"block",
	    Usage:	"Get latest block or a block given by an index or hash id",
	    Aliases:	[]string{"b"},
	    Action: showBlock,
	    Flags: []cli.Flag{
            cli.IntFlag{
                Name: "index",
		        Value: -1,
		        Usage: "give this block index",
		    },
            cli.StringFlag{
                Name: "hash",
		        Usage: "give block hash id to show",
		    },
	    },
    },
}
