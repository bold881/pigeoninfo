import { Component, OnInit, ViewChild, OnChanges } from '@angular/core';

import { MatTableDataSource, MatPaginator, MatSort,} from '@angular/material';

import { News } from '../news';
import { NewsService } from '../news.service';

@Component({
  selector: 'app-news',
  templateUrl: './news.component.html',
  styleUrls: ['./news.component.css']
})
export class NewsComponent implements OnInit{
  news: News[];
  dayNow = '';

  constructor(private newsService: NewsService) {
  }

  ngOnInit() {
    this.getNews('');
  }

  getNews(day): void {
    if(!day) {
      var d = new Date();
      day =  d.getFullYear() + '-' + (d.getMonth()+1) + '-' + d.getDate();
      this.dayNow = day;
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
    if(!this.dayNow) {
      this.dayNow = Date.now().toString();
    }
    var day = Date.parse(this.dayNow);
    day = day - 86400000;
    var preDay = new Date(day);
    this.dayNow =  preDay.getFullYear() + '-' + (preDay.getMonth()+1)
     + '-' + preDay.getDate();
    
    this.getNews(this.dayNow);
  }

  onNextClick() {
    if(!this.dayNow) {
      this.dayNow = Date.now().toString();
    }
    var day = Date.parse(this.dayNow);
    day = day + 86400000;
    var preDay = new Date(day);
    this.dayNow =  preDay.getFullYear() + '-' + (preDay.getMonth()+1)
     + '-' + preDay.getDate();
    
    this.getNews(this.dayNow);
  }
}