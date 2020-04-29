package go_ini

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

/**
	@iniPath string 配置文件路径
    @dstStruct interface{} 映射结构体
*/
func Parse(iniPath string, dstStruct interface{}) (err error) {
	var (
		iniData       []byte
		lineArr       []string
		typeInfo      reflect.Type
		typeStruct    reflect.Type
		lastFieldName string
	)

	typeInfo = reflect.TypeOf(dstStruct)
	if typeInfo.Kind() != reflect.Ptr {
		err = fmt.Errorf("Parse() second param is not pointer")
		return
	}

	typeStruct = typeInfo.Elem()
	if typeStruct.Kind() != reflect.Struct {
		err = fmt.Errorf("Parse() second param is not struct")
		return
	}

	if iniData, err = ioutil.ReadFile(iniPath); err != nil {
		err = fmt.Errorf("read config file failed:%v", err)
		return
	}

	lineArr = strings.Split(string(iniData), "\n") //行数据
	//fmt.Printf("%#v",lineArr)

	for index, value := range lineArr {
		line := strings.TrimSpace(value)
		if len(line) == 0 {
			continue
		}

		//如果是注释，直接忽略
		if line[0] == ';' || line[0] == '#' {
			continue
		}

		//解析节点
		if line[0] == '[' {
			lastFieldName, err = parseSection(line, typeStruct)
			//fmt.Printf("lastFieldName:%v\n", lastFieldName)
			if err != nil {
				err = fmt.Errorf("%v lineno:%d", err, index+1)
				return
			}
			continue
		}

		//解析选项
		err = parseItem(lastFieldName, line, dstStruct)
		if err != nil {
			err = fmt.Errorf("%v lineno:%d", err, index+1)
			return
		}
	}

	return
}

func parseSection(line string, typeInfo reflect.Type) (fieldName string, err error) {

	if line[0] == '[' && len(line) <= 2 {
		err = fmt.Errorf("syntax error, invalid section:%s", line)
		return
	}

	if line[0] == '[' && line[len(line)-1] != ']' {
		err = fmt.Errorf("syntax error, invalid section:%s", line)
		return
	}

	if line[0] == '[' && line[len(line)-1] == ']' {
		sectionName := strings.TrimSpace(line[1 : len(line)-1])
		if len(sectionName) == 0 {
			err = fmt.Errorf("syntax error, invalid section:%s", line)
			return
		}

		for i := 0; i < typeInfo.NumField(); i++ {
			field := typeInfo.Field(i)
			tagValue := field.Tag.Get("ini")

			if tagValue != "" {

				if tagValue == sectionName {
					fieldName = field.Name
					break
				}

			} else {

				tagValue = strings.ToLower(field.Name)
				if tagValue == sectionName {
					fieldName = field.Name
					break
				}

			}
		}
	}

	return
}

func parseItem(lastFieldName, line string, dstStruct interface{}) (err error) {
	index := strings.Index(line, "=")
	if index == -1 {
		err = fmt.Errorf("syntax error, line:%s", line)
		return
	}

	key := strings.TrimSpace(line[0:index])
	val := strings.TrimSpace(line[index+1:])

	if len(key) == 0 {
		err = fmt.Errorf("syntax error, line:%s", line)
		return
	}

	resultValue := reflect.ValueOf(dstStruct)
	sectionValue := resultValue.Elem().FieldByName(lastFieldName)

	sectionType := sectionValue.Type()
	if sectionType.Kind() != reflect.Struct {
		err = fmt.Errorf("field:%s must be struct", lastFieldName)
		return
	}

	keyFieldName := ""
	defaultValue := ""

	for i := 0; i < sectionType.NumField(); i++ {

		field := sectionType.Field(i)               //获取字段
		tagValue := field.Tag.Get("ini")        //获取tag:ini
		defaultValue = field.Tag.Get("default") //获取默认值tag:default

		if tagValue != "" {
			if tagValue == key {
				keyFieldName = field.Name
				break
			}
		} else {
			tagValue = strings.ToLower(field.Name)
			if tagValue == key {
				keyFieldName = field.Name
				break
			}
		}
	}

	if len(keyFieldName) == 0 {
		return
	}

	fieldValue := sectionValue.FieldByName(keyFieldName)
	if fieldValue == reflect.ValueOf(nil) {
		return
	}
	//获取选项字段的类型,并设置值(如果该项值为空,则设置默认值)
	switch fieldValue.Type().Kind() {
	case reflect.String:
		//检查该值是否为空,如果为空,则查看其是否有默认值,如果有则设置,没有则设置空字符串
		if defaultValue != "" && val == "" {
			fieldValue.SetString(defaultValue)
		} else {
			fieldValue.SetString(val)
		}
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:

		if defaultValue != "" && val == "" {
			defaultInt, errRet := strconv.ParseInt(defaultValue, 10, 64) //字符串转10进制数字
			if errRet != nil {
				err = errRet
				return
			}
			fieldValue.SetInt(defaultInt)
		} else {
			intVal, errRet := strconv.ParseInt(val, 10, 64) //字符串转10进制数字
			if errRet != nil {
				err = errRet
				return
			}
			fieldValue.SetInt(intVal)
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:

		if defaultValue != "" && val == "" {
			defaultUint, errRet := strconv.ParseUint(defaultValue, 10, 64) //字符串转10进制数字
			if errRet != nil {
				err = errRet
				return
			}
			fieldValue.SetUint(defaultUint)
		} else {
			uIntVal, errRet := strconv.ParseUint(val, 10, 64)
			if errRet != nil {
				err = errRet
				return
			}
			fieldValue.SetUint(uIntVal)
		}

	case reflect.Float32, reflect.Float64:
		if defaultValue != "" && val == "" {
			defaultFloat, errRet := strconv.ParseFloat(defaultValue, 64) //字符串转10进制数字
			if errRet != nil {
				err = errRet
				return
			}
			fieldValue.SetFloat(defaultFloat)
		} else {
			floatVal, errRet := strconv.ParseFloat(val, 64)
			if errRet != nil {
				return
			}
			fieldValue.SetFloat(floatVal)
		}

	default:
		err = fmt.Errorf("unsupport type:%v", fieldValue.Type().Kind())
	}

	return
}
