# 包管理示例

本页面演示Go NPM SDK的高级包管理功能。

## 带选项安装包

### 开发依赖

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 安装开发依赖
    devPackages := []string{"jest", "typescript", "@types/node", "eslint"}
    
    for _, pkg := range devPackages {
        fmt.Printf("正在安装%s作为开发依赖...\n", pkg)
        err = client.InstallPackage(ctx, pkg, npm.InstallOptions{
            SaveDev:   true,
            SaveExact: true,
        })
        if err != nil {
            log.Printf("安装%s失败: %v", pkg, err)
        } else {
            fmt.Printf("%s安装成功！\n", pkg)
        }
    }
}
```

### 生产依赖

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 安装生产依赖
    prodPackages := map[string]string{
        "express":    "^4.18.0",
        "lodash":     "^4.17.21",
        "axios":      "^1.0.0",
        "dotenv":     "^16.0.0",
    }
    
    for pkg, version := range prodPackages {
        fmt.Printf("正在安装%s@%s...\n", pkg, version)
        err = client.InstallPackage(ctx, pkg+"@"+version, npm.InstallOptions{
            SaveDev:   false,
            SaveExact: false,
        })
        if err != nil {
            log.Printf("安装%s失败: %v", pkg, err)
        } else {
            fmt.Printf("%s@%s安装成功！\n", pkg, version)
        }
    }
}
```

## 包信息和搜索

### 获取包信息

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 获取包信息
    packageName := "express"
    info, err := client.GetPackageInfo(ctx, packageName)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("包: %s@%s\n", info.Name, info.Version)
    fmt.Printf("描述: %s\n", info.Description)
    fmt.Printf("主页: %s\n", info.Homepage)
    fmt.Printf("许可证: %s\n", info.License)
    
    if info.Author != nil {
        fmt.Printf("作者: %s <%s>\n", info.Author.Name, info.Author.Email)
    }
    
    if info.Repository != nil {
        fmt.Printf("仓库: %s\n", info.Repository.URL)
    }
    
    fmt.Printf("关键词: %v\n", info.Keywords)
    
    // 显示最新版本
    fmt.Println("\n最新版本:")
    for tag, version := range info.DistTags {
        fmt.Printf("  %s: %s\n", tag, version)
    }
}
```

### 搜索包

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 搜索包
    query := "react testing"
    results, err := client.Search(ctx, query)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("'%s'的搜索结果:\n\n", query)
    
    for i, result := range results[:10] { // 显示前10个结果
        fmt.Printf("%d. %s@%s\n", i+1, result.Package.Name, result.Package.Version)
        fmt.Printf("   描述: %s\n", result.Package.Description)
        fmt.Printf("   评分: %.2f (质量: %.2f, 流行度: %.2f, 维护: %.2f)\n",
            result.Score.Final,
            result.Score.Detail.Quality,
            result.Score.Detail.Popularity,
            result.Score.Detail.Maintenance)
        
        if result.Package.Author != nil {
            fmt.Printf("   作者: %s\n", result.Package.Author.Name)
        }
        
        fmt.Println()
    }
}
```

## 包更新和维护

