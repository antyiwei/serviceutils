package main

import (
	"net/http"

	"git.witalk.cn/base/wsclientsdk"
	"github.com/luaxlou/goutils/tools/logutils"
)

type TestReq struct {
	Word string
}

func main() {

	wsc := wsclientsdk.New("127.0.0.1:9999",
		http.Header{
			"APP": []string{
				"wsclienttest",
			}, "MER": []string{
				"weiliao",
			},
		})

	wsc.Ready = func() {

		for i := 0; i < 10; i++ {

			wsc.SendRequest("some", TestReq{
				Word: "hello1",
			}, func(code int, r *wsclientsdk.WsResult) {

				logutils.PrintObj(code)

			})
		}

	}

	select {}
}
