//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/client/client.go
//
package rpc

import (
	"io"
	"net/rpc"
	"os"
	"sync"

	"github.com/elliottpolk/beam/log"

	"github.com/pkg/errors"
)

const BlockSize int64 = 512 * 1024

type Client struct {
	addr string
	*rpc.Client
}

func DialHttp(addr string) (*Client, error) {
	conn, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{addr, conn}, nil
}

func (c *Client) Open(path string) (int64, error) {
	res := &Response{}
	if err := c.Call("Rpc.Open", FileRequest{Name: path}, &res); err != nil {
		return 0, err
	}

	return res.Id, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) CloseSession(id int64) error {
	res := &Response{}
	if err := c.Call("Rpc.Close", Request{id}, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) Stat(path string) (*StatResponse, error) {
	res := &StatResponse{}
	if err := c.Call("Rpc.Stat", FileRequest{Name: path}, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Read(id int64, buf []byte) (int, error) {
	res := &ReadResponse{Data: buf}
	if err := c.Call("Rpc.Read", ReadRequest{Id: id, Size: cap(buf)}, &res); err != nil {
		return 0, err
	}

	return res.Size, nil
}

func (c *Client) ReadAt(id, offset int64, size int) ([]byte, error) {
	res := &ReadResponse{Data: make([]byte, size)}

	err := c.Call("Rpc.ReadAt", ReadRequest{Id: id, Size: size, Offset: offset}, &res)
	if res.EOF {
		err = errors.Wrap(err, io.EOF.Error())
	}

	if size != res.Size {
		return res.Data[:res.Size], err
	}

	return res.Data, nil
}

func (c *Client) GetBlock(id int64, block int64) ([]byte, error) {
	return c.ReadAt(id, block*BlockSize, int(BlockSize))
}

func (c *Client) Download(from, to string) error {
	return c.DownloadAt(from, to, 0)
}

func (c *Client) DownloadAt(from, to string, block int64) error {
	stat, err := c.Stat(from)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.Errorf("%s is a directory", from)
	}

	blocks := stat.Size / BlockSize
	if stat.Size%BlockSize != 0 {
		blocks++
	}

	out, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer out.Close()

	id, err := c.Open(from)
	if err != nil {
		return err
	}
	defer c.CloseSession(id)

	var wg sync.WaitGroup

	for i := block; i < blocks; i++ {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()

			buf, err := c.GetBlock(id, i)
			if err != nil && err != io.EOF {
				log.Error(err, "unable to retrieve block")
				return
			}

			if _, err := out.WriteAt(buf, i*BlockSize); err != nil {
				log.Error(err, "unable to write block")
				return
			}

			if i%((blocks-id)/100+1) == 0 {
				log.Infof("downloading %s [%d/%d] blocks", from, i-block+1, blocks-block)
			}
		}(i)
	}

	wg.Wait()
	return nil
}
