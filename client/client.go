//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/client/client.go
//
package client

import (
	"github.com/elliottpolk/beam/rpc"
)

func DialAndDownload(addr, from, to string) error {
	return DialAndDownloadAt(addr, from, to, 0)
}

func DialAndDownloadAt(addr, from, to string, block int64) error {
	c, err := rpc.DialHttp(addr)
	if err != nil {
		return err
	}

	return c.DownloadAt(from, to, block)
}
