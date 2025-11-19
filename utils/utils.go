package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func GetLocalIp() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}

func CopyMap2Sync[K string, V any](origin map[K]V, new *sync.Map) {
	for k, v := range origin {
		new.Store(k, v)
	}
}

func Bytes2String(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d Bytes", n)
	}
	if n < 1048576 {
		return fmt.Sprintf("%.3f KB", float64(n)/1024)
	}
	if n < 1073741824 {
		return fmt.Sprintf("%.3f MB", float64(n)/1048576)
	}
	if n < 1099511627776 {
		return fmt.Sprintf("%.3f GB", float64(n)/1073741824)
	}
	return fmt.Sprintf("%.3f TB", float64(n)/1099511627776)
}

// compare version: a == b return 0, a > b return 1, a < b return -1
func VersionCompare(a, b, sep string) int {
	if a == b {
		return 0
	}
	if b == "" {
		return 1
	}
	if a == "" {
		return -1
	}
	if sep == "" {
		sep = "."
	}
	aTmp := strings.Split(a, sep)
	bTmp := strings.Split(b, sep)
	al := len(aTmp)
	bl := len(bTmp)
	for i := range min(al, bl) {
		ai := aTmp[i]
		bi := bTmp[i]
		if ai == bi {
			continue
		}

		an, aerr := strconv.Atoi(ai)
		bn, berr := strconv.Atoi(bi)
		if aerr != nil && berr == nil {
			return -1
		}
		if aerr == nil && berr != nil {
			return 1
		}
		if aerr == nil && berr == nil {
			if an > bn {
				return 1
			}
			return -1
		}

		ax := []rune(ai)[0]
		bx := []rune(bi)[0]
		if ax > bx {
			return 1
		}
		return -1
	}

	if al > bl {
		return 1
	}
	return -1
}

// check port is occupied
func PortCheck(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	defer l.Close()
	return true
}

// ToStr
func ToStr(i any) string {
	if i == nil {
		return ""
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64:
		return fmt.Sprintf("(%g+%gi)", real(v.Complex()), imag(v.Complex()))
	case reflect.Complex128:
		return fmt.Sprintf("(%g+%gi)", real(v.Complex()), imag(v.Complex()))
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Slice, reflect.Map, reflect.Struct, reflect.Array:
		str, _ := json.Marshal(i)
		return string(str)
	default:
		return ""
	}
}

func StringfyError(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
