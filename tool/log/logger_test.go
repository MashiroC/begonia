package log

import "testing"

/*
 @Author: as
 @Date: Creat in 19:26 2021/3/7
 @Description: begonia
*/
func TestLevel_String(t *testing.T) {
l:=DefaultNewLogger()
	l.SetCaller()
	l.Info("123")
}