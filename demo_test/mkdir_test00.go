package main

import (
	"fmt"
	"os"
)

/**
func Mkdir(name string, perm FileMode) error
创建名称为name的目录，权限设置是perm，例如0777
func MkdirAll(path string, perm FileMode) error
根据path创建多级子目录，例如vickywang/test1/test2。
func Remove(name string) error
删除名称为name的目录，当目录下有文件或者其他目录时会出错
func RemoveAll(path string) error
根据path删除多级子目录，如果path是单个名称，那么该目录下的子目录全部删除。
*/
func main() {
	os.Mkdir("vickywang", 0777)
	os.MkdirAll("vickywang/test1/test2", 0777)
	err := os.Remove("vickywang")
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll("vickywang")
}
