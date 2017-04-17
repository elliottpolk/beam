//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/respond/error.go
//
package respond

import (
	"fmt"
	"net/http"
)

func WithError(w http.ResponseWriter, statuscode int, format string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}
