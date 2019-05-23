package globalmw

import (
	"backend/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

type pathParser struct {
	port string
}

var skipPaths = []string{
	"/metadata",
	"/metrics",
	"/ping",
}

func (p *pathParser) ServeHTTP(ctx *bm.Context) {
	path := ctx.Request.RequestURI
	method := ctx.Request.Method
	if res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/metadata", p.port)); err != nil {
		ctx.String(400, "查询API列表发生错误：%v", err)
	} else {
		defer res.Body.Close()
		var jsonApis map[string]interface{}
		if body, err := ioutil.ReadAll(res.Body); err != nil {
			ctx.String(400, "读取API列表发生错误：%v", err)
		} else if err := json.Unmarshal(body, &jsonApis); err != nil {
			ctx.String(400, "解析API列表发生错误：%v", err)
		} else {
			pVarReg := regexp.MustCompile(":\\w+")
			// 遍历所有开放的API
			for p, m := range jsonApis["data"].(map[string]interface{}) {
				m = m.(map[string]interface{})["method"]
				if method != m {
					continue
				}
				// 如果路径被指定为跳过，则跳过
				if utils.Includes(skipPaths, p) {
					continue
				}
				// 替换所有路径中的路径变量为正则表达式：\w+
				sPathReg := pVarReg.ReplaceAllStringFunc(p, func(k string) string {
					fmt.Printf("Key: %s\n", k)
					return "\\w+"
				})
				// 如果做了替换的路径跟原路径是一样的，则意味着该路径没有路径变量
				if sPathReg == p {
					continue
				}
				if matched, err := regexp.MatchString(sPathReg, path); !matched || err != nil {
					continue
				}
				pathReg := regexp.MustCompile(sPathReg)
				for _, v := range pathReg.FindAllString(path, -1) {
					fmt.Printf("Value: %s\n", v)
				}
				// TODO
				fmt.Println(p)
				fmt.Println(m)
				fmt.Println("========================================")
			}
		}
	}
	ctx.Next()
}

func ParsePath(port string) (p *pathParser) {
	return &pathParser{
		port: port,
	}
}
