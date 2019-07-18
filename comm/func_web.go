package comm

import (
	"crypto/md5"
	"fmt"
	"log"
	"lottery/conf"
	"lottery/models"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func ClientIp(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func Redirect(writer http.ResponseWriter, url string) {
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusFound)
}

func GetLoginUser(r *http.Request) *models.ObjLoginUser {
	c, err := r.Cookie("lottery_loginuser")
	if err != nil {
		return nil
	}
	params, err := url.ParseQuery(c.Value)
	if err != nil {
		return nil
	}
	uid, err := strconv.Atoi(params.Get("uid"))
	if err != nil {
		return nil
	}
	now, err := strconv.Atoi(params.Get("now"))
	if err != nil {
		return nil
	}
	loginuser := &models.ObjLoginUser{
		Uid:      uid,
		Username: params.Get("name"),
		Now:      now,
		Ip:       ClientIp(r),
		Sign:     params.Get("sign"),
	}
	if sign := createLoginUserSign(loginuser); sign != loginuser.Sign {
		log.Println("comm func_web.createLoginUserSign not signed", sign, loginuser.Sign)
		return nil
	}

	return loginuser
}

func SetLoginUser(w http.ResponseWriter, loginUser *models.ObjLoginUser) {
	if loginUser == nil || loginUser.Uid < 1 {
		c := &http.Cookie{
			Name:   "lottery_loginuser",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		return
	}
	if loginUser.Sign == "" {
		loginUser.Sign = createLoginUserSign(loginUser)
	}
	params := url.Values{}
	params.Add("uid", strconv.Itoa(loginUser.Uid))
	params.Add("now", strconv.Itoa(loginUser.Now))
	params.Add("name", loginUser.Username)
	params.Add("sign", loginUser.Sign)
	params.Add("ip", loginUser.Ip)
	c := &http.Cookie{
		Name:  "lottery_loginuser",
		Value: params.Encode(),
	}
	http.SetCookie(w, c)
}

func createLoginUserSign(loginUser *models.ObjLoginUser) string {
	str := fmt.Sprintf("uid=%d&username=%s&secret=%s%now=%d",
		loginUser.Uid, loginUser.Username, conf.CookieSecret, loginUser.Now)
	fmt.Println(str)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}
