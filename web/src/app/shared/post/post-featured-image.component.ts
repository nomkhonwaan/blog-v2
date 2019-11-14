import { PostComponent } from './post.component';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-post-featured-image',
  template: `
    <div>
      <img *ngIf="src" [src]="src">
    </div>
  `,
  styles: [
    `
      :host {
        display: block;
      }
    `,
  ],
})
export class PostFeaturedImageComponent extends PostComponent implements OnInit {

  src: string;

  ngOnInit(): void {
    if (this.post.featuredImage.slug.length > 0) {
      this.src = `/api/v2/storage/${this.post.featuredImage.slug}`;
    }
  }

}
