package tmpl

import (
	"github.com/dustin/go-humanize"
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
	Link string
	Image []byte
	Src string
	Author string
	Title string
	Description string
	GlobalRank int32
	PageViews string
	CountryCode string
	CountryName string
	CountryImg string
}

func CreateWpPostTmpl(post WpPost) string {
	tmpl := `
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
					<p>17</p>
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
						<a class="ca-list-item-img" href="#"><img style="max-width: 230px" src="` + link.Src + `" /></a>
						<div class="ca-list-item-content">
							<div class="ca-list-item-user">
								<figure style="background: url('/wp-content/uploads/2021/02/default.jpg') no-repeat center / cover;"></figure>
								<p class="ca-list-item-by">Asked by: ` + link.Author + `</p>
								<div class="ca-list-item-tag">Explainer</div>
							</div>
							<p class="ca-list-item-title">` + link.Title + `</p>
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
							<p class="ca-list-item-info-item-value"><img src="/wp-content/uploads/flags/` + link.CountryImg + `" /> <span style="color: #000;">` + link.CountryName + `</span></p>
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
