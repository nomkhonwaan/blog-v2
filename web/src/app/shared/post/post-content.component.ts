import { Component, OnInit, Directive, ElementRef } from '@angular/core';

import { PostComponent } from './post.component';

@Directive({
  selector: 'img',
})
export class HTMLImageDirective implements OnInit {

  constructor(private el: ElementRef) { }

  ngOnInit(): void {
    console.log(this.el)
    console.log(this.el.nativeElement.getAttribute('src'));
  }

}

@Component({
  selector: 'app-post-content',
  template: `
    <article [innerHTML]="content"></article>
  `,
  styleUrls: ['./post-content.component.scss'],
})
export class PostContentComponent extends PostComponent implements OnInit {

  content: string;

  ngOnInit(): void {
    this.content = this.post.html.replace(new RegExp('/api/v1/attachments', 'g'), 'https://www.nomkhonwaan.com/api/v1/attachments');
  }
}
