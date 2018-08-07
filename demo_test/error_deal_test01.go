package main

/*
type error interface {
	Error() string
}
*/
type SyntaxError struct {
	msg    string // 错误描述
	Offset int64  // 错误发生的位置
}
type Human struct {
	name string
}

func (e *SyntaxError) Error() string { return e.msg }

func (h *Human) Decode(human *Human) *SyntaxError { // 错误，将可能导致上层调用者err!=nil的判断永远为true。
	var err *SyntaxError // 预声明错误变量
	//if 出错条件 {
	//	err = &SyntaxError{}
	//}
	return err // 错误，err永远等于非nil，导致上层调用者err!=nil的判断始终为true
}
func main() {
	h := Human{}
	val := Human{}
	if err := h.Decode(&val); err != nil {
		//if serr, ok := err.(*json.SyntaxError); ok {
		//	line, col := findLine(f, serr.Offset)
		//	return fmt.Errorf("%s:%d:%d: %v", f.Name(), line, col, err)
		//}
		//return err
	}
}
