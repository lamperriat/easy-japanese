# Easy Japanese
Run locally:
```bash
# backend
cd backend
export API_KEYS="comma separated list of apikeys"
go run main.go

# frontend
REACT_APP_API_KEY=your_dev_key
npm install
npm start
```

语言/Language: 仅中文/Chinese Only Now
本项目的目的在于帮助每一位日语学习者简单建立起自己的知识体系。作为一个刚接触日语不久的学习者，对我而言缺乏一款可以有效帮助我复习、整理知识（包括词汇、语法点、例句）的工具，因而产生了开发本工具的想法。以web的形式搭建，一是能力有限，希望以最简单的形式使我自己在各个平台上都可以使用该工具；二是希望将项目部署在服务器上，和朋友一起学习，即未来可能添加少量的社交相关功能。

您可以：
* 使用我们提供的单词库、例句库，其中例句均由我个人在学习的过程中收集，如果您也是初学者，我想这会对您有很大帮助
* 建立自己的知识库，您可以在浏览我们的知识库时将任何您看到的知识加入自己的知识库，我们设计了一个简单的算法帮助您复习任何新学到的知识。您当然也可以在学习的过程中，将自己所学记录到知识库中。我们的复习算法会在你的知识库中按照权重(即你所熟悉的程度)随机抽取
* (待定)使用基于LLM的功能，需要您自行提供对其openai接口的LLM API的API key

Roadmap:
* 基本的字典CRUD操作
* 例句库
* 每个用户专用的数据库
* 按权重随机抽取的算法
* 基于LLM的功能
* 社交相关功能

我并不确定本项目是否会对您的日语学习产生我所想象的帮助，故暂时只提供中文
I am not sure if this project is beneficial to your Japanese study or not. English translation may be available after I am satisfied with our progress.
