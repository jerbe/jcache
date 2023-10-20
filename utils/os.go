package utils

import "os"

/*
*

	@author : Jerbe - The porter from Earth
	@time : 2023/9/19 21:40
	@describe :
*/

// Hostname 获取主机名
func Hostname() string {
	host, _ := os.Hostname()
	return host
}
