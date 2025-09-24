# translate-sidecar

This repo contains a small example CLI that calls the [Google Cloud Translate API](https://cloud.google.com/translate?hl=en) with [Sidecar](https://github.com/agentio/sidecar).

You'll need a Google Cloud project and [service credentials](https://cloud.google.com/iam/docs/service-account-creds), but once you have that, the rest is easy:

```
$ make
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
mkdir -p genproto
protoc proto/google/cloud/translate/v3/translation_service.proto \
--proto_path='proto' \
--go_opt='module=github.com/agentio/translate-sidecar/genproto' \
--go_opt=Mgoogle/cloud/translate/v3/translation_service.proto=github.com/agentio/translate-sidecar/genproto/translatepb \
--go_out='genproto'
go install ./...

$ translate-sidecar --parent projects/agentio "I just called a Google gRPC API" --credentials credentials.json
{
  "translations":  [
    {
      "translatedText":  "Acabo de llamar a una API gRPC de Google"
    }
  ]
}
```

In the example above, my project is `agentio` and my downloaded service account credentials are in a local file that I named `credentials.json`.
