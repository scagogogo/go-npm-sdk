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

	fmt.Println("=== Go NPM SDK 便携版使用示例 ===")

	// 1. 创建便携版管理器
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("获取用户目录失败: %v", err)
	}

	portableDir := filepath.Join(homeDir, ".go-npm-sdk-example", "portable")
	portableManager, err := npm.NewPortableManager(portableDir)
	if err != nil {
		log.Fatalf("创建便携版管理器失败: %v", err)
	}

	fmt.Printf("便携版目录: %s\n", portableDir)

	// 2. 列出已安装的版本
	fmt.Println("\n2. 列出已安装的版本...")
	configs, err := portableManager.List()
	if err != nil {
		log.Printf("列出版本失败: %v", err)
	} else {
		if len(configs) == 0 {
			fmt.Println("没有已安装的便携版")
		} else {
			fmt.Printf("已安装 %d 个版本:\n", len(configs))
			for _, config := range configs {
				fmt.Printf("  - v%s (路径: %s)\n", config.Version, config.InstallPath)
			}
		}
	}

	// 3. 安装便携版Node.js
	fmt.Println("\n3. 安装便携版Node.js...")
	version := "18.17.0" // 指定版本

	// 检查是否已安装
	if portableManager.IsVersionInstalled(version) {
		fmt.Printf("版本 %s 已安装\n", version)
	} else {
		fmt.Printf("正在安装版本 %s...\n", version)
		
		// 进度回调
		progress := func(message string) {
			fmt.Printf("  %s\n", message)
		}

		config, err := portableManager.Install(ctx, version, progress)
		if err != nil {
			log.Fatalf("安装便携版失败: %v", err)
		}

		fmt.Printf("安装成功！\n")
		fmt.Printf("  版本: %s\n", config.Version)
		fmt.Printf("  Node.js路径: %s\n", config.NodePath)
		fmt.Printf("  npm路径: %s\n", config.NpmPath)
	}

	// 4. 获取配置信息
	fmt.Println("\n4. 获取配置信息...")
	config, err := portableManager.GetConfig(version)
	if err != nil {
		log.Printf("获取配置失败: %v", err)
	} else {
		fmt.Printf("版本 %s 配置:\n", version)
		fmt.Printf("  安装路径: %s\n", config.InstallPath)
		fmt.Printf("  Node.js路径: %s\n", config.NodePath)
		fmt.Printf("  npm路径: %s\n", config.NpmPath)
		fmt.Printf("  安装时间: %s\n", config.InstallDate)
	}

	// 5. 使用便携版创建npm客户端
	fmt.Println("\n5. 使用便携版创建npm客户端...")
	client, err := portableManager.CreateClient(version)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	// 6. 测试npm功能
	fmt.Println("\n6. 测试npm功能...")
	
	// 检查npm是否可用
	if !client.IsAvailable(ctx) {
		log.Printf("便携版npm不可用")
		return
	}

	// 获取版本
	npmVersion, err := client.Version(ctx)
	if err != nil {
		log.Printf("获取npm版本失败: %v", err)
	} else {
		fmt.Printf("npm版本: %s\n", npmVersion)
	}

	// 7. 创建测试项目
	fmt.Println("\n7. 使用便携版创建测试项目...")
	projectDir := filepath.Join(os.TempDir(), "portable-npm-test")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Fatalf("创建项目目录失败: %v", err)
	}
	defer os.RemoveAll(projectDir) // 清理

	// 初始化项目
	initOptions := npm.InitOptions{
		Name:        "portable-test",
		Version:     "1.0.0",
		Description: "便携版npm测试项目",
		WorkingDir:  projectDir,
		Force:       true,
	}

	if err := client.Init(ctx, initOptions); err != nil {
		log.Printf("初始化项目失败: %v", err)
	} else {
		fmt.Println("项目初始化成功！")
	}

	// 安装一个简单的包
	installOptions := npm.InstallOptions{
		WorkingDir: projectDir,
	}

	testPackage := "lodash"
	fmt.Printf("正在安装 %s...\n", testPackage)
	if err := client.InstallPackage(ctx, testPackage, installOptions); err != nil {
		log.Printf("安装 %s 失败: %v", testPackage, err)
	} else {
		fmt.Printf("%s 安装成功！\n", testPackage)
	}

	// 8. 设置为默认版本
	fmt.Println("\n8. 设置为默认版本...")
	if err := portableManager.SetAsDefault(version); err != nil {
		log.Printf("设置默认版本失败: %v", err)
	} else {
		fmt.Printf("版本 %s 已设置为默认版本\n", version)
		fmt.Printf("默认路径: %s\n", portableManager.GetDefaultPath())
	}

	// 9. 演示多版本管理
	fmt.Println("\n9. 多版本管理演示...")
	
	// 尝试安装另一个版本（如果需要）
	anotherVersion := "16.20.0"
	fmt.Printf("检查版本 %s...\n", anotherVersion)
	
	if !portableManager.IsVersionInstalled(anotherVersion) {
		fmt.Printf("版本 %s 未安装，可以使用以下代码安装:\n", anotherVersion)
		fmt.Printf("  config, err := portableManager.Install(ctx, \"%s\", progress)\n", anotherVersion)
	} else {
		fmt.Printf("版本 %s 已安装\n", anotherVersion)
		
		// 创建另一个版本的客户端
		anotherClient, err := portableManager.CreateClient(anotherVersion)
		if err != nil {
			log.Printf("创建另一个版本的客户端失败: %v", err)
		} else {
			if anotherVersion, err := anotherClient.Version(ctx); err == nil {
				fmt.Printf("另一个版本的npm: %s\n", anotherVersion)
			}
		}
	}

	// 10. 清理演示
	fmt.Println("\n10. 清理选项...")
	fmt.Println("如果要卸载便携版，可以使用:")
	fmt.Printf("  err := portableManager.Uninstall(\"%s\")\n", version)
	fmt.Println("注意: 本示例不会自动清理便携版安装")

	fmt.Println("\n=== 便携版示例完成 ===")
	fmt.Printf("便携版安装目录: %s\n", portableDir)
	fmt.Printf("测试项目目录: %s (将自动清理)\n", projectDir)
}
