// Time : 2020/9/26 21:43
// Author : Kieran

// ids
package ids

import 	uuid "github.com/satori/go.uuid"


// ids.go something

func New() string {
	return uuid.NewV4().String()
}