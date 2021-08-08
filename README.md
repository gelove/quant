# 接口文档
[binance-docs](https://binance-docs.github.io/apidocs/spot/cn/#ed913b7357)

# 主网和测试网
| Spot API URL                         | Spot Test Network URL               |
| ------------------------------------ | ----------------------------------- |
| https://api.binance.com/api          | https://testnet.binance.vision/api  |
| wss://stream.binance.com:9443/ws     | wss://testnet.binance.vision/ws     |
| wss://stream.binance.com:9443/stream | wss://testnet.binance.vision/stream |

# 使用
## 获取最佳反弹
make bounce MOMENT_MAC=2021-06-22T16:00:00

MOMENT_MAC 指定上次大跌最低的具体时间

## 获取最优振幅
make swing [COUNT=30]

## 运行网格交易
make start
