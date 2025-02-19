# zhihu-crawler
## 背景

以前在知乎写文章，现在不想在知乎写了，都爬下来以后自己写博客了。所有就有了这个项目。

## 操作步骤

使用`cookieEdit`复制导出的cookie放入根目录下的`zhihu.json`文件内，没有就创建一个`zhihu.json`

然后运行项目

```
go run /cmd/mian.go
```

接口触发爬取动作

```
curl --location --request POST 'http://127.0.0.1:12345/api/crawler/zhihu'
```

项目依赖MySQL，爬取后的内容会存下来。你可以直接在表中导出

![image-20241212165806131](D:\Desktop\GitHub\go-crawler\assets\image-20241212165806131.png)

## 导出文章

使用无头浏览器爬取知乎文章信息，然后使用https://github.com/chenluda/zhihu-download下载文章内容，具体可以看这个项目

你也可以直接使用根目录的`zhihu-download`做了些小优化，日志和下载方面能方便些
