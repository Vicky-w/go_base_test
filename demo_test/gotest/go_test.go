package gotest

import (
	"testing"
)

/**
文件名必须是_test.go结尾的，这样在执行go test的时候才会执行到相应的代码
必须import testing这个包
测试用例函数必须是Test开头
函数中通过调用testing.T的Error, Errorf, FailNow, Fatal, FatalIf方法，说明测试不通过，调用Log方法用来记录测试的信息
*/

func Test_Division_1(t *testing.T) {
	if i, e := Division(6, 2); i != 3 || e != nil { //try a unit test on function
		t.Error("除法函数测试没通过") // 如果不是如预期的那么就报错
	} else {
		t.Log("第一个测试通过了") //记录一些你期望记录的信息
	}
}

//func Test_Division_2(t *testing.T) {
//	t.Error("就是不通过")
//}

func Test_Division_2(t *testing.T) {
	if _, e := Division(6, 0); e == nil { //try a unit test on function
		t.Error("Division did not work as expected.") // 如果不是如预期的那么就报错
	} else {
		t.Log("one test passed.", e) //记录一些你期望记录的信息
	}
}

/**

go test
go test -v

*/
