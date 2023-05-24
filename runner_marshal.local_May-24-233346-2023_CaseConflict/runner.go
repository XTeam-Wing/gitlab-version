package Runner

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	Silent      bool
	UrlList     goflags.StringSlice
	ShowCVE     bool
	Concurrency int
	Timeout     int
	Data        map[string]Version
	Output      string
	Debug       bool
}

const ruleAddr = "raw.githubusercontent.com/XTeam-Wing/gitlab-version/main/gitlab_hashes.json"

func (r *Runner) Run() (err error) {
	var wg sync.WaitGroup
	threads := make(chan struct{}, r.Concurrency)
	err = r.GetLatestHash()
	if err != nil {
		return err
	}
	var results []string
	for _, url := range r.UrlList {
		wg.Add(1)
		threads <- struct{}{}
		go func(url string) {
			defer wg.Done()
			body, err := r.GetBody(url)
			if err != nil {
				gologger.Error().Msgf("get %s failed: %v", url, err)
				return
			}
			build, result, err := r.Detect(body)
			if len(result) == 0 {
				gologger.Debug().Msgf("detect %s failed: %v", url, err)
				return
			}
			outputStr := fmt.Sprintf("target:%s version:%s build:%s", url, strings.Join(result, "||"), build)
			gologger.Info().Msgf(outputStr)
			results = append(results, outputStr)

			if err != nil {
				gologger.Error().Msgf("detect %s failed: %v", url, err)
			}
			<-threads
		}(url)
	}

	wg.Wait()

	if r.Output != "" {
		err := WriteFile(r.Output, []byte(strings.Join(results, "\n")), 0777)
		if err != nil {
			gologger.Error().Msgf("write file failed: %v", err)
		}
	}
	return err
}

func (r *Runner) GetLatestHash() (err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(r.Timeout) * time.Second}

	gologger.Debug().Msgf("get %s", ruleAddr)
	resp, err := client.Get(ruleAddr)
	if err != nil {
		gologger.Error().Msgf("get latest hash failed: %v", err)
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gologger.Error().Msgf("read latest hash failed: %v", err)
		return err
	}
	err = json.Unmarshal(body, &r.Data)
	if err != nil {
		gologger.Error().Msgf("decode latest hash failed: %v", err)
		return err
	}

	return nil
}

func (r *Runner) GetBody(url string) (result string, err error) {
	// 忽略https证书错误
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	gologger.Debug().Msgf("get %s", url)
	resp, err := client.Get(fmt.Sprintf("%s%s", url, "/assets/webpack/manifest.json"))

	defer resp.Body.Close()
	if err != nil {
		gologger.Error().Msgf("get %s failed: %v", url, err)
		return result, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		gologger.Error().Msgf("read body failed: %v", err)
		return result, err
	}
	return string(body), err
}
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
