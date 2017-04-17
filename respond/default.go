//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/respond/default.go
//
package respond

import (
	"fmt"
	"net/http"
)

func WithDefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprint(w, "ok")
}
