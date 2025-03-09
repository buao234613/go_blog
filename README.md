## 快速开始
在启动之前，请确保填写好*config.json*配置文件：
```json
{
    "authConfig": {
        "Username": "admin",
        "Password": "123456",
        "Enable2FA": false
    },
    "databaseConfig": {
        "Driver": "mysql",
        "Host": "127.0.0.1",
        "Port": "3306",
        "Database": "tomato",
        "Username": "root",
        "Password": "msql",
        "Charset": "utf8mb4"
    },
    "user": {
        "Username": "xxxx",
        "Gravatar": "/",
        "Url": "xxxx",
        "Description": "(⊙o⊙)？"
    },
    "site": {
        "Path": "http://localhost:7891/",
        "Title": "xxxxx",
        "Icon": "",
        "Bili": "",
        "Github": "",
        "Twitter": "/",
        "Mail": ""
    }
// 新增加es查询功能
    "esConfig":{
        "Host": "127.0.0.1",
        "Port": "9200"
    }
}
```
第一次启动时，将*Enable2FA*设置为*false*，在登录后台查看二次验证的密钥后，根据需求改为true。

