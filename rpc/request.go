//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/rpc/request.go
//
package rpc

type Request struct {
	Id int64
}

type FileRequest struct {
	Name string
}

type GetRequest struct {
	Id      int64
	BlockId int
}
