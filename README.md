# 上课啦我爱记单词
本项目仅用于技术学习，请勿用于任何不合规用途。
## 使用方式
创建```secret.env```文件，并如下填写：
```azure
API_KEY="sk-GJjPL7CEaaaaaaaaaaaaaaazOaspTpwZ57wq8h"
API_URL="https://api.openai.com/v1/chat/completions"
TOKENS="aaaaaaaa-2bb4-4f27-ba0a-08e6a9414989"
WEEK=9
MODE=0
TIME=100
```
其中 ```API_KEY```为openai的apikey，```API_URL```为openai接口调用的endpoint，``` TOKENS```为上课啦登录后的权限凭证，``` WEEK```是考试周次，```MODE```为考试模式，0为自测，1为考试，```TIME```为提交等待时间，防止提交过快导致出现不可预知问题，建议```TIME```为300以上。

完成环境文件配置后，在终端打开文件夹，运行```loveWord.exe```,等待提交即可。

