const zhCN: Record<string, string> = {
  // Status labels
  "Downloading": "下载中",
  "Paused": "已暂停",
  "Queued": "排队中",
  "Completed": "已完成",
  "Error": "出错",

  // Action tooltips
  "Pause": "暂停",
  "Resume": "恢复",
  "Open folder": "打开文件夹",
  "Open file": "打开文件",
  "Cancel": "取消",

  // Connection status
  "Connected": "已连接",
  "Offline": "离线",
  "Invalid": "无效",

  // View tabs
  "Active": "活动中",
  "History": "历史",

  // Empty states
  "No active downloads": "没有活跃的下载",
  "Downloads will appear here automatically": "下载将自动显示在此处",
  "No history downloads": "没有历史下载",
  "Completed downloads will appear here": "已完成的下载将显示在此处",

  // Duplicate modal
  "Duplicate Download": "重复下载",
  "This file is already being downloaded:": "此文件已在下载中：",
  "Skip": "跳过",
  "Download Anyway": "仍然下载",

  // Settings - labels
  "Intercept Downloads": "拦截下载",
  "Show Notifications": "显示通知",
  "Server": "服务器",
  "Server URL": "服务器 URL",
  "Auth Token": "身份验证令牌",

  // Settings - buttons
  "Save": "保存",
  "Delete": "删除",

  // Settings - status messages
  "Saving...": "保存中...",
  "Saved": "已保存",
  "Failed to save": "保存失败",
  "Removing...": "移除中...",
  "Removed": "已移除",
  "Failed to remove": "移除失败",

  // Settings - help & validation
  "Token can be obtained from TUI > Settings > Extension": "令牌可在 TUI > 设置 > 扩展 中获取",
  "Invalid Token": "无效的令牌",
  "Token is Required": "需要令牌",

  // Settings - placeholders
  "http://127.0.0.1:1700": "http://127.0.0.1:1700",
  "Enter your token": "输入您的令牌",

  // Settings - support
  "Support": "支持",
  "Report a Bug": "报告问题",
  "Firefox Add-ons": "Firefox 附加组件",

  // Background notifications
  "Server not running": "服务器未运行",
  "Download started: {filename}": "下载已开始：{filename}",
  "Failed to start download: {error}": "启动下载失败：{error}",
  "Pending download not found": "未找到待处理的下载",

  // History timestamp
  "Unknown time": "未知时间",
  "Today at ": "今天 ",
  "Yesterday at ": "昨天 ",

  // ETA
  "> 1 week": "> 1 周",
  "> 1 day": "> 1 天",

  // General
  "Unknown": "未知",
};

let isZH = false;
if (typeof navigator !== 'undefined') {
  isZH = navigator.language.startsWith('zh');
}
// For background script (no navigator), check browser.i18n
if (typeof browser !== 'undefined' && browser.i18n) {
  try {
    if (browser.i18n.getUILanguage().startsWith('zh')) isZH = true;
  } catch { /* ignore */ }
}

export function t(key: string): string {
  if (!isZH) return key;
  return zhCN[key] ?? key;
}
