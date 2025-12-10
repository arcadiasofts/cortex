package gateway

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHub),     // 허브 생성
	fx.Provide(NewHandler), // 핸들러 생성
	fx.Invoke(StartHub),    // 허브 실행 (Invoke는 앱 시작 시 자동 실행됨)
)

// 허브를 고루틴으로 실행시키는 헬퍼 함수
func StartHub(h *Hub) {
	go h.Run()
}
