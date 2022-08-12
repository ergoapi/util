package file

import (
	"bufio"
	"os"
)

/* https://github.com/zgg2001/StringFinderZ/blob/master/src/detect/detect_binary.go
https://blog.csdn.net/qq_45698148/article/details/120930607
*/

// IsBinary returns true if the file is binary.
func IsBinary(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	r := bufio.NewReader(file)
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)

	whiteByte := 0
	for i := 0; i < n; i++ {
		if (buf[i] >= 0x20 && buf[i] <= 0xff) ||
			buf[i] == 9 ||
			buf[i] == 10 ||
			buf[i] == 13 {
			whiteByte++
		} else if buf[i] <= 6 || (buf[i] >= 14 && buf[i] <= 31) {
			return true
		}
	}

	if whiteByte >= 1 {
		return false
	}
	return true
}
