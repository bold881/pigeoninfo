import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { NewsService } from '../news.service';
import { News } from "../news";
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.css']
})
export class SearchComponent implements OnInit, OnDestroy {
  keyword: string;
  news: News[];
  sub: Subscription;

  constructor(
    private route: ActivatedRoute,
    private newsService: NewsService
  ) { }

  ngOnInit() {
    this.getSearchResults();
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }

  getSearchResults() {
    this.sub = this.route.queryParams.subscribe(
      params => {
        this.keyword = params['keyword'];
        if (!this.keyword) {
          return;
        }
        this.newsService.getNewsOfSearch(this.keyword)
          .subscribe(newses => {
            if (newses) {
              newses.forEach(function (item) {
                item.content = item.content.substring(0, 200);
              });
              this.news = newses;
            }
          });
      }
    );
  }
}
