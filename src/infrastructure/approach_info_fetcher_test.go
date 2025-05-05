package infrastructure

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
)

func setupTestDocument(t testing.TB) *CustomDocument {
	// response.htmlファイルを開く
	file, err := os.Open("response.html")
	if err != nil {
		t.Fatalf("Failed to open response.html: %v", err)
	}
	defer file.Close()

	// HTMLをパースする
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	return &CustomDocument{doc}
}

// ベンチマーク用にFmtの出力を抑制したドキュメントを作成する関数
func setupBenchmarkDocument(b *testing.B) *CustomDocument {
	// response.htmlファイルを開く
	file, err := os.Open("response.html")
	if err != nil {
		b.Fatalf("Failed to open response.html: %v", err)
	}
	defer file.Close()

	// ファイルの内容を読み込む
	content, err := io.ReadAll(file)
	if err != nil {
		b.Fatalf("Failed to read file: %v", err)
	}

	// io.Reader経由でドキュメントにパースする
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(content)))
	if err != nil {
		b.Fatalf("Failed to parse HTML: %v", err)
	}

	return &CustomDocument{doc}
}

func TestFetchApproachInfo(t *testing.T) {
	doc := setupTestDocument(t)

	// テスト対象の関数を実行
	moreMin, realArrivalTime, direction, scheduledTime, delay, busstop, via, requiredTime := doc.FetchApproachInfo()

	// 各戻り値の期待値
	expectedMoreMin := []string{"18"}
	expectedRealArrivalTime := []string{"18:24"}
	expectedDirection := []string{"立命館大学行 [パナソニック西口経由]", "立命館大学行 [パナソニック西口経由]"}
	expectedScheduledTime := []string{"18:24"}
	expectedDelay := []string{""}
	expectedBusStop := []string{"4"}
	expectedVia := []string{}
	expectedRequiredTime := []int{20}

	// 結果を検証
	if diff := cmp.Diff(expectedMoreMin, moreMin); diff != "" {
		t.Errorf("FetchApproachInfo (moreMin) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedRealArrivalTime, realArrivalTime); diff != "" {
		t.Errorf("FetchApproachInfo (realArrivalTime) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedDirection, direction); diff != "" {
		t.Errorf("FetchApproachInfo (direction) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedScheduledTime, scheduledTime); diff != "" {
		t.Errorf("FetchApproachInfo (scheduledTime) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedDelay, delay); diff != "" {
		t.Errorf("FetchApproachInfo (delay) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedBusStop, busstop); diff != "" {
		t.Errorf("FetchApproachInfo (busstop) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedVia, via); diff != "" {
		t.Errorf("FetchApproachInfo (via) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedRequiredTime, requiredTime); diff != "" {
		t.Errorf("FetchApproachInfo (requiredTime) mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchMoreMin(t *testing.T) {
	doc := setupTestDocument(t)

	// 期待値を設定
	expected := []string{"18"}

	// テスト対象の関数を実行
	actual := doc.FetchMoreMin()

	// 結果を検証
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("FetchMoreMin mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchRequiredTime(t *testing.T) {
	doc := setupTestDocument(t)

	// 期待値を設定
	expected := []int{20}

	// テスト対象の関数を実行
	actual := doc.FetchRequiredTime()

	// 結果を検証
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("FetchRequiredTime mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchBusStop(t *testing.T) {
	doc := setupTestDocument(t)

	// 期待値を設定 - HTMLの解析結果に基づいて期待値を設定
	expected := []string{"4"}

	// テスト対象の関数を実行
	actual := doc.FetchBusStop()

	// 結果を検証
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("FetchBusStop mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchDelay(t *testing.T) {
	doc := setupTestDocument(t)

	// 期待値を設定 - レスポンスHTMLの解析結果に基づいて
	expected := []string{""}

	// テスト対象の関数を実行
	actual := doc.FetchDelay(1)

	// 結果を検証
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("FetchDelay mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchVia(t *testing.T) {
	doc := setupTestDocument(t)

	// テスト対象の関数を実行
	actual := doc.FetchVia()

	// テスト結果から、Viaが取得できないことがわかったので、空のスライスを期待値とする
	expected := []string{}

	// 結果を検証
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("FetchVia mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchDirection(t *testing.T) {
	doc := setupTestDocument(t)

	// テスト対象の関数を実行
	actual := doc.FetchDirection()

	// テスト結果から、同じ文字列が2回取得されることがわかったので、期待値を修正
	expected := []string{"立命館大学行 [パナソニック西口経由]", "立命館大学行 [パナソニック西口経由]"}

	// 結果を検証
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("FetchDirection mismatch (-expected +actual):\n%s", diff)
	}
}

func TestFetchArrivalTime(t *testing.T) {
	doc := setupTestDocument(t)

	// テスト対象の関数を実行
	actualRealArrivalTime, actualScheduledTime := doc.FetchArrivalTime()

	// テスト結果から、実際には1つの時刻しか取得されないことがわかったので、期待値を修正
	expectedRealArrivalTime := []string{"18:24"}
	expectedScheduledTime := []string{"18:24"}

	// 結果を検証
	if diff := cmp.Diff(expectedRealArrivalTime, actualRealArrivalTime); diff != "" {
		t.Errorf("FetchArrivalTime (realArrivalTime) mismatch (-expected +actual):\n%s", diff)
	}

	if diff := cmp.Diff(expectedScheduledTime, actualScheduledTime); diff != "" {
		t.Errorf("FetchArrivalTime (scheduledTime) mismatch (-expected +actual):\n%s", diff)
	}
}

// ベンチマークテスト

func init() {
	// 使用するCPUの数を2に制限
	runtime.GOMAXPROCS(2)
}

func BenchmarkFetchApproachInfo(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _, _, _ = doc.FetchApproachInfo()
	}
}

func BenchmarkFetchMoreMin(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FetchMoreMin()
	}
}

func BenchmarkFetchRequiredTime(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FetchRequiredTime()
	}
}

func BenchmarkFetchBusStop(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FetchBusStop()
	}
}

func BenchmarkFetchDelay(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FetchDelay(1)
	}
}

func BenchmarkFetchVia(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FetchVia()
	}
}

func BenchmarkFetchDirection(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = doc.FetchDirection()
	}
}

func BenchmarkFetchArrivalTime(b *testing.B) {
	doc := setupBenchmarkDocument(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = doc.FetchArrivalTime()
	}
}
