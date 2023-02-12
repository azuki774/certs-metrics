# certs-metrics

証明書 `(ca.crt)` の残り有効期限を表す Prometheus metrics を出力するだけのプログラム

### Usage
証明書ファイルは引数で複数指定可

- 証明書の読み込みが正常にできるかを確認するコマンド
    ```
    $ ./build/bin/certs-metrics check <ca.crt>

    Ex:
    $ ./certs-metrics check ca.crt /root/ca2.crt
    ```

- 証明書の読み込みのメトリクスサーバを起動するコマンド
    ```
    $ ./build/bin/certs-metrics start <ca.crt>

    Ex:
    $ ./certs-metrics start ca.crt /root/ca2.crt
    ```


### metrics
|  metrics name  |  pattern  | Description
| ---- | ---- | ---- |
|  cert_not_before  |  Gauge  | 証明書の `NotBefore` を UnixTime で表示 |
|  cert_not_after  |  Gauge  | 証明書の `NotAfter` を UnixTime で表示 |
|  cert_valid_period  |  Gauge  | 証明書の有効期限の長さを分で表示|
|  cert_remaining_valid_period  |  Gauge  | 証明書の残り有効期限を分で表示|
