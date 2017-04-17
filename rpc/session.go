//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/rpc/session.go
//
package rpc

import (
	"os"
	"sync"
)

type session struct {
	*sync.Mutex

	files   map[int64]*os.File
	counter int64
}

func (s *session) add(f *os.File) int64 {
	s.Lock()
	defer s.Unlock()

	s.counter++
	s.files[s.counter] = f

	return s.counter
}

func (s *session) get(id int64) *os.File {
	s.Lock()
	defer s.Unlock()

	return s.files[id]
}

func (s *session) delete(id int64) {
	s.Lock()
	defer s.Unlock()

	if f, exist := s.files[id]; exist {
		f.Close()
		delete(s.files, id)
	}
}

func (s *session) len() int {
	return len(s.files)
}
