package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogginMiddleware struct {
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func NewLogginMiddleware() *LogginMiddleware {
	return &LogginMiddleware{}
}

func (m *LogginMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()

		// 复制请求体（Body 只能读取一次）
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		pr := "[HTTP Request] %s %s"
		params := []any{
			r.Method,
			r.URL.String(),
		}

		// 记录请求信息
		headerStr := ""
		for k, v := range r.Header {
			headerStr += fmt.Sprintf("%s: %s, ", k, strings.Join(v, ","))
		}

		if headerStr != "" {
			pr += ", Header: %s"
			params = append(params, headerStr)
		}

		tokenStr := strings.TrimSpace(strings.TrimLeft(r.Header.Get("Authorization"), "Bearer"))
		if parsed, _ := ParseBindEmailToken(tokenStr, "YnVkb25n"); parsed != nil {
			pr += ", UID: %s, Email: %s "
			params = append(params, parsed.UserId, parsed.Email)
		}

		if len(bodyBytes) > 0 {
			bodyBytes = bytes.ReplaceAll(bodyBytes, []byte("\n"), []byte(""))
			bodyBytes = bytes.ReplaceAll(bodyBytes, []byte("\r"), []byte(""))
			pr += ", Body: %s"
			params = append(params, string(bodyBytes))
		}

		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next(rw, r)

		pr += ", Status: %d"
		params = append(params, rw.statusCode)

		var respBody ResponseBody
		if rw.body.Len() > 0 {
			json.Unmarshal(rw.body.Bytes(), &respBody)
			respBytes := bytes.ReplaceAll(rw.body.Bytes(), []byte("\n"), []byte(""))
			respBytes = bytes.ReplaceAll(respBytes, []byte("\r"), []byte(""))
			pr += ", Response: %s"
			params = append(params, string(respBytes))
		}

		if respBody.Code != 0 {
			logx.Errorf(pr, params...)
		} else {
			logx.Infof(pr, params...)
		}
	}
}

type ResponseBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ParseBindEmailToken(tokenStr string, secret string) (*BindEmailClaims, error) {
	claims := &BindEmailClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type BindEmailClaims struct {
	Provider string `json:"provider"`
	OpenID   string `json:"openID"`
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	jwt.RegisteredClaims
}

type MyClaims struct {
	Provider string `json:"provider"`
	OpenID   string `json:"openID"`
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	RoleId   string `json:"roleId"`
	UserName string `json:"userName"`
	Site     string `json:"site"`
	Version  string `json:"version"`
	UserType string `json:"userType"`
	jwt.RegisteredClaims
}
