//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/rpc/rpc.go
//
package rpc

import (
	"io"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"

	"github.com/elliottpolk/beam/log"

	"github.com/pkg/errors"
)

type Rpc struct {
	*session

	addr string
	dir  string
}

func RegisterAndHandle(addr, dir string) error {
	s := &session{&sync.Mutex{}, make(map[int64]*os.File), 0}
	if err := rpc.Register(&Rpc{s, addr, dir}); err != nil {
		return err
	}

	rpc.HandleHTTP()
	return nil
}

func (r *Rpc) Open(req FileRequest, res *Response) error {
	file, err := os.Open(filepath.Join(r.dir, req.Name))
	if err != nil {
		return err
	}

	res.Id = r.add(file)
	res.Result = true

	log.Infof("open %s, session %d", req.Name, res.Id)
	return nil
}

func (r *Rpc) Close(req Request, res *Response) error {
	r.delete(req.Id)
	res.Result = true

	log.Infof("closing session %d", req.Id)
	return nil
}

func (r *Rpc) Stat(req FileRequest, res *StatResponse) error {
	info, err := os.Stat(filepath.Join(r.dir, req.Name))
	if os.IsNotExist(err) {
		return err
	}

	res.Type = DirectoryType
	if !info.IsDir() {
		res.Type = FileType
		res.Size = info.Size()
	}
	res.LastModified = info.ModTime()

	log.Infof("stat %s, %+v", req.Name, res)
	return nil
}

func (r *Rpc) Read(req ReadRequest, res *ReadResponse) error {
	file := r.get(req.Id)
	if file == nil {
		return errors.New("rpc.Open must be called first")
	}

	res.Data = make([]byte, req.Size)
	n, err := file.Read(res.Data)
	if err != nil {
		if err != io.EOF {
			return err
		}

		res.EOF = true
	}

	res.Size = n
	res.Data = res.Data[:res.Size]

	// log.Infof("read session %d, read %d[bytes]", req.Id, res.Size)
	return nil
}

func (r *Rpc) ReadAt(req ReadRequest, res *ReadResponse) error {
	file := r.get(req.Id)
	if file == nil {
		return errors.New("Rpc.Open must be called first")
	}

	res.Data = make([]byte, req.Size)
	n, err := file.ReadAt(res.Data, req.Offset)
	if err != nil {
		if err != io.EOF {
			return err
		}

		res.EOF = true
	}

	res.Size = n
	res.Data = res.Data[:n]

	// log.Infof("read session %d, offset %d, n %d", req.Id, req.Offset, res.Size)
	return nil
}
