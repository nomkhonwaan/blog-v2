import { PostComponent } from './post.component';
import { Component, OnInit, Input, HostBinding } from '@angular/core';

@Component({
  selector: 'app-post-featured-image',
  template: `
    <div *ngIf="coverMode; else nonCoverMode">
      <ng-content></ng-content>
    </div>

    <ng-template #nonCoverMode>
      <img *ngIf="src" [src]="src">
    </ng-template>
  `,
  styles: [
    `
      :host {
        display: block;
        background: center/cover no-repeat fixed;
        background-image: url('http://placekitten.com/860/256');
      }
    `,
    `
      :host.-cover {
        margin: -8rem -3.2rem 2.4rem -3.2rem;
        padding: 0 3.2rem 3.2rem 3.2rem;
      }
    `,
    `
      :host.-cover ::ng-deep app-post-title,
      :host.-cover ::ng-deep app-post-tags a {
        color: #fafafa;
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

  /**
   * Used to indicate whether featured image should display as cover or not
   */
  @Input()
  @HostBinding('class.-cover')
  coverMode = false;

  @HostBinding('style.background')
  src: string;

  ngOnInit(): void {
    if (this.post.featuredImage.slug.length > 0) {
      this.src = `/api/v2/storage/${this.post.featuredImage.slug}`;
      this.src = 'http://localhost:8080' + this.src;
    }
  }

}
