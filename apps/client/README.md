# Client

The client requires the OTel collector binary. You need to download it to the `bin` directory:

```
mkdir bin
cd bin
curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.92.0/otelcol-contrib_0.92.0_darwin_amd64.tar.gz
tar -xvf otelcol-contrib_0.92.0_darwin_amd64.tar.gz
```
