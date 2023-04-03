# Busdes!
## 機能要件
* 立命館大学⇆南草津駅のバス時刻を取得して返す
* バスのリアルタイム接近情報を取得して返す



## アプリケーション
- [Android](https://play.google.com/store/apps/details?id=busdes.rits.jp&hl=ja&gl=US)
- [iOS](https://apps.apple.com/jp/app/busdes-%E3%83%90%E3%82%B9%E3%81%A7%E3%81%99/id1491015874)

## スクレイピング先
- [途中で停車するバス停はここから](https://ohmitetudo-bus.jorudan.biz/nsresult)
- [リアルタイム接近情報はここから](https://ohmitetudo-bus.jorudan.biz/busstatedtl)
- [時刻表はここから](https://ohmitetudo-bus.jorudan.biz/diagrampoledtl)
- [停留所の情報もここから頑張れば取れる](https://ohmitetudo-bus.jorudan.biz/diagrampoledtl)

## インフラ
- [GCP](https://cloud.google.com/)
- [CloudFlare](https://www.cloudflare.com/)
- [upstash](https://upstash.com/)