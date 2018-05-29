# wechat-robot

#### 项目介绍

> 项目采用[Gin实现](https://github.com/gin-gonic/gin#quick-start)，该项目原始的仓库不在这里。
> 属于私密项目

#### 使用软件

> - Redis：NoSQL中的经典成员，个人用的比较多
> - MySql：用过数据库的应该知道
> - Gin：Web框架，好评不错。最近自己“抄袭”的另外一个项目用这个框架实现
> - Gorm：数据库
> - [WeChat SDK](https://github.com/silenceper/wechat)

#### 安装教程

本机运行时候需要修改`refect.go`里面的`GetFilePath`方法

```go
func GetFilePath(fileName string) string {
	// 本地运行使用，我懒得弄了。将就先
	//return filepath.Join(".", string(os.PathSeparator), fileName)
	// 线上运行使用下面语句
	return filepath.Join(filepath.Dir(getExePath()), string(os.PathSeparator), fileName)
}
```

```bash
dep ensure -update
调用 linux.bat 编译Linux下运行的版本。这个是在windows环境下编译的。
```

#### 参考资料

- Gopkg.toml里面有，懒得写。没有这么多时间。先将就看吧。
- golang123
- beego里面的案例
- iris的案例
 
#### 效果体验

 微信公众号：`記憶先生(o-jiyi)`
 
 
#### 吐槽方式

> - 这样的机会我会留下来么？不存在的。
> - o(∩_∩)o 哈哈
