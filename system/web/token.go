package web

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"acmweb/system/config"
)

// 指定加密密钥
var jwtSecret = []byte("87953z29310576g68x666zxx9z305722")

// Claims 是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	UserId   int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

type ZToken struct{}

func NewToken() *ZToken {
	return &ZToken{}
}

// Generate 根据用户的用户名和密码产生token
func (c *ZToken) Generate(userId int, username, password string) (string, error) {
	// 设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(config.CONFIG.Application.ExpireTime) * time.Hour)

	claims := Claims{
		UserId:   userId,
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "ACMManager",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// Parse 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func (c *ZToken) Parse(token string) (*Claims, error) {

	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*ZToken
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err

}
