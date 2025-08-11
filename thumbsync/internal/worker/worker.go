package worker

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"thumbsync/internal/cache"
	"thumbsync/internal/jobqueue"
	"time"

	"github.com/disintegration/imaging"
	"github.com/vbauerster/mpb/v8"
)

var dirOnce sync.Once

var httpTransport = &http.Transport{
	DisableKeepAlives:     true,
	TLSHandshakeTimeout:   5 * time.Second,
	ResponseHeaderTimeout: 10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	IdleConnTimeout:       5 * time.Second,
	MaxIdleConnsPerHost:   2,
}

var httpClient = &http.Client{
	Transport: httpTransport,
	Timeout:   30 * time.Second,
}

// Pool координирует группу воркеров.

type Pool struct {
	q      *jobqueue.Queue
	outDir string
	cache  *cache.Cache

	bar *mpb.Bar

	mu   sync.Mutex
	errs map[string]error

	wg sync.WaitGroup

	Dl   uint64 // downloaded (attempted)
	Ok   uint64 // successful
	Fail uint64 // failed
}

func NewPool(q *jobqueue.Queue, outDir string, bar *mpb.Bar) *Pool {
	return &Pool{
		q: q, outDir: outDir, bar: bar,
		cache: cache.New(),
		errs:  make(map[string]error),
	}
}

func (p *Pool) Run(n int) {
	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *Pool) Wait() { p.wg.Wait() }

func (p *Pool) Errors() map[string]error { return p.errs }

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		url, ok := p.q.Pop()
		if !ok {
			return
		}
		atomic.AddUint64(&p.Dl, 1)
		if _, hit := p.cache.Load(url); hit {
			atomic.AddUint64(&p.Ok, 1)
			p.bar.Increment()
			continue
		}
		img, err := download(url)
		if err != nil {
			p.recordErr(url, err)
			p.bar.Increment()
			continue
		}
		thumb := imaging.Thumbnail(img, 128, 128, imaging.Lanczos)
		buf, err := encodeJPEG(thumb)
		if err != nil {
			p.recordErr(url, err)
			p.bar.Increment()
			continue
		}
		p.cache.Store(url, buf)
		if err := p.save(url, buf); err != nil {
			p.recordErr(url, err)
			p.bar.Increment()
			continue
		}
		atomic.AddUint64(&p.Ok, 1)
		p.bar.Increment()
	}
}

func (p *Pool) recordErr(url string, err error) {
	atomic.AddUint64(&p.Fail, 1)
	p.mu.Lock()
	p.errs[url] = err
	p.mu.Unlock()
}

func (p *Pool) save(url string, data []byte) error {
	dirOnce.Do(func() { _ = os.MkdirAll(p.outDir, 0o755) })
	name := sha1Hash(url) + ".jpg"
	path := filepath.Join(p.outDir, name)
	return os.WriteFile(path, data, 0o644)
}

func download(url string) (image.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req) // <-- именно здесь выполняем запрос
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	const max = 10 << 20 // 10 MiB
	limited := io.LimitReader(resp.Body, max)

	img, _, err := image.Decode(limited)
	return img, err
}

func encodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := imaging.Encode(&buf, img, imaging.JPEG); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sha1Hash(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
