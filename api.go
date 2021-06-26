package daolicloud

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
    bc "daolicloud/blockchain"

    "go.dedis.ch/cothority/v3"
    "go.dedis.ch/onet/v3"
    "go.dedis.ch/onet/v3/log"
    "go.dedis.ch/onet/v3/network"
)

// ServiceName is used for registration on the onet.
const ServiceName = "LotMintService"

// Client is a structure to communicate with daolicloud service
type Client struct {
    *onet.Client
}

// NewClient instantiates a new daolicloud.Client
func NewClient() *Client {
	return &Client{Client: onet.NewClient(cothority.Suite, ServiceName)}
}

// Clock chooses one server from the Roster at random. It
// sends a Clock to it, which is then processed on the server side
// via the code in the service package.
//
// Clock will return the time in seconds it took to run the protocol.
func (c *Client) Clock(r *onet.Roster) (*ClockReply, error) {
    dst := r.RandomServerIdentity()
    log.Lvl4("Sending message to", dst)
    reply := &ClockReply{}
    err := c.SendProtobuf(dst, &Clock{r}, reply)
    if err != nil {
        return nil, err
    }
    return reply, nil
}

// Count will return the number of times `Clock` has been called on this
// service-node.
func (c *Client) Count(si *network.ServerIdentity) (int, error) {
    reply := &CountReply{}
    err := c.SendProtobuf(si, &Count{}, reply)
    if err != nil {
        return -1, err
    }
    return reply.Count, nil
}

// Peer will operation the communication with nodes
func (c *Client) Peer(r *onet.Roster, p interface{}) (*PeerReply, error) {
    dst := r.RandomServerIdentity()
    log.Lvl4("Sending message to", dst)
    reply := &PeerReply{}
    err := c.SendProtobuf(dst, p, reply)
    if err != nil {
        return nil, err
    }
    return reply, nil
}

func (c *Client) CreateGenesisBlock(r *onet.Roster) (*bc.Block, error) {
    dst := r.RandomServerIdentity()
    log.Lvl4("Sending message to", dst)
    reply := &bc.Block{}
    err := c.SendProtobuf(dst, &GenesisBlockRequest{}, reply)
    if err != nil {
        return nil, err
    }
    return reply, nil
}

func (c *Client) GetBlockByID(r *onet.Roster, blockID bc.BlockID) (*bc.Block, error) {
    dst := r.RandomServerIdentity()
    log.Lvl4("Sending message to", dst)
    reply := &bc.Block{}
    err := c.SendProtobuf(dst, &BlockByIDRequest{blockID}, reply)
    if err != nil {
        return nil, err
    }
    return reply, nil
}

func (c *Client) GetBlockByIndex(r *onet.Roster, blockIndex uint64) (*bc.Block, error) {
    dst := r.RandomServerIdentity()
    log.Lvl4("Sending message to", dst)
    reply := &bc.Block{}
    err := c.SendProtobuf(dst, &BlockByIndexRequest{blockIndex}, reply)
    if err != nil {
        return nil, err
    }
    return reply, nil
}

func (c *Client) GetLatestBlock(r *onet.Roster) (*bc.Block, error) {
    dst := r.RandomServerIdentity()
    log.Lvl4("Sending message to", dst)
    reply := &bc.Block{}
    err := c.SendProtobuf(dst, &BlockLatestRequest{}, reply)
    if err != nil {
        return nil, err
    }
    return reply, nil
}
