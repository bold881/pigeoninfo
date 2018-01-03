import { Injectable, Output, EventEmitter } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { News } from './news';

import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap } from 'rxjs/operators';

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type': 'application/x-www-form-urlencoded'
  })
}

@Injectable()
export class NewsService {

  private newsUrl = "http://101.200.47.113:4567/newsofday";

  static wsUrl = "ws://localhost:4567/echo";

  @Output() static change: EventEmitter<string> = new EventEmitter();

  constructor(
    private http: HttpClient,
  ) { }

  getNews(day): Observable<News[]> {
    if (!day) {
      return
    }
    return this.http.post<News[]>(this.newsUrl,
      day,
      httpOptions)
      .pipe(
      tap(sret => console.log(sret)),
      catchError(this.handleError('getNewses', []))
      );
  }

  initNewsWebsocket() {
    if ("WebSocket" in window) {
      // Let us open a web socket
      var ws = new WebSocket(NewsService.wsUrl);

      ws.onopen = function () {
        ws.send("Message to send");
      };

      ws.onmessage = function (evt) {
        //var received_msg = evt.data;
        NewsService.change.emit(evt.data);
       };

      ws.onclose = function () {
      };

      window.onbeforeunload = function (event) {
        ws.close();
      };
    } else {
    }
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      // this.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
