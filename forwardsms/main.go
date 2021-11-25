package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func sendWechat(url string, SenderNumber string, ReceivingDateTime string, TextDecoded string, rule string) {

	smssend := fmt.Sprintf(`{"msgtype":"text","text":{"content":"触发规则: %s \n发送时间: %s \n发送人：%s \n短信内容：%s" \n}}`, rule, ReceivingDateTime, SenderNumber, TextDecoded)
	var jsonStr = []byte(smssend)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func main() {
	for true {
		//读取发送配置文件
		config := map[string]interface{}{}
		viperconfig := viper.New()
		viperconfig.SetConfigName("forward")
		viperconfig.SetConfigType("yaml")
		viperconfig.AddConfigPath("/data/config")
		viperconfig.ReadInConfig()
		viperconfig.Unmarshal(&config)

		//读取当前发送规则配置文件
		status := map[string]interface{}{}
		viperStatus := viper.New()
		viperStatus.SetConfigName("status")
		viperStatus.SetConfigType("yaml")
		viperStatus.AddConfigPath("/data/config")
		viperStatus.ReadInConfig()
		viperStatus.Unmarshal(&status)

		//写入测试
		// fmt.Println(viperStatus.Get("id"))
		// viperStatus.Set("id", 2)
		// if err := viperStatus.WriteConfig(); err != nil {
		// 	fmt.Println(err)
		// }

		//数据库查询
		db, err := gorm.Open(sqlite.Open("/data/db/sms.db"), &gorm.Config{})
		if err != nil {
			panic("数据库文件不存在")
		}
		inbox := []map[string]interface{}{}

		db.Debug().Select("ID", "TextDecoded", "SenderNumber", "ReceivingDateTime").Table("inbox").Where("id > ?", viperStatus.Get("id")).Find(&inbox)

		//循环获取到的所有短信内容
		for _, inboxfor := range inbox {

			//循环所有规则库
			for _, config_for := range config {
				//如果匹配规则符合all就发送默认通知
				if config_for.(map[string]interface{})["rule"].(string) == "all" {
					sendWechat(config_for.(map[string]interface{})["url"].(string), inboxfor["SenderNumber"].(string), inboxfor["ReceivingDateTime"].(string), inboxfor["TextDecoded"].(string), config_for.(map[string]interface{})["rule"].(string))
				}
				//匹配关键字给对应的机器人
				if strings.Contains(inboxfor["TextDecoded"].(string), config_for.(map[string]interface{})["rule"].(string)) {
					sendWechat(config_for.(map[string]interface{})["url"].(string), inboxfor["SenderNumber"].(string), inboxfor["ReceivingDateTime"].(string), inboxfor["TextDecoded"].(string), config_for.(map[string]interface{})["rule"].(string))
					// config_for.map[string]
				}

			}

			//写入最后一次发短信的ID
			viperStatus.Set("id", inboxfor["ID"])
			if err := viperStatus.WriteConfig(); err != nil {
				fmt.Println(err)
			}

		}
		//清空变量
		inbox = nil
		viperconfig = nil
		viperStatus = nil
		time.Sleep(5 * time.Second)
	}

}
