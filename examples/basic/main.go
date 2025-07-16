package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/go-npm-sdk/pkg/npm"
)

func main() {
	ctx := context.Background()

	// 创建npm客户端
	client, err := npm.NewClient()
	if err != nil {
		log.Fatalf("创建npm客户端失败: %v", err)
	}

	fmt.Println("=== Go NPM SDK 基本使用示例 ===")

	// 1. 检查npm是否可用
	fmt.Println("\n1. 检查npm是否可用...")
	if !client.IsAvailable(ctx) {
		fmt.Println("npm不可用，正在尝试安装...")
		if err := client.Install(ctx); err != nil {
			log.Fatalf("安装npm失败: %v", err)
		}
		fmt.Println("npm安装成功！")
	} else {
		fmt.Println("npm已可用")
	}

	// 2. 获取npm版本
	fmt.Println("\n2. 获取npm版本...")
	version, err := client.Version(ctx)
	if err != nil {
		log.Printf("获取npm版本失败: %v", err)
	} else {
		fmt.Printf("npm版本: %s\n", version)
	}

	// 3. 创建测试项目目录
	fmt.Println("\n3. 创建测试项目...")
	projectDir := filepath.Join(os.TempDir(), "npm-sdk-test")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Fatalf("创建项目目录失败: %v", err)
	}
	defer os.RemoveAll(projectDir) // 清理

	fmt.Printf("项目目录: %s\n", projectDir)

	// 4. 初始化项目
	fmt.Println("\n4. 初始化npm项目...")
	initOptions := npm.InitOptions{
		Name:        "test-project",
		Version:     "1.0.0",
		Description: "Go NPM SDK测试项目",
		Author:      "Go NPM SDK",
		License:     "MIT",
		WorkingDir:  projectDir,
		Force:       true, // 跳过交互式提示
	}

	if err := client.Init(ctx, initOptions); err != nil {
		log.Printf("初始化项目失败: %v", err)
	} else {
		fmt.Println("项目初始化成功！")
	}

	// 5. 安装依赖包
	fmt.Println("\n5. 安装依赖包...")
	installOptions := npm.InstallOptions{
		WorkingDir: projectDir,
		SaveDev:    false,
	}

	packages := []string{"lodash", "axios"}
	for _, pkg := range packages {
		fmt.Printf("正在安装 %s...\n", pkg)
		if err := client.InstallPackage(ctx, pkg, installOptions); err != nil {
			log.Printf("安装 %s 失败: %v", pkg, err)
		} else {
			fmt.Printf("%s 安装成功！\n", pkg)
		}
	}

	// 6. 安装开发依赖
	fmt.Println("\n6. 安装开发依赖...")
	devInstallOptions := npm.InstallOptions{
		WorkingDir: projectDir,
		SaveDev:    true,
	}

	devPackages := []string{"jest", "eslint"}
	for _, pkg := range devPackages {
		fmt.Printf("正在安装开发依赖 %s...\n", pkg)
		if err := client.InstallPackage(ctx, pkg, devInstallOptions); err != nil {
			log.Printf("安装开发依赖 %s 失败: %v", pkg, err)
		} else {
			fmt.Printf("开发依赖 %s 安装成功！\n", pkg)
		}
	}

	// 7. 列出已安装的包
	fmt.Println("\n7. 列出已安装的包...")
	listOptions := npm.ListOptions{
		WorkingDir: projectDir,
		Depth:      0,
		JSON:       false,
	}

	packageList, err := client.ListPackages(ctx, listOptions)
	if err != nil {
		log.Printf("列出包失败: %v", err)
	} else {
		fmt.Printf("已安装 %d 个包:\n", len(packageList))
		for _, pkg := range packageList {
			fmt.Printf("  - %s@%s\n", pkg.Name, pkg.Version)
		}
	}

	// 8. 搜索包
	fmt.Println("\n8. 搜索包...")
	searchResults, err := client.Search(ctx, "react")
	if err != nil {
		log.Printf("搜索包失败: %v", err)
	} else {
		fmt.Printf("搜索 'react' 找到 %d 个结果:\n", len(searchResults))
		for i, result := range searchResults {
			if i >= 3 { // 只显示前3个结果
				break
			}
			fmt.Printf("  - %s@%s: %s\n",
				result.Package.Name,
				result.Package.Version,
				result.Package.Description)
		}
	}

	// 9. 获取包信息
	fmt.Println("\n9. 获取包信息...")
	packageInfo, err := client.GetPackageInfo(ctx, "lodash")
	if err != nil {
		log.Printf("获取包信息失败: %v", err)
	} else {
		fmt.Printf("包信息:\n")
		fmt.Printf("  名称: %s\n", packageInfo.Name)
		fmt.Printf("  版本: %s\n", packageInfo.Version)
		fmt.Printf("  描述: %s\n", packageInfo.Description)
		fmt.Printf("  许可证: %s\n", packageInfo.License)
		if packageInfo.Homepage != "" {
			fmt.Printf("  主页: %s\n", packageInfo.Homepage)
		}
	}

	// 10. 卸载包
	fmt.Println("\n10. 卸载包...")
	uninstallOptions := npm.UninstallOptions{
		WorkingDir: projectDir,
	}

	packageToUninstall := "axios"
	fmt.Printf("正在卸载 %s...\n", packageToUninstall)
	if err := client.UninstallPackage(ctx, packageToUninstall, uninstallOptions); err != nil {
		log.Printf("卸载 %s 失败: %v", packageToUninstall, err)
	} else {
		fmt.Printf("%s 卸载成功！\n", packageToUninstall)
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Printf("测试项目位置: %s\n", projectDir)
	fmt.Println("注意: 临时目录将在程序结束时自动清理")
}
