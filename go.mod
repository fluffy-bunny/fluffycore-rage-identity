module github.com/fluffy-bunny/fluffycore-rage-identity

go 1.26.0

replace github.com/alexedwards/argon2id => github.com/ghstahl/argon2id v0.0.0-20260623181421-c6025230e214

//replace github.com/fluffy-bunny/fluffycore => ../fluffycore

//replace github.com/fluffy-bunny/viperEx => ../ViperEx

require (
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/alexedwards/argon2id v1.0.0
	github.com/coreos/go-oidc/v3 v3.20.0
	github.com/domodwyer/mailyak/v3 v3.6.2
	github.com/eko/gocache/lib/v4 v4.4.0
	github.com/eko/gocache/store/go_cache/v4 v4.2.10
	github.com/fluffy-bunny/fluffy-dozm-di v0.0.10
	github.com/fluffy-bunny/fluffycore v1.0.303
	github.com/go-oauth2/oauth2/v4 v4.5.4
	github.com/go-webauthn/webauthn v0.17.4
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gogo/status v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.30.0
	github.com/labstack/echo/v5 v5.3.1
	github.com/lestrrat-go/jwx/v3 v3.2.0
	github.com/nicksnyder/go-i18n/v2 v2.6.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/reugn/async v0.10.0
	github.com/rs/xid v1.6.0
	github.com/rs/zerolog v1.35.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.12.1
	github.com/sumup/aaguids-go v0.0.0-20240628162135-8d4ff0ce96a0
	github.com/swaggo/echo-swagger/v2 v2.0.1
	github.com/swaggo/swag v1.16.6
	github.com/wagslane/go-password-validator v0.3.0
	github.com/wamuir/svg-qr-code v0.0.0-20210725140500-9525ec975db7
	github.com/xlzd/gotp v0.1.0
	golang.org/x/oauth2 v0.37.0
	golang.org/x/text v0.41.0
	google.golang.org/genproto/googleapis/api v0.0.0-20260819154853-08b0e4226688
	google.golang.org/grpc v1.83.2
	google.golang.org/protobuf v1.36.12
	maragu.dev/gomponents v1.3.0
)

require (
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	github.com/DataDog/datadog-agent/comp/core/tagger/origindetection v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/opentelemetry-mapping-go/otlp/attributes v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/proto v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/template v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/trace v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/trace/log v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/trace/otel v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/trace/stats v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/trace/traceutil v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.77.3 // indirect
	github.com/DataDog/datadog-agent/pkg/version v0.77.3 // indirect
	github.com/DataDog/datadog-go/v5 v5.8.3 // indirect
	github.com/DataDog/dd-trace-go/contrib/google.golang.org/grpc/v2 v2.3.0 // indirect
	github.com/DataDog/dd-trace-go/instrumentation/testutils/grpc/v2 v2.3.0 // indirect
	github.com/DataDog/dd-trace-go/v2 v2.9.2 // indirect
	github.com/DataDog/go-libddwaf/v4 v4.9.0 // indirect
	github.com/DataDog/go-runtime-metrics-internal v0.0.4-0.20260217080614-b0f4edc38a6d // indirect
	github.com/DataDog/go-sqllexer v0.1.13 // indirect
	github.com/DataDog/go-tuf v1.1.1-0.5.2 // indirect
	github.com/DataDog/gostackparse v0.7.0 // indirect
	github.com/DataDog/sketches-go v1.4.8 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/auth0/go-jwt-middleware/v2 v2.3.1 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/ebitengine/purego v0.10.2 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/fgprof v0.9.5 // indirect
	github.com/fluffy-bunny/viperEx v0.0.40 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.2 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v1.0.0 // indirect
	github.com/go-openapi/jsonreference v1.0.0 // indirect
	github.com/go-openapi/spec v0.22.9 // indirect
	github.com/go-openapi/swag/conv v0.28.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.28.0 // indirect
	github.com/go-openapi/swag/loading v0.28.0 // indirect
	github.com/go-openapi/swag/pools v0.28.0 // indirect
	github.com/go-openapi/swag/stringutils v0.28.0 // indirect
	github.com/go-openapi/swag/typeutils v0.28.0 // indirect
	github.com/go-openapi/swag/yamlutils v0.28.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/go-webauthn/x v0.2.6 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang-migrate/migrate/v4 v4.19.1 // indirect
	github.com/golang/mock v1.7.0-rc.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/google/pprof v0.0.0-20250423184734-337e5dd93bb4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/hashicorp/go-version v1.8.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.14.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.8.3 // indirect
	github.com/jellydator/ttlcache/v2 v2.11.1 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/dsig v1.3.0 // indirect
	github.com/lestrrat-go/dsig-secp256k1 v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/httprc/v3 v3.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.1.7 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lestrrat-go/option/v2 v2.0.0 // indirect
	github.com/linkdata/deadlock v0.5.5 // indirect
	github.com/lufia/plan9stats v0.0.0-20260802145828-341c2f0c90b5 // indirect
	github.com/madflojo/tasks v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/maxence-charriere/go-app/v10 v10.1.11 // indirect
	github.com/mdaverde/jsonpath v0.2.1 // indirect
	github.com/minio/simdjson-go v0.4.5 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nats.go v1.53.1 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/outcaste-io/ristretto v0.2.3 // indirect
	github.com/pb33f/ordered-map/v2 v2.3.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/petermattis/goid v0.0.0-20260226131333-17d1149c6ac6 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/profile v1.7.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/power-devops/perfstat v0.0.0-20260805114148-88456608a4f6 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/quasoft/memstore v0.0.0-20191010062613-2bce066d2b0b // indirect
	github.com/richardartoul/molecule v1.0.1-0.20240531184615-7ca0df43c0b3 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.10.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shirou/gopsutil/v4 v4.26.7 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/sv-tools/openapi v0.4.0 // indirect
	github.com/swaggo/files/v2 v2.0.2 // indirect
	github.com/swaggo/swag/v2 v2.0.0-rc5 // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
	github.com/tklauser/go-sysconf v0.4.0 // indirect
	github.com/tklauser/numcpus v0.12.0 // indirect
	github.com/tkuchiki/go-timezone v0.2.2 // indirect
	github.com/tkuchiki/parsetime v0.3.0 // indirect
	github.com/trailofbits/go-mutexasserts v0.0.0-20250514102930-c1f3d2e37561 // indirect
	github.com/valyala/fastjson v1.6.10 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/collector/component v1.51.1-0.20260205185216-81bc641f26c0 // indirect
	go.opentelemetry.io/collector/featuregate v1.51.1-0.20260205185216-81bc641f26c0 // indirect
	go.opentelemetry.io/collector/pdata v1.51.1-0.20260205185216-81bc641f26c0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.145.1-0.20260205185216-81bc641f26c0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.71.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/host v0.71.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.71.0 // indirect
	go.opentelemetry.io/otel v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.46.0 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	go.opentelemetry.io/otel/sdk v1.46.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.46.0 // indirect
	go.opentelemetry.io/otel/trace v1.46.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	go.yaml.in/yaml/v4 v4.0.0-rc.2 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/exp v0.0.0-20260529124908-c761662dc8c9 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260825221802-da73d73af1c5 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.74.8 // indirect
	gopkg.in/go-jose/go-jose.v2 v2.6.3 // indirect
	gopkg.in/ini.v1 v1.67.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
