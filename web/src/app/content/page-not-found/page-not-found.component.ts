import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-page-not-found',
  templateUrl: './page-not-found.component.html',
})
export class PageNotFoundComponent implements OnInit {

  constructor(private title: Title) { }

  ngOnInit(): void {
    this.title.setTitle(`Page not found - ${environment.title}`);
  }

}
