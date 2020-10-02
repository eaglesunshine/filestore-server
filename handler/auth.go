package handler

import (
	"net/http"
)

//HTTPInterceptor: http请求拦截器
//装饰器(decorator)是一个这样的函数：它的参数是具体类型的函数，并且返回值也是和参数相同类型的函数。
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")

			if len(username) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w, r)
		})
}
