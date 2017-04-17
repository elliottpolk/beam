//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/respond/json.go
//
package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		fmt.Fprint(w, `{"success": true}`)
		return
	}

	out, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		WithError(w, http.StatusInternalServerError, "unable to convert results to json: %v\n", err)
		return
	}

	fmt.Fprint(w, string(out))
}
