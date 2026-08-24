# semver-resolver

语义化版本（SemVer）约束求解 CLI（纯 Go 标准库，离线可用）。

给定一组版本约束与一个候选版本列表，求解满足约束的最高版本；或判断单个版本是否满足约束。

支持的约束语法：

- `*`：任意版本
- `1.2.3`：精确等于
- `^1.2.3`：兼容区间（>=1.2.3, <2.0.0；主版本为 0 时按次/修订号收窄）
- `~1.2.3`：近似区间（>=1.2.3, <1.3.0）
- `>=1.0.0` `<2.0.0`：比较器（逗号或空格分隔表示「且」）
- `A || B`：两组约束的「或」

## 用法

```bash
# 从候选版本中求解满足约束的最高版本
semver-resolver resolve --constraint "^1.2.3" --versions 1.2.0,1.2.3,1.9.0,2.0.0

# 判断单个版本是否满足约束
semver-resolver check --constraint ">=1.0.0,<2.0.0" --version 1.9.0

# 指定输出格式
semver-resolver resolve --constraint "^1.2.3" --versions 1.2.3,1.9.0 --format json
```

`example/versions.txt` 提供了一个候选版本示例清单。
