import { Injectable } from '@angular/core';
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

  constructor(
    private http: HttpClient
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
