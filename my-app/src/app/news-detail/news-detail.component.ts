import { Component, OnInit, Input } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';

import { News } from '../news';
import { NewsService } from '../news.service';


@Component({
  selector: 'app-news-detail',
  templateUrl: './news-detail.component.html',
  styleUrls: ['./news-detail.component.css']
})
export class NewsDetailComponent implements OnInit {
  @Input() news: News;

  constructor(
    private route: ActivatedRoute,
    private newsService: NewsService,
    private location: Location
  ) { }

  ngOnInit() {
    this.getNews()
  }

  getNews(): void {
    const id = this.route.snapshot.paramMap.get('id');
    this.newsService.getNewsDetail(id)
      .subscribe(news => {
        NewsDetailComponent.newsReformat(news);
        this.news = news;
      });
  }

  goBack(): void {
    this.location.back();
  }

  // reformat news content 
  static newsReformat(item: News) {
    item.content = this.strReformat(item.content, /\s{2}/);
    item.content = this.strReformat(item.content, /\s{4}/);
  }

  static strReformat(content, target) {
    var arr = content.split(target);
    var i = 0;
    var text = "";
    for (; i < arr.length; i++) {
      text += "<p>" + arr[i] + "</p>";
    }
    if (text) {
      return text;
    }
    return content;
  }

}
