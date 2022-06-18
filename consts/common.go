package consts

const (
	// ゲームを強制削除するまでの時間(秒): ゲーム作成から4時間
	TimeAutoDelete = 14400

	MutexTTL        = 300 // mutexTTL(ミリ秒)
	MutexRetryCount = 3   // mutex最大リトライ回数
	MutexRetrySpan  = 100 // mutexリトライ感覚(ミリ秒)
)
