package shared

import "github.com/gin-gonic/gin"

func NamedRoute(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, ok := c.Get("namedRoute"); ok {
			path = val.(string) + "/" + path
		}

		c.Set("namedRoute", path)

		c.Next()

		c.Writer.Header().Set("X-Backend-Route", path)
	}
}

type CustomRoutes gin.RouterGroup

func (r *CustomRoutes) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	h := append(handlers, NamedRoute(relativePath))

	return (*gin.RouterGroup)(r).POST(relativePath, h...)
}

func (r *CustomRoutes) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	h := append(handlers, NamedRoute(relativePath))

	return (*gin.RouterGroup)(r).GET(relativePath, h...)
}

func (r *CustomRoutes) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	h := append(handlers, NamedRoute(relativePath))

	return (*gin.RouterGroup)(r).PUT(relativePath, h...)
}

func (r *CustomRoutes) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	h := append(handlers, NamedRoute(relativePath))

	return (*gin.RouterGroup)(r).PATCH(relativePath, h...)
}

func (r *CustomRoutes) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	h := append(handlers, NamedRoute(relativePath))

	return (*gin.RouterGroup)(r).DELETE(relativePath, h...)
}
