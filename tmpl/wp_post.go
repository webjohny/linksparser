package tmpl

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/dustin/go-humanize"
	"math/rand"
	"os"
	"strconv"
)

type WpPost struct {
	Title string
	Content string
	Url string
	AskedBy string
	Text string
	Links []*LinkResult
	CatId int
	Image string
	Date string
}

type LinkResult struct {
	Link string `json:"link"`
	Image []byte `json:"-"`
	Src string `json:"src"`
	Author string `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	GlobalRank int32 `json:"global_rank"`
	PageViews string `json:"page_views"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	CountryImg string `json:"country_img"`
}

func CreateWpPostTmpl(post WpPost) string {
	randInt := rand.Intn(90) + 9
	tmpl := `
		<style>
		.ca-list-item-title a{
			color: #005edf;
			opacity: 0.8;
		}
		.ca-list-item-title a:hover{
			opacity: 1;
		}
		.ca-list-item-title a:after{
			content: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAUCAYAAACNiR0NAAABN2lDQ1BBZG9iZSBSR0IgKDE5OTgpAAAokZWPv0rDUBSHvxtFxaFWCOLgcCdRUGzVwYxJW4ogWKtDkq1JQ5ViEm6uf/oQjm4dXNx9AidHwUHxCXwDxamDQ4QMBYvf9J3fORzOAaNi152GUYbzWKt205Gu58vZF2aYAoBOmKV2q3UAECdxxBjf7wiA10277jTG+38yH6ZKAyNguxtlIYgK0L/SqQYxBMygn2oQD4CpTto1EE9AqZf7G1AKcv8ASsr1fBBfgNlzPR+MOcAMcl8BTB1da4Bakg7UWe9Uy6plWdLuJkEkjweZjs4zuR+HiUoT1dFRF8jvA2AxH2w3HblWtay99X/+PRHX82Vun0cIQCw9F1lBeKEuf1UYO5PrYsdwGQ7vYXpUZLs3cLcBC7dFtlqF8hY8Dn8AwMZP/fNTP8gAAAAJcEhZcwAACxMAAAsTAQCanBgAAAXRaVRYdFhNTDpjb20uYWRvYmUueG1wAAAAAAA8P3hwYWNrZXQgYmVnaW49Iu+7vyIgaWQ9Ilc1TTBNcENlaGlIenJlU3pOVGN6a2M5ZCI/PiA8eDp4bXBtZXRhIHhtbG5zOng9ImFkb2JlOm5zOm1ldGEvIiB4OnhtcHRrPSJBZG9iZSBYTVAgQ29yZSA1LjYtYzE0NSA3OS4xNjM0OTksIDIwMTgvMDgvMTMtMTY6NDA6MjIgICAgICAgICI+IDxyZGY6UkRGIHhtbG5zOnJkZj0iaHR0cDovL3d3dy53My5vcmcvMTk5OS8wMi8yMi1yZGYtc3ludGF4LW5zIyI+IDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PSIiIHhtbG5zOnhtcD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLyIgeG1sbnM6ZGM9Imh0dHA6Ly9wdXJsLm9yZy9kYy9lbGVtZW50cy8xLjEvIiB4bWxuczpwaG90b3Nob3A9Imh0dHA6Ly9ucy5hZG9iZS5jb20vcGhvdG9zaG9wLzEuMC8iIHhtbG5zOnhtcE1NPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vIiB4bWxuczpzdEV2dD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL3NUeXBlL1Jlc291cmNlRXZlbnQjIiB4bXA6Q3JlYXRvclRvb2w9IkFkb2JlIFBob3Rvc2hvcCBDQyAyMDE5IChXaW5kb3dzKSIgeG1wOkNyZWF0ZURhdGU9IjIwMjEtMDQtMDVUMTk6MjI6MzMrMDM6MDAiIHhtcDpNb2RpZnlEYXRlPSIyMDIxLTA0LTA1VDE5OjI3OjIxKzAzOjAwIiB4bXA6TWV0YWRhdGFEYXRlPSIyMDIxLTA0LTA1VDE5OjI3OjIxKzAzOjAwIiBkYzpmb3JtYXQ9ImltYWdlL3BuZyIgcGhvdG9zaG9wOkNvbG9yTW9kZT0iMyIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDphYWZiMzA2OC1jZTk1LTQwNDgtYThhNC0wY2JmOGUyNWNlMzYiIHhtcE1NOkRvY3VtZW50SUQ9ImFkb2JlOmRvY2lkOnBob3Rvc2hvcDo0ZTQyMTYyYy1jYjVjLWU4NDEtODg5Mi05N2UzZmI5OGUyYjEiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1wLmRpZDpiZjcyNmEyNS0wNWI5LTc2NDEtODEwZS1hZGZjZDk3OGNjOWUiPiA8eG1wTU06SGlzdG9yeT4gPHJkZjpTZXE+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJjcmVhdGVkIiBzdEV2dDppbnN0YW5jZUlEPSJ4bXAuaWlkOmJmNzI2YTI1LTA1YjktNzY0MS04MTBlLWFkZmNkOTc4Y2M5ZSIgc3RFdnQ6d2hlbj0iMjAyMS0wNC0wNVQxOToyMjozMyswMzowMCIgc3RFdnQ6c29mdHdhcmVBZ2VudD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTkgKFdpbmRvd3MpIi8+IDxyZGY6bGkgc3RFdnQ6YWN0aW9uPSJzYXZlZCIgc3RFdnQ6aW5zdGFuY2VJRD0ieG1wLmlpZDphYWZiMzA2OC1jZTk1LTQwNDgtYThhNC0wY2JmOGUyNWNlMzYiIHN0RXZ0OndoZW49IjIwMjEtMDQtMDVUMTk6Mjc6MjErMDM6MDAiIHN0RXZ0OnNvZnR3YXJlQWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCBDQyAyMDE5IChXaW5kb3dzKSIgc3RFdnQ6Y2hhbmdlZD0iLyIvPiA8L3JkZjpTZXE+IDwveG1wTU06SGlzdG9yeT4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz42nF0EAAABUElEQVQ4ja3VsUscQRjG4ed0IZjAHSFY2FhIJDZpNWJnI2dlb+qA1yQQ0trYmBBCKiF/RGzPWNkIntgoRMJBhFSSoJ2dIWexczA3ubu9BV9Y9tt9Z3683yw7U+l0Ou5TlblGaWAV2xhHOvlXFj3UMFkAu8QDbAwa0AV+wasR0m1iC89wgKnEv8mwE2CfcFQAPAn3NhZxgbHI72Ty+Lt4O0LCWJ8TGNGL45KwFtZC3cTPUP/rrmG1BOwQ86H+hlU8wRVq2aBZQ5J1YXuoh/oaL1AvA0yT1RO/hdaowCMshLopb7Ov/vtKiaZxFsH2I9hrvCsLfIznUbKVyPuID+mEopZPQ6J1vEy8Nh4OAt4OgTbDleqvfIPoUbflSkHSfprAo37A3/LfrminibWEWXxNjQzL+I4/JRP+QKMf8BxP8QYzhrdfka9dC+9D3Tvgvo+AO9mHPpalCzRaAAAAAElFTkSuQmCC');
			margin: 0 3px 0 5px;
			width: 20px;
			height: 20px;
			display: inline-block;
		}
		</style>
		<!-- wp:html -->
		<div class="custom-article">
			<div class="ca-header">
				<div class="ca-header-content">
					<div class="ca-header-user">
						<figure style="background: url('/wp-content/uploads/2021/02/default.jpg') no-repeat center / cover;"></figure>
						<p class="ca-header-by">Asked by: ` + post.AskedBy + `</p>
						<div class="ca-header-tag red">Questioner</div>
						<div class="ca-header-tag blue">General</div>
					</div>
					<h1 class="ca-header-title">` + post.Title + `</h1>
					<p class="ca-header-description">` + post.Content + `</p>
					<div class="ca-header-date">
						<img src="/wp-content/uploads/2021/02/date.png" />
						<p>Last Updated: [post_published]</p>
					</div>
				</div>
				<div class="ca-header-vote">
					<a href="#" class="upvote">
						<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 292.362 292.361">
							<path d="M286.935 197.287L159.028 69.381c-3.613-3.617-7.895-5.424-12.847-5.424s-9.233 1.807-12.85 5.424L5.424 197.287C1.807 200.904 0 205.186 0 210.134s1.807 9.233 5.424 12.847c3.621 3.617 7.902 5.425 12.85 5.425h255.813c4.949 0 9.233-1.808 12.848-5.425 3.613-3.613 5.427-7.898 5.427-12.847s-1.814-9.23-5.427-12.847z"></path>
						</svg>
					</a>
					<p>` + strconv.Itoa(randInt) + `</p>
					<a href="#" class="downvote">
						<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 292.362 292.362">
							<path d="M286.935 69.377c-3.614-3.617-7.898-5.424-12.848-5.424H18.274c-4.952 0-9.233 1.807-12.85 5.424C1.807 72.998 0 77.279 0 82.228c0 4.948 1.807 9.229 5.424 12.847l127.907 127.907c3.621 3.617 7.902 5.428 12.85 5.428s9.233-1.811 12.847-5.428L286.935 95.074c3.613-3.617 5.427-7.898 5.427-12.847 0-4.948-1.814-9.229-5.427-12.85z"></path>
						</svg>
					</a>
				</div>
			</div>
			<div class="ca-text">
				` + post.Text + `
			</div>
			<div class="ca-list">`

	links := post.Links
	if len(links) > 0 {
		for _, link := range links {
			globalRank := humanize.Comma(int64(link.GlobalRank))
			tmpl += `
				<div class="ca-list-item">
					<div class="ca-list-item-row">
						<a class="ca-list-item-img" rel="nofollow" href="` + link.Link + `" target="_blank">
							<img style="max-width: 230px" src="` + link.Src + `" />
						</a>
						<div class="ca-list-item-content">
							<div class="ca-list-item-user">
								<figure style="background: url('/wp-content/uploads/2021/02/default.jpg') no-repeat center / cover;"></figure>
								<p class="ca-list-item-by">Asked by: ` + link.Author + `</p>
								<div class="ca-list-item-tag">Explainer</div>
							</div>
							<h3 class="ca-list-item-title">
								<a rel="nofollow" href="` + link.Link + `" target="_blank">` + link.Title + `</a>
							</h3>
							<p class="ca-list-item-link">` + link.Link + `</p>
							<p class="ca-list-item-text">` + link.Description + `</p>
						</div>
					</div>
					<div class="ca-list-item-info">
						<div class="ca-list-item-info-item">
							<p class="ca-list-item-info-item-value"><span style="background: #f0ad4e;">` + globalRank + `</span></p>
							<p class="ca-list-item-info-item-text">Global Rank</p>
						</div>
						<div class="ca-list-item-info-item">
							<p class="ca-list-item-info-item-value"><span style="background: #495b79;">` + link.PageViews + `</span></p>
							<p class="ca-list-item-info-item-text">Pageviews</p>
						</div>
						<div class="ca-list-item-info-item">
							<p class="ca-list-item-info-item-value">` + ConvertFlag(link.CountryImg) + ` <span style="color: #000;">` + link.CountryName + `</span></p>
							<p class="ca-list-item-info-item-text">Top Country</p>
						</div>
						<div class="ca-list-item-info-item">
							<p class="ca-list-item-info-item-value"><span style="background: #7ed321;">Up</span></p>
							<p class="ca-list-item-info-item-text">Site Status</p>
						</div>
						<div class="ca-list-item-info-item">
							<p class="ca-list-item-info-item-value"><span style="background: #5bc0de;">[ca_random_time]h ago</span></p>
							<p class="ca-list-item-info-item-text">Last Pinged</p>
						</div>
					</div>
				</div>
			`
		}
	}
				tmpl += `
			</div>
		</div>
		<!-- /wp:html -->
	`

	return tmpl
}


func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ConvertFlag(country string) string {
	if !Exists("./tmpl/png/" + country) {
		return ""
	}
	imgFile, err := os.Open("./tmpl/png/" + country) // a QR code image

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer imgFile.Close()

	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	img2html := "<img src=\"data:image/png;base64," + imgBase64Str + "\" />"

	return img2html
}
