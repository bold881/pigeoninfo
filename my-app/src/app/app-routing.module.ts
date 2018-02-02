import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { NewsComponent } from './news/news.component';
import { NewsDetailComponent } from './news-detail/news-detail.component';
import { SearchComponent } from './search/search.component';

const routes: Routes = [
  {path: 'news', component: NewsComponent },
  {path: 'news-detail/:id', component: NewsDetailComponent},
  {path: 'search', component: SearchComponent},
];


@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule { }
