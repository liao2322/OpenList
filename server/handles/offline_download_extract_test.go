package handles

import "testing"

func TestCollectOfflineURLsExtractMagnetFromRawText(t *testing.T) {
	req := AddOfflineDownloadReq{
		RawText: "first magnet:?xt=urn:btih:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA&dn=a second magnet:?xt=urn:btih:BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
	}
	urls := collectOfflineURLs(req)
	if len(urls) != 2 {
		t.Fatalf("expected 2 urls, got %d: %#v", len(urls), urls)
	}
}

func TestCollectOfflineURLsDeduplicate(t *testing.T) {
	m := "magnet:?xt=urn:btih:CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC&dn=c"
	req := AddOfflineDownloadReq{
		Urls:    []string{m, "  " + m + "  "},
		RawText: m,
	}
	urls := collectOfflineURLs(req)
	if len(urls) != 1 {
		t.Fatalf("expected 1 url after dedup, got %d: %#v", len(urls), urls)
	}
}

func TestCollectOfflineURLsKeepNormalURL(t *testing.T) {
	req := AddOfflineDownloadReq{Urls: []string{"https://example.com/file.iso"}}
	urls := collectOfflineURLs(req)
	if len(urls) != 1 || urls[0] != "https://example.com/file.iso" {
		t.Fatalf("expected normal url kept, got %#v", urls)
	}
}
