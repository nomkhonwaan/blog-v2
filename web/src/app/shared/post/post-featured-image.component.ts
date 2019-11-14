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
    `
      img {
        max-width: 100%;
      }
    `,
  ],
})
export class PostFeaturedImageComponent extends PostComponent implements OnInit {

  src: string;

  ngOnInit(): void {
    if (this.post.featuredImage.slug.length > 0) {
      this.src = `http://localhost:8080/api/v2/storage/${this.post.featuredImage.slug}`;
      // this.src = `/api/v2/storage/${this.post.featuredImage.slug}`;
    }
  }

}
