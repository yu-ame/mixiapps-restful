package config

import (
    "io/ioutil"
    "log"
    "encoding/json"
)

var config map[string]interface{}

func init(){
    conf,_ := ioutil.ReadFile("./configs/config.json")

    log.Printf(string(conf))
    var f interface{}
    json.Unmarshal(conf, &f)
    config = f.(map[string]interface{})  
}

func GetString(name string) string{
    return config[name].(string)    
}

