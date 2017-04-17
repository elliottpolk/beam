//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/server/server.go
//
package server

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/elliottpolk/beam/log"
	"github.com/elliottpolk/beam/respond"
	"github.com/elliottpolk/beam/rpc"
)

func mkdir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}

	return nil
}

func ListenAndServe(addr, dir string) error {
	if err := mkdir(dir); err != nil {
		return err
	}

	http.HandleFunc("/ls", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("request - %+v", r)

		if r.Method != http.MethodGet {
			respond.WithError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		info, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Error(err, "unable to read directory")
			respond.WithError(w, http.StatusInternalServerError, "unable to read directory")

			return
		}

		files := make([]string, 0)
		for _, i := range info {
			files = append(files, i.Name())
		}

		respond.WithJson(w, files)
	})

	if err := rpc.RegisterAndHandle(addr, dir); err != nil {
		return err
	}

	return http.ListenAndServe(addr, nil)
}
