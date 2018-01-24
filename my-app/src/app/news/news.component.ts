import {
  Component, OnInit, ViewChild, OnChanges, OnDestroy,
  DoCheck, AfterContentInit, AfterContentChecked,
  AfterViewInit, AfterViewChecked, ElementRef,
  HostBinding, HostListener
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
  newsLite: NewsLite;

  //@ViewChild("news-detail") newsDetail: ElementRef;
  @HostListener('window:scroll', ['$event']) onScrollEvent($event) {
    if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
      this.getMoreData();
    }
  }

  constructor(private newsService: NewsService,
    private er: ElementRef
  ) { }

  ngOnInit() {
    //this.getNews('');
    this.getMoreData();
    this.test();
    NewsService.change.subscribe(d => {
      if (!this.dayNowChinese) {
        var day = Date.parse(this.dayNow);
        this.dayNowChinese = this.getChineseDayFromDate(new Date(day));
      }
      var tmp = JSON.parse(d);

      if (tmp.meta.includes(this.dayNowChinese)) {
        this.newsLite = JSON.parse(d);
      } else if (tmp.meta.includes(this.dayNow)) {
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
          newses.forEach(function (item) {
            item.content = item.content.substring(0, 200);
          });
          this.news = newses;
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
    if (item) {
      this.selectedNews = item;
      //NewsComponent.newsReformat(item);
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
    this.newsLite = { title: "", meta: "" };
    //this.getNews(this.dayNow);
    this.news = [];
    this.getMoreData();
  }

  // get year-month-day from given Date
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
      } else {
        d = td;
      }

      return dd.getFullYear() + '-' + m + '-' + d;
    }
  }

  // get 年月日from given Date
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
      } else {
        d = td;
      }

      return dd.getFullYear() + '年' + m + '月' + d + '日';
    }
  }

  getMoreData() {
    var ts;
    if (this.news) {
      ts = this.news[this.news.length - 1].sztime;
    } else {
      ts = "";
    }
    this.newsService.getNewsOfLimit(ts)
      .subscribe(
      newses => {
        if (newses) {
          newses.forEach(function (item) {
            item.content = item.content.substring(0, 200);
          });
          if(this.news) {
            this.news = this.news.concat(newses.slice(1));
          } else {
            this.news = newses;
          }
        }
      });
  }
}