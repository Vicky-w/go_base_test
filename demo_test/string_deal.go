package main

import (
	"fmt"
	"strconv"
	"strings"
)

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func main() {

	/**
	func Contains(s, substr string) bool
	字符串s中是否包含substr，返回bool值
	*/
	fmt.Println(strings.Contains("seafood", "foo")) //true
	fmt.Println(strings.Contains("seafood", "bar")) //false
	fmt.Println(strings.Contains("seafood", ""))    //true
	fmt.Println(strings.Contains("", ""))           //true
	/**
	func Join(a []string, sep string) string
	字符串链接，把slice a通过sep链接起来
	*/
	s := []string{"foo", "bar", "baz"}
	fmt.Println(strings.Join(s, ", ")) //foo, bar, baz

	/**
	func Index(s, sep string) int
	在字符串s中查找sep所在的位置，返回位置值，找不到返回-1
	*/

	fmt.Println(strings.Index("chicken", "ken")) //4
	fmt.Println(strings.Index("chicken", "dmr")) //-1

	/**
	func Repeat(s string, count int) string
	重复s字符串count次，最后返回重复的字符串
	*/
	fmt.Println("ba" + strings.Repeat("na", 2)) //banana

	/**
	func Replace(s, old, new string, n int) string
	在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	*/
	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))      //oinky oinky oink
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1)) //moo moo moo
	/*
		func Split(s, sep string) []string
		把s字符串按照sep分割，返回slice
	*/
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))                        //["a" "b" "c"]
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a ")) //["" "man " "plan " "canal panama"]
	fmt.Printf("%q\n", strings.Split(" xyz ", ""))                         //[" " "x" "y" "z" " "]
	fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))            //[""]
	/*
		func Trim(s string, cutset string) string
		在s字符串的头部和尾部去除cutset指定的字符串
	*/
	fmt.Printf("[%q]", strings.Trim(" !!! Achtung !!! ", "! ")) //return string     ["Achtung"]
	/**
	func Fields(s string) []string
	去除s字符串的空格符，并且按照空格分割返回slice
	*/
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   ")) //["foo" "bar" "baz"]

	fmt.Println("\n============================Append 系列函数将整数等转换为字符串后，添加到现有的字节数组中==================================")
	str := make([]byte, 0, 100)
	str = strconv.AppendInt(str, 4567, 10)
	str = strconv.AppendBool(str, false)
	str = strconv.AppendQuote(str, "abcdefg")
	str = strconv.AppendQuoteRune(str, '单')
	fmt.Println(string(str))
	fmt.Println("============================Format 系列函数把其他类型的转换为字符串==================================")
	a := strconv.FormatBool(false)
	b := strconv.FormatFloat(123.23, 'g', 12, 64)
	c := strconv.FormatInt(1234, 10)
	d := strconv.FormatUint(12345, 10)
	e := strconv.Itoa(1023)
	fmt.Println(a, b, c, d, e)
	fmt.Println("============================Parse 系列函数把字符串转换为其他类型==================================")
	a2, err := strconv.ParseBool("false")
	checkError(err)
	b2, err := strconv.ParseFloat("123.23", 64)
	checkError(err)
	c2, err := strconv.ParseInt("1234", 10, 64)
	checkError(err)
	d2, err := strconv.ParseUint("12345", 10, 64)
	checkError(err)
	e2, err := strconv.Atoi("1023")
	checkError(err)
	fmt.Println(a2, b2, c2, d2, e2)
}
