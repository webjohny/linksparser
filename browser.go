package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/performance"
	"github.com/chromedp/cdproto/security"
	"github.com/webjohny/chromedp"
	"linksparser/services"
	"log"
	"strings"
	"time"
)

type Browser struct {
	interceptionID fetch.RequestID

	BrowserContextID cdp.BrowserContextID

	cancelBrowser context.CancelFunc
	cancelTask context.CancelFunc

	Proxy *Proxy
	ctx context.Context

	isOpened bool
	streamId int
	limit int64
}

func (b *Browser) Init() bool {
	if b.isOpened {
		return true
	}

	if b.limit < 1 {
		b.limit = 60
	}

	if b.Proxy == nil || b.Proxy.Host == "" {
		fmt.Println("CHECK PROXY")
		b.Proxy = NewProxy()
		// Подключаемся к прокси
		if b.Proxy == nil {
			return false
		}

		b.Proxy.setTimeout(b.streamId, 5)

		if !b.checkProxy(b.Proxy) {
			return false
		}
	}

	if b.ctx == nil {
		fmt.Println("NEW INSTANCE")

		options := b.setOpts(b.Proxy)

		if CONF.Env == "local" {
			options = append(options, chromedp.Flag("headless", false))
		}
		// Запускаем контекст браузера
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
		b.cancelBrowser = cancel

		// Устанавливаем собственный logger
		//, chromedp.WithDebugf(log.Printf)
		taskCtx, cancel := chromedp.NewContext(allocCtx)
		b.cancelTask = cancel
		b.ctx = taskCtx

		if err := chromedp.Run(taskCtx,
			chromedp.Sleep(time.Second),
			b.setProxyToContext(b.Proxy),
		); err != nil {
			log.Println("Browser.Init.HasError", err)
			return false
		}
	}

	fmt.Println("RETURN TRUE")

	b.isOpened = true
	return true
}

func (b *Browser) checkProxy(proxy *Proxy) bool {
	if proxy == nil {
		return false
	}

	fmt.Println(proxy.LocalIp)

	options := b.setOpts(proxy)
	//@toDo убрать коммент
	if CONF.Env == "local" {
		options = append(options, chromedp.Flag("headless", false))
	}

	keyWords := []string{
		"whats+my+ip",
		"ssh+run+command",
		"how+work+with+git",
		"bitcoin+price+2013+year",
		"онлайн+обменник+крипта+рубль",
		"где+купить+акции",
		"i+want+to+spend+crypto",
	}

	// Запускаем контекст браузера
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	b.cancelBrowser = cancel

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	b.cancelTask = cancelTask
	b.ctx = taskCtx

	var searchHtml string

	fmt.Println(services.ArrayRand(keyWords))
	if err := chromedp.Run(taskCtx,
		b.setProxyToContext(proxy),
		b.runWithTimeOut(10, false, chromedp.Tasks{
			// Устанавливаем страницу для парсинга
			chromedp.Navigate("https://www.google.com/search?q=" + services.ArrayRand(keyWords)),
			//chromedp.Navigate("https://deelay.me/23545/google.com"),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.OuterHTML("body", &searchHtml, chromedp.ByQuery),
		}),
	); err != nil {
		log.Println("Browser.checkProxy.HasError", err)
		return false
	}

	if searchHtml != "" && !b.CheckCaptcha(searchHtml) {
		fmt.Println("YES)))")
		return true
	}

	return false
}


func (b *Browser) CheckCaptcha(html string) bool {
	return strings.Contains(html,"g-recaptcha") && strings.Contains(html,"data-sitekey")
}

func (b *Browser) setProxyToContext(proxy *Proxy) chromedp.Tasks {
	fmt.Print(proxy.Login, proxy.Password)
	if CONF.Env != "local" {
		return chromedp.Tasks{
			network.Enable(),
			performance.Enable(),
			page.SetLifecycleEventsEnabled(true),
			security.SetIgnoreCertificateErrors(true),
			emulation.SetTouchEmulationEnabled(false),
			network.SetCacheDisabled(true),
			chromedp.Authentication(proxy.Login, proxy.Password),
		}
	}else{
		return chromedp.Tasks{}
	}
}

func (b *Browser) setOpts(proxy *Proxy) []chromedp.ExecAllocatorOption {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("password-store", false),
	)

	if proxy != nil && CONF.Env != "local" {
		proxyScheme := proxy.LocalIp

		if proxyScheme != "" {
			opts = append(opts, chromedp.ProxyServer(proxyScheme))
		}

		if proxy.Agent != "" {
			opts = append(opts, chromedp.UserAgent(proxy.Agent))
		}
	}
	return opts
}

func (b *Browser) runWithTimeOut(timeout time.Duration, isStrict bool, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var check bool
		time.AfterFunc(timeout * time.Second, func(){
			if !check {
				if isStrict {
					b.cancelTask()
				} else {
					b.Reload()
				}
			}
		})

		err := tasks.Do(ctx)
		if err != nil {
			if "page load error net::ERR_ABORTED" != err.Error() {
				fmt.Println("ERR.Browser.runWithTimeOut.1", err)
				b.cancelTask()
				return err
			}else{
				time.Sleep(2 * time.Second)
			}
		}
		check = true
		fmt.Println("RUN_WITH_TIMEOUT")
		return nil
	}
}

func (b *Browser) Cancel() {
	if b.cancelTask != nil {
		b.cancelTask()
	}

	if b.isOpened {
		if b.Proxy != nil && b.Proxy.LocalIp != "" {
			b.Proxy.freeProxy()
		}
		b.isOpened = false
	}
	b.ctx = nil
}

func (b *Browser) Reload() bool {
	b.Cancel()
	time.Sleep(time.Second)
	return b.Init()
}

func (b *Browser) ChangeTab() {

}

func (b *Browser) ScreenShot(url string) (*[]byte, error) {
	if url == "" {
		return nil, errors.New("undefined url")
	}

	var buf []byte
	if err := chromedp.Run(b.ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctxt context.Context) error {
			fmt.Println("FIRST")
			_, viewLayout, contentRect, err := page.GetLayoutMetrics().Do(ctxt)
			if err != nil {
				return err
			}

			v := page.Viewport{
				X:      contentRect.X,
				Y:      contentRect.Y,
				Width:  viewLayout.ClientWidth, // or contentRect.Width,
				Height: viewLayout.ClientHeight,
				Scale:  1,
			}
			log.Printf("Capture %#v", v)
			buf, err = page.CaptureScreenshot().WithClip(&v).Do(ctxt)
			if err != nil {
				return err
			}
			return nil
		}),
	); err != nil {
		log.Println("Browser.ScreenShotSave.HasError", err)
		return nil, err
	}

	if buf != nil && len(buf) < 1 {
		return nil, errors.New("undefined image")
	}

	return &buf, nil
}
