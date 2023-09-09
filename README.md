# Login-STU

## 一、使用说明

- 用于在无头环境，比如Linux服务器上登录STU校园网；
- 用于在桌面环境双击登录校园网，简化了打开浏览器的步骤。

## 二、 不支持的系统（不完整列表）

Go 1.21 不支持的所有系统；因为软件使用 Go 编写， Go 版本是 1.21。

> Go 1.20 是支持 Microsoft Windows 7 / 8 / Server 2008 / Server 2012 的最后一个版本。自 Go 1.21 开始，用户需要在 Windows 10 或 Windows Server 2016 及更高版本上运行。
> 
> Go 1.20 也放弃了对 macOS 10.13 和 10.14 系列的支持。

## 三、How

URL: https://a.stu.edu.cn:444/ac_portal/login.php
Method: POST

### 1. Login

```go
payload := url.Values{
    "opr":         {"pwdLogin"},
    "userName":    {username},
    "pwd":         {password},
    "ipv4or6":     {""},
    "rememberPwd": {"0"},
}

```

### 2. Logout 

```go

payload := url.Values{

    "opr":      {"logout"}, 
    "ipv4or6" : {""}
}

```

