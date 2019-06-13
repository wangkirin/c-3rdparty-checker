package main

import (
	"io/ioutil"
	"os"
	"strings"
	"fmt"
	"io"
	"bufio"
	"reflect"
	"errors"
	"path/filepath"
	"sort"
)

func main() {
	dirpath:="./"
	if len(os.Args)>1{
		dirpath=os.Args[1]
	}
	dirpath,_=filepath.Abs(dirpath)
	fmt.Printf("Scan Path: %s \n",dirpath)
	allHeaderFiles,err:=GetAllHeaderFiles(dirpath,dirpath)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
	allFiles,err:=GetAllFiles(dirpath)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}

	m:=make(map[string]string)

	for _,filePath:=range allFiles {
		includes,err:=GetAllIncludesByfile(filePath)
		if err!=nil {
			fmt.Println(err.Error())
			return
		}
		for _,include:=range includes {
			bool,_:=Contain(include,allHeaderFiles)
			if !bool {
				m[include]=""
			}
		}
	}

	var out[]string
	for k,_:=range m {
		out=append(out,k)
	}
	sort.Strings(out)
	for _,s:=range out{
		fmt.Println(s)
	}

}




//https://blog.csdn.net/robertkun/article/details/78744464
//获取指定目录下的所有文件,包含子目录下的文件
func GetAllHeaderFiles(dirPth string, rootPath string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllHeaderFiles(dirPth + PthSep + fi.Name(),rootPath)
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".h")
			if ok {
				path:=strings.TrimPrefix(strings.TrimPrefix(dirPth,rootPath),"/")
				files = append(files, path+"/"+fi.Name())
				GetAllIncludesByfile(dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllHeaderFiles(table,rootPath)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			files = append(files, dirPth+PthSep+fi.Name())

		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}



func GetAllIncludesByfile(filePath string)(inclues []string,err error)  {
	f,err:=os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil,err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		content:=string(a)
		if strings.HasPrefix(content,"#include") {
			inclues=append(inclues,ParseInclue(content))
		}
	}
	return inclues,nil
}

func ParseInclue(line string)string  {
	slice:=strings.Fields(line)
	if len(slice)<2 {
		return ""
	}
	if strings.HasPrefix(slice[1],"<") {
		out:=strings.TrimPrefix(strings.TrimSuffix(slice[1],">"),"<")
		return out
	}
	if strings.HasPrefix(slice[1],"\"") {
		out:=strings.TrimPrefix(strings.TrimSuffix(slice[1],"\""),"\"")
		return out
	}
	return ""
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
//https://www.cnblogs.com/zsbfree/archive/2013/05/23/3094993.html
func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}