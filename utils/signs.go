/**
 * Project Name:wechat-robot
 * File Name:signs.go
 * Package Name:utils
 * Date:2019年07月16日 14:36
 * Function:
 * Copyright (c) 2019, Jason.Wang All Rights Reserved.
 */
package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// GetSign get the sign info
func GetSign(srcdata interface{}, bizkey string) string {
	md5ctx := md5.New()

	switch v := reflect.ValueOf(srcdata); v.Kind() {
	case reflect.String:
		md5ctx.Write([]byte(v.String() + bizkey))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Map:
		orderStr := orderParam(v.Interface(), bizkey)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Struct:
		orderStr := Struct2map(v.Interface(), bizkey)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	default:
		return ""
	}
}

func orderParam(source interface{}, bizKey string) (returnStr string) {
	switch v := source.(type) {
	case map[string]string:
		keys := make([]string, 0, len(v))

		for k := range v {
			if k == "sign" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var buf bytes.Buffer
		for _, k := range keys {
			if v[k] == "" {
				continue
			}
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}

			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v[k])
		}
		buf.WriteString(bizKey)
		returnStr = buf.String()
	case map[string]interface{}:
		keys := make([]string, 0, len(v))

		for k := range v {
			if k == "sign" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var buf bytes.Buffer
		for _, k := range keys {
			if v[k] == "" {
				continue
			}
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(toString(v[k]))
		}
		buf.WriteString(bizKey)
		returnStr = buf.String()
	}
	return
}

func GetSignEncodeToUpper(srcdata interface{}, bizkey string) string {
	sign := GetSignEncode(srcdata, bizkey)
	return strings.ToUpper(sign)
}

func GetSignEncode(srcdata interface{}, bizkey string) string {
	md5ctx := md5.New()
	switch v := reflect.ValueOf(srcdata); v.Kind() {
	case reflect.String:
		md5ctx.Write([]byte(v.String() + bizkey))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Map:
		orderStr := orderParam(v.Interface(), bizkey)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Struct:
		orderStr := Struct2mapEncode(v.Interface(), bizkey)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	default:
		return ""
	}
}

func Struct2mapEncode(content interface{}, bizKey string) string {
	url := url.Values{}
	var val map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&val); err != nil {
			fmt.Println(err)
		} else {
			for k, v := range val {
				val[k] = v
			}
		}
	}
	for k, v := range val {
		// 去除冗余未赋值struct
		if v == "" {
			continue
		}
		url.Add(k, toString(v))
	}
	body := url.Encode()
	return body + bizKey
}

func Struct2Map(content interface{}) map[string]string {
	result := make(map[string]string)
	var val map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&val); err != nil {
			fmt.Println(err)
		} else {
			for k, v := range val {
				val[k] = v
			}
		}
	}
	for k, v := range val {
		// 去除冗余未赋值struct
		if v == "" {
			continue
		}
		result[k] = toString(v)
	}
	return result
}

func Struct2map(content interface{}, bizKey string) string {
	var tempArr []string
	temString := ""
	var val map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&val); err != nil {
			fmt.Println(err)
		} else {
			for k, v := range val {
				val[k] = v
			}
		}
	}
	i := 0
	for k, v := range val {
		// 去除冗余未赋值struct
		if v == "" {
			continue
		}
		i++
		tempArr = append(tempArr, k+"="+toString(v))
	}
	sort.Strings(tempArr)
	for n, v := range tempArr {
		if n+1 < len(tempArr) {
			temString = temString + v + "&"
		} else {
			temString = temString + v + bizKey
		}
	}
	return temString
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'E', -1, 64)
}
func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func boolToString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func toString(arg interface{}) string {
	switch arg.(type) {
	case bool:
		return boolToString(arg.(bool))
	case float32:
		return floatToString(float64(arg.(float32)))
	case float64:
		return floatToString(arg.(float64))
		//case complex64:
		//  p.fmtComplex(complex128(f), 64, verb)
		//case complex128:
		//  p.fmtComplex(f, 128, verb)
	case int:
		return intToString(int64(arg.(int)))
	case int8:
		return intToString(int64(arg.(int8)))
	case int16:
		return intToString(int64(arg.(int16)))
	case int32:
		return intToString(int64(arg.(int32)))
	case int64:
		return intToString(int64(arg.(int64)))
	default:
		return fmt.Sprint(arg)
	}
}

func combinePath(pre string, path string) string {
	if pre != "" && path != "" {
		return pre + "." + path
	}
	return pre + path
}

//将一个map[string]interface打平
func FlatMap(prefix string, mapData map[string]interface{}) map[string]interface{} {
	v := reflect.ValueOf(mapData)
	res := make(map[string]interface{})
	foreachObj(prefix, v, res)
	return res
}

func foreachObj(pre string, v reflect.Value, res map[string]interface{}) {
	switch v.Kind() {
	case reflect.Ptr:
		foreachObj(pre, v.Elem(), res)
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			foreachObj(combinePath(pre, strconv.Itoa(i)), v.Index(i), res)
		}
	case reflect.Struct:
		vType := v.Type()
		for i := 0; i < v.NumField(); i++ {
			foreachObj(combinePath(pre, vType.Field(i).Name), v.Field(i), res)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			foreachObj(combinePath(pre, key.String()), v.MapIndex(key), res)
		}
	case reflect.Interface:
		foreachObj(combinePath(pre, ""), v.Elem(), res)
	default: // float, complex, bool, chan, string,int,func, interface
		res[pre] = v.Interface()
	}
}

func getTplExpressions(str string) []string {
	reg_str := `\$\{.*?\}`
	re, _ := regexp.Compile(reg_str)
	all := re.FindAll([]byte(str), 2)
	keyArrays := make([]string, 0)
	for _, item := range all {
		itemStr := string(item)
		if len(itemStr) > 3 {
			itemStr = itemStr[2 : len(itemStr)-1]
			keyArrays = append(keyArrays, itemStr)
		}
	}
	return keyArrays
}

// 将tpl中的占位符 替换为真实值 ${data.0.att1}
func ParseTpl(tpl string, data map[string]interface{}) string {
	if len(tpl) < 4 {
		return tpl
	}
	expressions := getTplExpressions(tpl)
	data = FlatMap("", data)
	for _, exp := range expressions {
		//fmt.Println("exp",exp)
		exp = strings.TrimSpace(exp)
		tpl = strings.Replace(tpl, "${"+exp+"}", toString(data[exp]), -1)
	}
	return tpl
}
