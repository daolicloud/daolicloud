package utils

import (
    "net/url"

    "go.dedis.ch/kyber/v3"
    "go.dedis.ch/kyber/v3/suites"
    "go.dedis.ch/kyber/v3/util/encoding"
    "go.dedis.ch/onet/v3/network"
    "golang.org/x/xerrors"
)

// Convert url string info peer struct
func ConvertPeerURL(peerURL string) (*network.ServerIdentity, error) {
    parse, err := url.Parse(peerURL)
    if err != nil {
        return nil, xerrors.Errorf("url parse error: %v", err)
    }
    suite, err := suites.Find("Ed25519")
    if err != nil {
        return nil, xerrors.Errorf("kyber suite: %v", err)
    }
    var point kyber.Point
    if parse.User.Username() != "" {
        var err error
        point, err = encoding.StringHexToPoint(suite, parse.User.Username())
        if err != nil {
	        return nil, xerrors.Errorf("parsing public key error: %v", err)
        }
    }
    var connType network.ConnType
    switch parse.Scheme {
    case "tcp":
	    connType = network.PlainTCP
    default:
        connType = network.TLS
    }
    si := network.NewServerIdentity(point, network.NewAddress(connType, parse.Host))
    return si, nil
}
