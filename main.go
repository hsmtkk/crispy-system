package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hsmtkk/crispy-system/sessionstore"
	"github.com/hsmtkk/crispy-system/userstore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const cookieName = "sessionID"

func main() {
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	sessionStore, err := initSessionStore()
	if err != nil {
		log.Fatal(err)
	}
	userStore := userstore.NewMemoryImpl()
	hdl := newHandler(sessionStore, userStore)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/login", hdl.getLogin)
	e.POST("/login", hdl.postLogin)
	e.GET("/increment", hdl.increment)
	e.GET("/logout", hdl.logout)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type myHandler struct {
	sessionStore sessionstore.SessionStore
	userStore    userstore.UserStore
}

func newHandler(sessionStore sessionstore.SessionStore, userStore userstore.UserStore) *myHandler {
	return &myHandler{sessionStore, userStore}
}

const loginHTML = `
<html>
 <body>
   <form method="POST" action="/login">
    <input type="text" name="userID">
	<input type="submit" value="submit">
   </form>
 </body>
</html>
`

func (h *myHandler) getLogin(ectx echo.Context) error {
	return ectx.HTML(http.StatusOK, loginHTML)
}

func (h *myHandler) postLogin(ectx echo.Context) error {
	userID := ectx.Request().FormValue("userID")
	if userID == "" {
		return fmt.Errorf("empty userID")
	}
	sessionID := uuid.NewString()
	sessionID, err := h.sessionStore.NewSession(ectx.Request().Context(), sessionID, userID)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = sessionID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	ectx.SetCookie(cookie)
	return ectx.Redirect(http.StatusMovedPermanently, "/increment")
}

func (h *myHandler) increment(ectx echo.Context) error {
	cookie, err := ectx.Cookie(cookieName)
	if err != nil {
		log.Printf("session cookie does not exist; %s", err.Error())
		return ectx.Redirect(http.StatusMovedPermanently, "/login")
	}
	sessionID := cookie.Value
	userID, err := h.sessionStore.GetUserID(ectx.Request().Context(), sessionID)
	if err != nil {
		log.Printf("session does not exist; %s", err.Error())
		return ectx.Redirect(http.StatusMovedPermanently, "/login")
	}
	count, err := h.userStore.Increment(userID)
	if err != nil {
		return err
	}
	return ectx.String(http.StatusOK, strconv.Itoa(count))
}

func (h *myHandler) logout(ectx echo.Context) error {
	cookie, err := ectx.Cookie(cookieName)
	if err != nil {
		log.Print("session cookie does not exist")
		return ectx.Redirect(http.StatusMovedPermanently, "/login")
	}
	sessionID := cookie.Value
	if err := h.sessionStore.DeleteSession(ectx.Request().Context(), sessionID); err != nil {
		return err
	}
	cookie.Name = cookieName
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	ectx.SetCookie(cookie)
	return ectx.Redirect(http.StatusMovedPermanently, "/login")
}

func requiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("you must define %s environment variable", key)
	}
	return val, nil
}

func getPort() (int, error) {
	portStr, err := requiredEnv("PORT")
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as int; %w", portStr, err)
	}
	return port, nil
}

func initSessionStore() (sessionstore.SessionStore, error) {
	useRedis := os.Getenv("REDIS_HOST")
	if useRedis == "" {
		return sessionstore.NewMemoryImpl(), nil
	} else {
		redisHost, err := requiredEnv("REDIS_HOST")
		if err != nil {
			return nil, err
		}
		redisPortStr, err := requiredEnv("REDIS_PORT")
		if err != nil {
			return nil, err
		}
		redisPassword, err := requiredEnv("REDIS_PASSWORD")
		if err != nil {
			return nil, err
		}
		redisPort, err := strconv.Atoi(redisPortStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s as int; %w", redisPortStr, err)
		}
		return sessionstore.NewRedisImpl(redisHost, redisPort, redisPassword), nil
	}
}