### 更新包

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 列出当前包
    packages, err := client.ListPackages(ctx, npm.ListOptions{
        Global: false,
        Depth:  0,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("当前包:")
    for _, pkg := range packages {
        fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
    }
    
    // 更新特定包
    packagesToUpdate := []string{"lodash", "axios", "express"}
    
    fmt.Println("\n正在更新包...")
    for _, pkg := range packagesToUpdate {
        fmt.Printf("正在更新%s...\n", pkg)
        err = client.UpdatePackage(ctx, pkg)
        if err != nil {
            log.Printf("更新%s失败: %v", pkg, err)
        } else {
            fmt.Printf("%s更新成功！\n", pkg)
        }
    }
}
```

### 卸载包

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 卸载开发依赖
    devPackagesToRemove := []string{"@types/jest", "ts-node"}
    
    for _, pkg := range devPackagesToRemove {
        fmt.Printf("正在卸载%s...\n", pkg)
        err = client.UninstallPackage(ctx, pkg, npm.UninstallOptions{
            SaveDev: true,
        })
        if err != nil {
            log.Printf("卸载%s失败: %v", pkg, err)
        } else {
            fmt.Printf("%s卸载成功！\n", pkg)
        }
    }
    
    // 卸载生产依赖
    prodPackagesToRemove := []string{"unused-package"}
    
    for _, pkg := range prodPackagesToRemove {
        fmt.Printf("正在卸载%s...\n", pkg)
        err = client.UninstallPackage(ctx, pkg, npm.UninstallOptions{
            SaveDev: false,
        })
        if err != nil {
            log.Printf("卸载%s失败: %v", pkg, err)
        } else {
            fmt.Printf("%s卸载成功！\n", pkg)
        }
    }
}
```

## 使用package.json

### 读取和修改package.json

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    // 加载现有package.json
    pkg := npm.NewPackageJSON("./package.json")
    
    err := pkg.Load()
    if err != nil {
        log.Fatal(err)
    }
    
    // 显示当前信息
    fmt.Printf("当前包: %s@%s\n", pkg.GetName(), pkg.GetVersion())
    fmt.Printf("描述: %s\n", pkg.GetDescription())
    fmt.Printf("作者: %s\n", pkg.GetAuthor())
    
    // 修改包信息
    pkg.SetDescription("我的项目的更新描述")
    pkg.SetAuthor("新作者 <new.author@example.com>")
    
    // 添加新依赖
    pkg.AddDependency("moment", "^2.29.0")
    pkg.AddDevDependency("nodemon", "^2.0.0")
    
    // 添加脚本
    pkg.AddScript("dev", "nodemon src/index.js")
    pkg.AddScript("build", "webpack --mode production")
    pkg.AddScript("test:watch", "jest --watch")
    
    // 添加关键词
    pkg.AddKeyword("nodejs")
    pkg.AddKeyword("javascript")
    pkg.AddKeyword("api")
    
    // 设置仓库信息
    repo := &npm.Repository{
        Type: "git",
        URL:  "https://github.com/username/repo.git",
    }
    pkg.SetRepository(repo)
    
    // 设置错误报告信息
    bugs := &npm.Bugs{
        URL:   "https://github.com/username/repo/issues",
        Email: "bugs@example.com",
    }
    pkg.SetBugs(bugs)
    
    // 设置主页
    pkg.SetHomepage("https://github.com/username/repo#readme")
    
    // 保存前验证
    if err := pkg.Validate(); err != nil {
        log.Fatal(err)
    }
    
    // 保存更改
    err = pkg.Save()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("package.json更新成功！")
}
```

## 全局包管理

### 管理全局包

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 安装全局包
    globalPackages := []string{
        "typescript",
        "nodemon",
        "@angular/cli",
        "create-react-app",
        "eslint",
    }
    
    fmt.Println("正在安装全局包...")
    for _, pkg := range globalPackages {
        fmt.Printf("正在全局安装%s...\n", pkg)
        err = client.InstallPackage(ctx, pkg, npm.InstallOptions{
            Global: true,
        })
        if err != nil {
            log.Printf("全局安装%s失败: %v", pkg, err)
        } else {
            fmt.Printf("%s全局安装成功！\n", pkg)
        }
    }
    
    // 列出全局包
    fmt.Println("\n正在列出全局包...")
    globalPkgs, err := client.ListPackages(ctx, npm.ListOptions{
        Global: true,
        Depth:  0,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("全局包:")
    for _, pkg := range globalPkgs {
        fmt.Printf("  %s@%s\n", pkg.Name, pkg.Version)
    }
}
```

## 注册表配置

### 使用自定义注册表

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
    client, err := npm.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // 从自定义注册表安装
    customRegistry := "https://registry.npmjs.org/"
    
    fmt.Println("正在从自定义注册表安装包...")
    err = client.InstallPackage(ctx, "lodash", npm.InstallOptions{
        Registry: customRegistry,
        SaveDev:  false,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("从自定义注册表安装包成功！")
    
    // 安装私有包
    privateRegistry := "https://npm.company.com/"
    
    fmt.Println("正在安装私有包...")
    err = client.InstallPackage(ctx, "@company/private-package", npm.InstallOptions{
        Registry: privateRegistry,
        SaveDev:  false,
    })
    if err != nil {
        log.Printf("安装私有包失败: %v", err)
    } else {
        fmt.Println("私有包安装成功！")
    }
}
```

## 最佳实践

1. **对关键依赖使用精确版本**: 对重要包使用`SaveExact: true`
2. **分离开发和生产依赖**: 使用适当的`SaveDev`设置
3. **验证package.json**: 保存更改前始终验证
4. **优雅地处理错误**: 检查特定错误类型
5. **使用适当的注册表**: 为私有包配置自定义注册表
6. **保持依赖更新**: 定期更新包以确保安全性
7. **清理未使用的包**: 删除不再需要的包

## 下一步

- [便携版安装示例](./portable-installation.md)
- [高级功能示例](./advanced-features.md)
