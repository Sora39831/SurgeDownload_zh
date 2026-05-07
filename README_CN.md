<div align="center">

# Surge

**用 Go 构建的极速 TUI 下载管理器，为高级用户而生**

[English](README.md) • 中文

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/SurgeDM/Surge)
[![Release](https://img.shields.io/github/v/release/SurgeDM/Surge?style=flat-square&color=blue)](https://github.com/SurgeDM/Surge/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/SurgeDM/Surge?style=flat-square&color=cyan)](go.mod)
[![License](https://img.shields.io/badge/License-MIT-grey.svg?style=flat-square)](LICENSE)
[![BuyMeACoffee](https://raw.githubusercontent.com/pachadotdev/buymeacoffee-badges/main/bmc-violet.svg)](https://www.buymeacoffee.com/surge.downloader)
[![Stars](https://img.shields.io/github/stars/SurgeDM/Surge?style=social)](https://github.com/SurgeDM/Surge/stargazers)

[安装](#安装) • [使用方法](#使用方法) • [主题](docs/THEMES.md) • [字体](docs/FONTS.md) • [性能测试](#性能测试) • [浏览器扩展](#浏览器扩展) • [设置](docs/SETTINGS.md) • [CLI 参考](docs/USAGE.md)

</div>

---

## 什么是 Surge？

Surge 专为偏爱键盘驱动工作流的高级用户设计。它拥有精美的**终端用户界面 (TUI)**，同时支持后台**无头服务端**和**命令行工具 (CLI)**，方便自动化操作。

![Surge 演示](assets/demo.gif)

---

## 为什么选择 Surge？

大多数浏览器单次下载只能建立一个连接。Surge 可以同时建立多个连接（最多 32 个），将文件分片并行下载。但我们不止于此：

- **极速性能：** 旨在最大化带宽利用率，以最快速度完成下载。
- **多镜像源：** 支持同时从多个源下载。Surge 会在所有可用镜像之间分配工作线程，并自动处理故障转移。
- **顺序下载：** 可选择严格按顺序下载文件（流模式），适合希望在下载过程中即时预览的媒体文件。
- **守护进程架构：** Surge 运行单一后台"引擎"。你可以打开 10 个不同的终端标签页并排队下载；它们全部汇入一个高效的管理器。
- **精美 TUI：** 基于 Bubble Tea 和 Lipgloss 构建，支持自定义调色板和完整的主题引擎。

关于我们如何让下载更快的深入介绍（如任务窃取和慢速工作线程处理），请参阅我们的 **[优化指南](docs/OPTIMIZATIONS.md)**。

---

## 支持项目

我们只是两名计算机专业的学生，在课程和考试之余构建 Surge。我们热爱这个项目，但维护如此规模的项目需要时间和资源。这就是需要你支持的地方！

如果 Surge 为你节省了时间，请考虑支持我们的开发！捐赠将直接用于：

- **发布扩展：** 支付 Chrome Web Store 费用，让你终于可以正式安装扩展（不再需要侧加载！）。
- **开发工具：** 采购 **GoReleaser Pro** 等工具的许可证，帮助我们自动化构建。
- **Debrid 集成：** 承担订阅费用，以便我们测试和构建原生的 Debrid 支持。

[**☕ 请我们喝杯咖啡**](https://www.buymeacoffee.com/surge.downloader)

_完全自愿——你的 star、issue 和贡献已经让我们非常感激！:)_

---

## 安装

Surge 支持多个平台，选择最适合你的安装方式。

| 平台 / 方式                         | 命令 / 说明                                                                         | 备注                                 |
| :---------------------------------- | :---------------------------------------------------------------------------------- | :----------------------------------- |
| **预编译二进制文件**                | [从 Releases 下载](https://github.com/SurgeDM/Surge/releases/latest)                | 最简单的方法，下载即可运行。         |
| **Arch Linux (AUR)**                | `yay -S surge`                                                                      | 通过 AUR 管理。                      |
| **macOS / Linux (Homebrew)**        | `brew install SurgeDM/tap/surge`                                                    | 推荐 Mac/Linux 用户使用。            |
| **Windows**                         | `winget install surge-downloader.surge`<br />或<br />`scoop install surge`          | 推荐 Windows 用户使用。              |
| **Dockerfile**                      | [参见说明](#4-使用-docker-compose-的服务器模式)                                     | 使用 Docker Compose 以服务端模式运行 |
| **Go Install**                      | `go install github.com/SurgeDM/Surge@latest`                                        | 需要 Go 1.25+                        |

---

## 使用方法

Surge 有两种主要模式：**TUI（交互式）** 和 **Server（无头服务端）**。

完整参考请参阅 **[主题指南](docs/THEMES.md)**、**[设置与配置指南](docs/SETTINGS.md)** 和 **[CLI 使用指南](docs/USAGE.md)**。

### 1. 交互式 TUI 模式

直接运行 `surge` 进入仪表盘。在这里你可以可视化进度、管理队列、查看速度图表。如果遇到任何问题，按 `?` 打开问题反馈向导。

```bash
# 启动 TUI
surge

# 启动 TUI 并排队下载
surge https://example.com/file1.zip https://example.com/file2.zip

# 结合 URL 和批量文件
surge https://example.com/file.zip --batch urls.txt
```

### 2. 服务器模式（无头）

适用于服务器、树莓派或后台进程。

```bash
# 启动服务器
surge server

# 启动服务器并添加下载
surge server https://url.com/file.zip

# 使用显式 API Token 启动
surge server --token <token>
```

`surge` 和 `surge server` 默认将 HTTP API 绑定到 `0.0.0.0`（所有网络接口）。
这意味着服务器可通过 `localhost`（127.0.0.1）以及你的局域网 IP 访问。

API 受 token 保护。运行以下命令生成/查看你的 token：

```bash
surge token
```

或者，你也可以在 TUI 中的 **设置 > 扩展** 中找到它。

### 3. 远程 TUI

连接到正在运行的 Surge 守护进程（本地或远程）。

```bash
# 连接到本地服务器（自动检测）
surge connect

# 连接到远程守护进程
surge connect 192.168.1.10:1700 --token <token>

# 等效的全局标志形式
surge --host 192.168.1.10:1700 --token <token>
```

默认情况下，`surge connect` 使用：

- `http://` 用于回环地址和私有 IP 目标
- `https://` 用于公网/主机名目标

### 4. 全局连接标志（CLI + TUI）

以下全局标志适用于所有命令：

- `--host <host:port>`：TUI 和 CLI 操作的目标服务器。
- `--token <token>`：用于认证的 Bearer Token。

环境变量回退：

- `SURGE_HOST`
- `SURGE_TOKEN`

### 5. 使用 Docker Compose 的服务器模式

下载 compose 文件并启动容器：

```bash
wget https://raw.githubusercontent.com/SurgeDM/Surge/refs/heads/main/docker/compose.yml
docker compose up -d
```

获取 API Token：

```bash
docker compose exec surge surge token
```

保存此 token——后续 API 请求认证和远程连接都需要用到它。

检查下载/API 是否可用：

```bash
docker compose exec surge surge ls
```

查看日志：

```bash
docker compose logs -f surge
```

---

## 字体

Surge 在 TUI 中内置了 Nerd Font，但实际字体选择由你的终端控制。
有关安装步骤和许可详情，请参阅 [docs/FONTS.md](docs/FONTS.md)。

---

## 性能测试

我们对 Surge 与常用工具进行了对比测试。得益于连接优化逻辑，Surge 的性能显著优于单连接工具。

| 工具              | 耗时              | 速度                  | 对比           |
| ----------------- | ----------------- | --------------------- | -------------- |
| **Surge**         | **28.93 秒**      | **35.40 MB/s**        | **—**          |
| aria2c            | 40.04 秒          | 25.57 MB/s            | 慢 1.38 倍     |
| curl              | 57.57 秒          | 17.79 MB/s            | 慢 1.99 倍     |
| wget              | 61.81 秒          | 16.57 MB/s            | 慢 2.14 倍     |

> _测试详情：1GB 文件，Windows 11，Ryzen 5 5600X，360 Mbps 网络。结果取 5 次平均值。_

欢迎你在自己的系统上对 Surge 进行性能测试！

---

## 浏览器扩展

Surge 扩展可以拦截浏览器下载并将它们直接发送到你的终端。它默认在 **1700** 端口与 Surge 客户端通信。

> [!IMPORTANT]
> 连接扩展与 Surge 服务器需要**认证 Token**。可在 TUI 的 **设置 > 扩展** 中获取，或运行 `surge token`。

### Chrome / Edge / Brave

1. 从最新的 GitHub Release 中下载 `extension-chrome.zip`。
2. 将其解压到磁盘的某个位置。
3. 打开浏览器，导航至 `chrome://extensions`。
4. 在右上角开启**「开发者模式」**。
5. 点击**「加载已解压的扩展程序」**。
6. 选择解压后的 `extension-chrome` 文件夹。
7. 点击浏览器工具栏中的 Surge 图标，在设置中输入你的**认证 Token**。

### Firefox

1. **稳定版：** [获取插件](https://addons.mozilla.org/en-US/firefox/addon/surge/)
2. **开发版：**
   - 从最新的 GitHub Release 中下载 `extension-firefox.zip`。
   - 导航至 `about:debugging#/runtime/this-firefox`。
   - 点击**「临时加载附加组件…」**。
   - 选择 zip 文件（或解压后选择 `manifest.json`）。
   - 点击浏览器工具栏中的 Surge 图标，在设置中输入你的**认证 Token**。

---

## 致谢

衷心感谢帮助我们构建和发布 Surge 的团队和赞助商：

- [Charm](https://charm.sh/) 提供了卓越的终端 UI 生态（Bubble Tea、Lip Gloss 等）。
- [GoReleaser Pro](https://goreleaser.com/pro/) 用于发布自动化（为开源项目免费提供）。

---

## 社区与贡献

我们热爱社区贡献！无论是 bug 修复、新功能，还是修正拼写错误，
随时欢迎提交 PR。简要指南请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)。

欢迎在 [Discussions](https://github.com/SurgeDM/Surge/discussions) 中提出任何问题或想法，也可以在 [X (Twitter)](https://x.com/SurgeDownloader) 上关注我们！

## 许可证

基于 MIT 许可证发布。更多信息请参阅 [LICENSE](https://github.com/SurgeDM/Surge/blob/main/LICENSE)。

---

<div align="center">
<a href="https://star-history.com/#SurgeDM/Surge&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=SurgeDM/Surge&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=SurgeDM/Surge&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=SurgeDM/Surge&type=Date" />
 </picture>
</a>

<br />
如果 Surge 帮你省了时间，不妨给它点个 ⭐，帮助更多人发现它！
</div>
