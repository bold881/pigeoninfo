import { Component, OnInit, ViewChild, OnChanges, OnDestroy,
          DoCheck, AfterContentInit, AfterContentChecked, 
          AfterViewInit, AfterViewChecked, ElementRef } from '@angular/core';

import { MatTableDataSource, MatPaginator, MatSort, } from '@angular/material';

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
  private selectedNews: News;

  //@ViewChild("news-detail") newsDetail: ElementRef;

  constructor(private newsService: NewsService, private er: ElementRef) {
  }

  ngOnInit() {
    this.getNews('');
  }

  getNews(day): void {
    if (!day) {
      var d = new Date();
      day = d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate();
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
    if (!this.dayNow) {
      this.dayNow = Date.now().toString();
    }
    var day = Date.parse(this.dayNow);
    day = day - 86400000;
    var preDay = new Date(day);
    this.dayNow = preDay.getFullYear() + '-' + (preDay.getMonth() + 1)
      + '-' + preDay.getDate();

    this.getNews(this.dayNow);
  }

  onNextClick() {
    if (!this.dayNow) {
      this.dayNow = Date.now().toString();
    }
    var day = Date.parse(this.dayNow);
    day = day + 86400000;
    var preDay = new Date(day);
    this.dayNow = preDay.getFullYear() + '-' + (preDay.getMonth() + 1)
      + '-' + preDay.getDate();

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
    if ("WebSocket" in window) {
      // Let us open a web socket
      var ws = new WebSocket("ws://localhost:4567/echo");

      ws.onopen = function () {
        // Web Socket is connected, send data using send()
        ws.send("Message to send");
        alert("Message is sent...");
      };

      ws.onmessage = function (evt) {
        var received_msg = evt.data;
        alert("Message is received...");
      };

      ws.onclose = function () {
        // websocket is closed.
        alert("Connection is closed...");
      };

      window.onbeforeunload = function (event) {
        ws.close();
      };
    } else {
      // The browser doesn't support WebSocket
      alert("WebSocket NOT supported by your Browser!");
    }
  }
}