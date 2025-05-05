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

## Cloud Trace の設定

このアプリケーションは Google Cloud Trace と統合されています。Cloud Run にデプロイする際に、以下の手順で設定してください：

1. サービスアカウントに必要な権限を付与する
```bash
# サービスアカウントに Cloud Trace 権限を付与
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
    --member=serviceAccount:YOUR_SERVICE_ACCOUNT_EMAIL \
    --role=roles/cloudtrace.agent
```

2. Cloud Run サービスをデプロイする際に、サービスアカウントを指定する
```bash
gcloud run deploy bustimer \
    --source . \
    --platform managed \
    --allow-unauthenticated \
    --service-account=YOUR_SERVICE_ACCOUNT_EMAIL \
    --region=asia-northeast1
```

3. デプロイ後、Google Cloud Console の "Cloud Trace" ページでトレースデータを確認できます。