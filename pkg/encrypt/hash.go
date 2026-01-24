package encrypt

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func Md5(str []byte) string {
	// 1. 创建一个新的MD5哈希计算实例（初始化MD5算法）
	h := md5.New()
	// 2. 将输入的字节数组写入哈希实例
	h.Write(str)
	// h.Sum(nil)：生成原始的16字节二进制MD5值（nil表示不追加到现有字节数组）
	// hex.EncodeToString：将16字节二进制数据转换成32位十六进制字符串
	return hex.EncodeToString(h.Sum(nil))
}

// hash加密
func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// hash校验
func ValidatePasswordHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
