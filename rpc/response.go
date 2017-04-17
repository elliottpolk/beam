//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/rpc/response.go
//
package rpc

import "time"

const (
	DirectoryType string = "directory"
	FileType      string = "file"
)

type Response struct {
	Id     int64
	Result bool
}

type GetResponse struct {
	BlockId int
	Size    int64
	Data    []byte
}

type ReadRequest struct {
	Id     int64
	Offset int64
	Size   int
}

type ReadResponse struct {
	Size int
	Data []byte
	EOF  bool
}

type StatResponse struct {
	Type         string
	Size         int64
	LastModified time.Time
}

func (r *StatResponse) IsDir() bool {
	return r.Type == DirectoryType
}
