package sdk

import (
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/analyze"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/manage"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"

	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/analyze/v1"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/manage/v1"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/rest"
	_ "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket"
)
