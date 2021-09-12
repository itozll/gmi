package dingtalk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var defaultService = &nilservice{}

type nilservice struct{}

func (*nilservice) Send(typ string, body interface{}) error {
	log.Println("type", typ, "body", body, "DingTalk:Nil")
	return nil
}

const uri = "https://oapi.dingtalk.com/robot/send?access_token="

// format: https://developers.dingtalk.com/document/app/custom-robot-access/title-72m-8ag-pqw
func SendByMarkdown(token string, markdown map[string]interface{}) error {
	return Send(token, "markdown", markdown)
}

// format: https://developers.dingtalk.com/document/app/custom-robot-access/title-72m-8ag-pqw
func Send(token string, typ string, body interface{}) error {
	v := map[string]interface{}{
		"msgtype": typ,
		typ:       body,
	}

	dat, err := json.Marshal(v)
	if err != nil {
		log.Println("value", v, "err_message", err, "DingTalk")
		return err
	}

	resp, err := http.Post(uri+token, "application/json", strings.NewReader(string(dat)))
	if err != nil {
		log.Println("value", v, "err_message", err, "DingTalk Post")
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("value", v, "err_message", err, "DingTalk")
		return err
	}

	if resp.StatusCode != 200 {
		log.Println("value", v, "status", resp.StatusCode, "DingTalk Response")
		return fmt.Errorf("status=%d", resp.StatusCode)
	}

	return nil
}
