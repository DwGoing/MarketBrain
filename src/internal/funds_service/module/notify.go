package module

import (
	"bytes"
	"errors"
	"net/http"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
type Notify struct {
	Config *Config `normal:""`
}

// @title	发送通知
// @param	Self	*Notify	模块实例
// @param	url		string	回调Url
// @param	data	[]byte	回调数据
func (Self *Notify) Send(url string, data []byte) error {
	if data == nil {
		data = []byte{}
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("request failed")
	}
	return nil
}
