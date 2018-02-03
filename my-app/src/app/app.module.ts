import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { NewsComponent } from './news/news.component'
import { BrowserAnimationBuilder } from '@angular/platform-browser/animations/src/animation_builder';

import { NewsService } from './news.service';
import { AppRoutingModule } from './/app-routing.module';
import { NewsDetailComponent } from './news-detail/news-detail.component';
import { SearchComponent } from './search/search.component';
import { AboutComponent } from './about/about.component';


@NgModule({
  declarations: [
    AppComponent,
    NewsComponent,
    NewsDetailComponent,
    SearchComponent,
    AboutComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    HttpClientModule,
    BrowserAnimationsModule,
    AppRoutingModule,
  ],
  providers: [NewsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
