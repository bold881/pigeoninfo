import {
  Component, OnInit, ViewChild, OnChanges, OnDestroy,
  DoCheck, AfterContentInit, AfterContentChecked,
  AfterViewInit, AfterViewChecked, ElementRef,
  HostBinding
} from '@angular/core';

import { MatTableDataSource, MatPaginator, MatSort, } from '@angular/material';

import { News, NewsLite } from '../news';
import { NewsService } from '../news.service';

@Component({
  selector: 'app-news',
  templateUrl: './news.component.html',
  styleUrls: ['./news.component.css']
})
export class NewsComponent implements OnInit {
  news: News[];
  dayNow = '';
  dayNowChinese = '';
  private selectedNews: News;

  @HostBinding('class.is-open')
  newsLite : NewsLite;

  //@ViewChild("news-detail") newsDetail: ElementRef;

  constructor(private newsService: NewsService, 
    private er: ElementRef
  ) {}

  ngOnInit() {
    this.getNews('');
    this.test();
    NewsService.change.subscribe(d=> {
      if(!this.dayNowChinese) {
        var day = Date.parse(this.dayNow);
        this.dayNowChinese = this.getChineseDayFromDate(new Date(day));
      }
      var tmp = JSON.parse(d);
      
      if(tmp.meta.includes(this.dayNowChinese)) {
        this.newsLite = JSON.parse(d);
      }
    });
  }

  getNews(day): void {
    if (!day) {
      var d = new Date();
      day = this.getDayFromDate(d);
    }
    this.dayNow = day;
    this.newsService.getNews(day)
      .subscribe(
      newses => {
        if (newses) {
          this.news = newses.reverse();
        }
      }
      );
  }

  onPreviousClick() {
    if (!this.dayNow) {
      this.dayNow = Date.now().toString();
    }
    var day = Date.parse(this.dayNow);
    day = day - 86400000;
    var preDay = new Date(day);
    this.dayNow = this.getDayFromDate(preDay);

    this.getNews(this.dayNow);
  }

  onNextClick() {
    if (!this.dayNow) {
      this.dayNow = Date.now().toString();
    }
    var day = Date.parse(this.dayNow);
    day = day + 86400000;
    var preDay = new Date(day);
    this.dayNow = this.getDayFromDate(preDay);

    this.getNews(this.dayNow);
  }

  onListNewsClicked(item: News) {
    console.log(item);
    if (item) {
      this.selectedNews = item;
      try {
        this.er.nativeElement.querySelector('#news-detail').scrollTop = 0;
      } catch (err) {
        console.log(err);
      }
    }
  }

  test() {
    this.newsService.initNewsWebsocket();
  }

  onNewsAvailable() {
    this.newsLite = {title:"", meta:""};
    this.getNews(this.dayNow);
  }

  getDayFromDate(dd: Date) {
    if (dd) {
      var m, d;
      var tm = dd.getMonth() + 1;
      if (tm < 10) {
        m = '0' + tm;
      }

      var td = dd.getDate();
      if (td < 10) {
        d = '0' + td;
      }
      
      return dd.getFullYear() + '-' + m + '-' + d;
    }
  }

  getChineseDayFromDate(dd: Date) {
    if (dd) {
      var m, d;
      var tm = dd.getMonth() + 1;
      if (tm < 10) {
        m = '0' + tm;
      }

      var td = dd.getDate();
      if (td < 10) {
        d = '0' + td;
      }
      
      return dd.getFullYear() + '年' + m + '月' + d  + '日';
    }
  }
}